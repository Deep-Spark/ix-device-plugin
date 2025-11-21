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
	Cfg      *config.Config
	Replicas int
	ixmlLock sync.Mutex
}

type DeviceList []*Device
type ReplicaDeviceMap map[string]ReplicaDevice
type Alias string
type moduleContext struct {
	ChipMap         map[string]*Chip
	unManagedChip   map[string]*Chip
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

func (d *Device) GenerateIndexs() []int {
	var ret []int
	for _, c := range d.Chips {
		ret = append(ret, int(c.Index))
	}
	return ret
}

func buildReplicaDevice(dev pluginapi.Device, parent *Device) *ReplicaDevice {
	return &ReplicaDevice{Device: dev, Parent: parent}
}

func buildChip(index uint, d ixml.Device) *Chip {
	var err error
	chip := Chip{Operations: d}

	chip.UUID, err = d.DeviceGetUUID()
	if err != nil {
		klog.Errorf("Failed to get device uuid: %v", err)
		return nil
	}

	chip.Name, err = d.DeviceGetName()
	if err != nil {
		klog.Errorf("Failed to get device name: %v", err)
	}

	chip.Minor, err = d.DeviceGetMinorNumber()
	if err != nil {
		klog.Errorf("Failed to get device minor number: %v", err)
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
				libctx.unManagedChip[chip.UUID] = chip
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

func (d *Device) UpdateHealth() bool {
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

func BuildDeviceSet(cfg *config.Config) *DeviceSet {
	chips, err := scanAllChips()
	if err != nil {
		return nil
	}
	var ds DeviceSet
	ds.Cfg = cfg
	ds.Replicas = cfg.Sharing.TimeSlicing.Replicas

	reconcileDeviceSet(&ds, chips)

	return &ds
}

func (d *DeviceSet) updateDeviceEvent() {
	klog.Info("Start Update DeviceSet")
	// Always do a full scan + full rebuild
	chips, err := scanAllChips()
	if err != nil {
		klog.Infof("get device count failed, Failed to update Udev event.")
		return
	}

	libctx.timeLck.Lock()
	libctx.timeArg = &chips
	if libctx.timeEventDone {
		// First event after idle: schedule a delayed rebuild
		libctx.timeEventDone = false
		libctx.timeProcessUdev = time.AfterFunc(libctx.duration, func() {
			// Actually rebuild after the debounce window
			local := *libctx.timeArg
			reconcileDeviceSet(d, local)

			libctx.timeLck.Lock()
			libctx.timeArg = nil
			libctx.timeEventDone = true
			libctx.timeLck.Unlock()
		})
	} else {
		// Merge multiple events: if another event arrives before the timer fires,
		// update timeArg with the latest chip list
		libctx.timeArg = &chips
	}
	libctx.timeLck.Unlock()
}

func (d *DeviceSet) UpdateUdev(dev *udev.Device) {
	action := dev.Action()
	switch action {
	case "add":
		klog.Infof("-- Add    -- udev event\n")
		d.updateDeviceEvent()
	case "remove":
		klog.Infof("-- Remove -- udev event\n")
		d.updateDeviceEvent()
	case "change":
		klog.Infof("-- Change -- udev event\n")
		d.updateDeviceEvent()
	default:
		klog.Infof("[%v] udev event (ignored)\n", action)
	}
}

func scanAllChips() ([]*Chip, error) {
	klog.Info("Start scan all chips")
	var chips []*Chip

	count, err := ixml.GetDeviceCount()
	if err != nil {
		klog.Infof("get device count failed.")
		return nil, err
	}
	klog.Infof("IXML device count = %d", count)

	libctx.ChipMap = make(map[string]*Chip)
	for i := uint(0); i < count; i++ {
		devHandler, err := ixml.NewDeviceByIndex(i)
		if err != nil {
			klog.Errorf("Failed to get device-%d handle: %v", i, err)
			continue
		}
		c := buildChip(i, devHandler)
		if c == nil {
			klog.Error("Undetected Chip")
			continue
		}
		chips = append(chips, c)
		libctx.ChipMap[c.UUID] = c
	}

	klog.Infof("Real device count = %d", len(chips))
	return chips, nil
}

func reconcileDeviceSet(ds *DeviceSet, chips []*Chip) {
	klog.Info("Reconcile DeviceSet")
	if ds == nil {
		return
	}
	ds.Lk.Lock()
	defer ds.Lk.Unlock()

	// rebuild DeviceSet
	ds.Devices = make(map[string]*Device)
	libctx.unManagedChip = make(map[string]*Chip)

	if ds.Cfg.Flags.SplitBoard {
		processSingleChip(ds, chips)
	} else {
		processMultiChip(ds, chips)
	}
	resetTopological(&ds.Devices)

	ds.Count = uint(len(chips))
	libctx.count = ds.Count

	ds.ShowLayout()
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
