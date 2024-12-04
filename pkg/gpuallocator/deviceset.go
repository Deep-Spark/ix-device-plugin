/*
Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gpuallocator

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitee.com/deep-spark/ix-device-plugin/pkg/config"
	"gitee.com/deep-spark/ix-device-plugin/pkg/ixml"

	udev "github.com/jochenvg/go-udev"
	"k8s.io/klog/v2"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type Chip struct {
	Name       string
	Minor      uint
	UUID       string
	Index      uint
	Operations ixml.Device
	pluginapi.Device
}

type ReplicaDevice struct {
	pluginapi.Device
	Parent *Device
}

type Device struct {
	Name      string
	UUID      string
	Minor     uint
	IsMulChip bool
	Index     *uint
	Exposed   []*ReplicaDevice

	Chips    map[string]*Chip
	Links    map[string][]P2PLink
	Replicas int
}

type DeviceSet struct {
	Devices  map[string]*Device
	Lk       sync.Mutex
	Count    uint
	cfg      *config.Config
	Replicas int
}

type DeviceList []*Device
type ReplicaDeviceMap map[string]ReplicaDevice
type Alias string
type moduleContext struct {
	ChipMap         map[string]*Chip
	unManagedChip   map[string]*Chip
	MainSet         *DeviceSet
	timeLck         sync.Mutex
	timeEventDone   bool
	timeProcessUdev *time.Timer
	duration        time.Duration
	timeArg         *[]*Chip
	count           uint
}

var (
	libctx = moduleContext{
		ChipMap:         make(map[string]*Chip),
		unManagedChip:   make(map[string]*Chip),
		count:           0,
		timeEventDone:   true,
		timeProcessUdev: nil,
		duration:        time.Duration(5) * time.Second,
		timeArg:         nil,
	}
)

func (a Alias) HasAlias() bool {
	slc := strings.Split(string(a), "::")
	if len(slc) == 2 {
		return true
	} else {
		return false
	}
}

func (a Alias) Prefix() string {
	cnt := strings.Split(string(a), "::")
	return cnt[0]
}

func (a Alias) Suffix() string {
	cnt := strings.Split(string(a), "::")
	return cnt[1]
}

func (a Alias) GetValue() (string, string) {
	cnt := strings.Split(string(a), "::")
	return cnt[0], cnt[1]
}

func buildReplicaDevice(dev pluginapi.Device, parent *Device) *ReplicaDevice {
	return &ReplicaDevice{Device: dev, Parent: parent}
}

func buildChip(index uint, d ixml.Device) *Chip {
	var err error
	chip := Chip{Operations: d}

	chip.Name, err = d.DeviceGetName()
	if err != nil {
		klog.Errorf("Failed to get device name: %v", err)
	}

	chip.Minor, err = d.DeviceGetMinorNumber()
	if err != nil {
		klog.Errorf("Failed to get device minor number: %v", err)
	}

	chip.UUID, err = d.DeviceGetUUID()
	if err != nil {
		klog.Errorf("Failed to get device uuid: %v", err)
	}

	hasNuma, numa, err := d.DeviceGetNumaNode()
	if err != nil {
		klog.Errorf("Failed to get pci info: %v", err)
	}

	health, err := d.DeviceGetHealth()
	herr := ixml.CheckDeviceError(health)
	if err != nil {
		klog.Warningf("Unhealthy: dev:%v   err:%v\n", chip.UUID, err)
		chip.Health = pluginapi.Unhealthy
	} else if len(herr) > 0 {
		klog.Warningf("Unhealthy: dev:%v   herr:%v\n", chip.UUID, herr)
		chip.Health = pluginapi.Unhealthy
	} else {
		chip.Health = pluginapi.Healthy
	}

	chip.Index = index
	chip.ID = chip.UUID
	if hasNuma {
		chip.Topology = &pluginapi.TopologyInfo{
			Nodes: []*pluginapi.NUMANode{
				{
					ID: int64(numa),
				},
			},
		}
	}

	klog.Infof("Detected Chip: %d, name: %s, uuid: %s  numa:%v", chip.Index, chip.Name, chip.UUID, numa)

	return &chip
}

func buildDevice(c *Chip, replicas int) *Device {
	dev := Device{
		Name:     c.Name,
		UUID:     c.UUID,
		Minor:    c.Minor,
		Index:    &c.Index,
		Replicas: replicas,
	}

	dev.Chips = make(map[string]*Chip)
	dev.Chips[dev.UUID] = c
	if replicas == 0 {
		replicaDev := buildReplicaDevice(c.Device, &dev)
		dev.Exposed = append(dev.Exposed, replicaDev)
	} else {
		for i := 0; i < replicas; i++ {
			replicasID := fmt.Sprintf("%s::%d", c.UUID, i)
			replicaDev := buildReplicaDevice(c.Device, &dev)
			replicaDev.ID = replicasID
			dev.Exposed = append(dev.Exposed, replicaDev)
		}
	}
	return &dev
}

func resetTopological(devs *map[string]*Device) {
	for _, d := range *devs {
		d.Links = make(map[string][]P2PLink)
	}

	for i, d1 := range *devs {
		c1 := d1.GetMasterChip()
		if c1 == nil {
			continue
		}
		for j, d2 := range *devs {
			c2 := d2.GetMasterChip()
			if c2 == nil {
				continue
			}
			if i != j {
				p2plink := GetP2PLink(c1.Operations, c2.Operations)
				if p2plink != P2PLinkUnknown {
					d1.Links[j] = append(d1.Links[j], P2PLink{d2, p2plink})
				}

				//TODO nvlink assign
			}
		}
	}
}

func processSingleChip(sideEffect *DeviceSet, ChipList []*Chip) {
	for _, chip := range ChipList {
		dev := buildDevice(chip, sideEffect.Replicas)
		sideEffect.Devices[dev.UUID] = dev
	}
}

// Shall lock gpuallocator global contexts
func processMultiChip(sideEffect *DeviceSet, ChipList []*Chip) {
	var Mul []*Chip
	for _, chip := range libctx.unManagedChip {
		Mul = append(Mul, chip)
	}
	for _, chip := range ChipList {
		isSupport, pos := chip.Operations.DeviceGetBoardPosition()
		if isSupport {
			if pos == 0 {
				dev := buildDevice(chip, sideEffect.Replicas)
				sideEffect.Devices[dev.UUID] = dev
				dev.IsMulChip = true
			} else {
				Mul = append(Mul, chip)
				libctx.unManagedChip[chip.ID] = chip
			}
		} else {
			dev := buildDevice(chip, sideEffect.Replicas)
			sideEffect.Devices[dev.UUID] = dev
		}
	}

	for _, dev := range sideEffect.Devices {
		if dev.IsMulChip {
			MasterChip := dev.GetMasterChip()
			if MasterChip == nil {
				continue
			}
			for i := 0; i < len(Mul); i++ {
				if Mul[i] == nil {
					continue
				}
				_, onSameBoard := ixml.GetDeviceOnSameBoard(MasterChip.Operations, Mul[i].Operations)
				if onSameBoard {
					dev.Chips[Mul[i].UUID] = Mul[i]
					if Mul[i].Health == pluginapi.Unhealthy {
						dev.SetUnHealth()
					}
					delete(libctx.unManagedChip, Mul[i].UUID)

					Mul[i] = nil
				}
			}

			//MR150 has 2 chips, only 1 work shall label unhealthy
			if len(dev.Chips) != 2 {
				dev.SetUnHealth()
			}
		}
	}

	for _, chip := range libctx.unManagedChip {
		klog.Info("Warning: still have chips is not recognized :%v", chip)
	}
}

func BuildDeviceSet(cfg *config.Config) *DeviceSet {
	var ret DeviceSet
	ret.Devices = make(map[string]*Device)
	var ChipList []*Chip
	ret.cfg = cfg
	ret.Replicas = cfg.Sharing.TimeSlicing.Replicas

	count, err := ixml.GetDeviceCount()
	if err != nil {
		klog.Infof("get device count failed.")
		return nil
	}

	for i := uint(0); i < count; i++ {
		devHandler, err := ixml.NewDeviceByIndex(i)
		if err != nil {
			klog.Errorf("Failed to get device-%d handle: %v", i, err)
			continue
		}

		dev := buildChip(i, devHandler)
		ChipList = append(ChipList, dev)
		libctx.ChipMap[dev.UUID] = dev
	}

	libctx.MainSet = &ret
	libctx.MainSet.Lk.Lock()

	if ret.cfg.Flags.SplitBoard {
		processSingleChip(&ret, ChipList)
	} else {
		processMultiChip(&ret, ChipList)
	}
	resetTopological(&ret.Devices)
	ret.Count = count
	libctx.count = count

	libctx.MainSet.Lk.Unlock()
	return &ret
}

func (d *Device) GetMasterChip() *Chip {
	chip, ok := d.Chips[d.UUID]
	if !ok {
		return nil
	} else {
		return chip
	}
}

func (d *Device) SetHealth() {
	for _, d := range d.Exposed {
		d.Health = pluginapi.Healthy
	}
}

func (d *Device) SetUnHealth() {
	for _, d := range d.Exposed {
		d.Health = pluginapi.Unhealthy
	}
}

func (d *Device) GenerateSpecList() []*pluginapi.DeviceSpec {
	var ret []*pluginapi.DeviceSpec

	for _, c := range d.Chips {
		d := pluginapi.DeviceSpec{}
		// Expose the device node for iluvatar pod.
		d.HostPath = config.HostPathPrefix + config.DeviceName + strconv.Itoa(int(c.Minor))
		d.ContainerPath = config.ContainerPathPrefix + config.DeviceName + strconv.Itoa(int(c.Minor))
		d.Permissions = "rw"
		ret = append(ret, &d)
	}
	return ret
}

func (d *Device) GenerateIDS() []string {
	var ret []string
	for _, c := range d.Chips {
		ret = append(ret, c.UUID)
	}
	return ret
}

func (d *Device) UpdateHelath() bool {
	ret := false

	//check wheter is offline card
	if d.IsMulChip && len(d.Chips) != 2 {
		if d.Exposed[0].Health == pluginapi.Healthy {
			d.SetUnHealth()
			ret = true
		}
		return ret
	}

	for _, c := range d.Chips {
		if c.Health == pluginapi.Unhealthy {
			if d.Exposed[0].Health == pluginapi.Healthy {
				d.SetUnHealth()
				ret = true
			}
			return ret
		}
	}

	if d.Exposed[0].Health == pluginapi.Unhealthy {
		d.SetHealth()
		ret = true
	}
	return ret
}

func (d *DeviceSet) CachedDevices() []*pluginapi.Device {
	var devs []*pluginapi.Device
	var cp *pluginapi.Device
	for _, d := range d.Devices {
		for _, replica := range d.Exposed {
			cp = new(pluginapi.Device)
			*cp = replica.Device
			devs = append(devs, cp)
		}
	}

	return devs
}

func (d *DeviceSet) DeviceExist(id string) bool {
	prefix := Alias(id).Prefix()

	if dev, ok := d.Devices[prefix]; ok {
		for _, d := range dev.Exposed {
			if d.ID == id {
				return true
			}
		}
	}

	return false
}

func timeUdev() {
	libctx.timeLck.Lock()
	ChipList := *libctx.timeArg

	libctx.MainSet.Lk.Lock()
	if libctx.MainSet.cfg.Flags.SplitBoard {
		processSingleChip(libctx.MainSet, ChipList)
	} else {
		processMultiChip(libctx.MainSet, ChipList)
	}
	resetTopological(&libctx.MainSet.Devices)

	libctx.MainSet.Count = libctx.count

	libctx.timeArg = nil
	libctx.timeEventDone = true

	libctx.MainSet.Lk.Unlock()
	libctx.timeLck.Unlock()
	libctx.MainSet.ShowLayout()
}

func (d *DeviceSet) updateDeviceEvent() {
	var ChipList *[]*Chip
	if libctx.timeArg != nil {
		ChipList = libctx.timeArg
	} else {
		ChipList = new([]*Chip)
	}

	klog.Infof("count %v\n", libctx.count)
	count, err := ixml.GetDeviceCount()
	if err != nil {
		klog.Infof("get device count failed, Failed to update Udev event.")
		return
	}
	klog.Infof("new count %v\n", count)
	if count == libctx.count {
		klog.Infof("count is not changed, do nothing\n")
		return
	} else {
		IndexSet := make(map[uint]bool)
		for _, dev := range libctx.ChipMap {
			newIdx, err := dev.Operations.DeviceGetIndex()
			if err != nil {
				klog.Errorf("Failed to get device-%d index, and set to 0xffff", dev.UUID)
				newIdx = 0xffff
			}
			IndexSet[newIdx] = true
			dev.Index = newIdx
		}
		for idx := uint(0); idx < count; idx++ {
			if _, ok := IndexSet[idx]; ok {
				continue
			}
			devHandler, err := ixml.NewDeviceByIndex(idx)
			if err != nil {
				klog.Errorf("Failed to get device-%d handle: %v", idx, err)
				continue
			}

			dev := buildChip(idx, devHandler)
			_, ok := libctx.ChipMap[dev.UUID]
			if ok {
				klog.Infof("Error: duplicated uuid-->  %v", dev.UUID)
			} else {
				klog.Infof("Added Device with the New uuid-->  %v", dev.UUID)

				libctx.timeLck.Lock()
				*ChipList = append(*ChipList, dev)
				libctx.timeLck.Unlock()

				libctx.ChipMap[dev.UUID] = dev
			}
		}

		libctx.timeLck.Lock()
		libctx.count = count
		if libctx.timeEventDone {
			libctx.timeEventDone = false
			libctx.timeProcessUdev = time.AfterFunc(libctx.duration, timeUdev)
		}
		libctx.timeArg = ChipList
		libctx.timeLck.Unlock()
	}
}

func (d *DeviceSet) UpdateUdev(dev *udev.Device) {
	action := dev.Action()
	switch action {
	case "add":
		klog.Infof("-- Add    -- udev event\n")
		d.updateDeviceEvent()
	case "remove":
		klog.Infof("-- Remove -- udev event\n")
	default:
		klog.Infof("[%v] udev event\n", action)
	}
}

func (d *DeviceSet) ShowLayout() {
	for id, dev := range d.Devices {
		klog.Infof("Dev ID:%v\n", id)
		for tid, linkLst := range dev.Links {
			klog.Infof("\t\t--->\tTarget ID:%v", tid)
			for _, link := range linkLst {
				klog.Infof("\t\t\t\t---->\tLink Type:%v", link.Type)
			}
		}
	}
}

func (d *DeviceSet) Filter(uuids []string) (DeviceList, error) {
	var filtered DeviceList
	for _, uuid := range uuids {
		if dev, ok := d.Devices[uuid]; ok {
			filtered = append(filtered, dev)
		} else {
			return nil, fmt.Errorf("no device with uuid: %v", uuid)
		}
	}

	return filtered, nil
}

func (d *DeviceSet) GetTotalCount() int {
	sum := 0
	for _, dev := range d.Devices {
		if dev.IsMulChip {
			sum += 2
		} else {
			sum += 1
		}
	}

	return sum
}

func (d *DeviceSet) BuildReplicaMap() ReplicaDeviceMap {
	res := make(ReplicaDeviceMap)
	for _, dev := range d.Devices {
		for _, replicaDev := range dev.Exposed {
			res[replicaDev.ID] = *replicaDev
		}
	}
	return res
}

func (dm ReplicaDeviceMap) Contains(ids ...string) bool {
	for _, id := range ids {
		if _, exists := dm[id]; !exists {
			return false
		}
	}
	return true
}

// Subset returns the subset of devices in Devices matching the provided ids.
// If any id in ids is not in Devices, then the subset that did match will be returned.
func (dm ReplicaDeviceMap) Subset(ids []string) ReplicaDeviceMap {
	res := make(ReplicaDeviceMap)
	for _, id := range ids {
		if dm.Contains(id) {
			res[id] = dm[id]
		}
	}
	return res
}

// Difference returns the set of devices contained in ds but not in ods.
func (dm ReplicaDeviceMap) Difference(ods ReplicaDeviceMap) ReplicaDeviceMap {
	res := make(ReplicaDeviceMap)
	for id := range dm {
		if !ods.Contains(id) {
			res[id] = dm[id]
		}
	}
	return res
}

func (dm ReplicaDeviceMap) GetIDs() []string {
	var res []string
	for _, d := range dm {
		res = append(res, d.ID)
	}
	return res
}

func PrefixUUID(uuids []string) ([]string, map[string][]string) {
	var res []string
	set := make(map[string][]string)
	for _, u := range uuids {
		key := Alias(u).Prefix()
		if _, ok := set[key]; !ok {
			res = append(res, key)
		}
		set[key] = append(set[key], u)
	}
	return res, set
}
