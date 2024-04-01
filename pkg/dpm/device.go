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

package dpm

import (
	"strings"

	"github.com/golang/glog"
        "gitee.com/deep-spark/ix-device-plugin/pkg/ixml"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const deviceName string = "iluvatar"

type device struct {
	pluginapi.Device

	name  string
	minor uint
	minorslice []uint
	uuid  string
	uuidslice []string
	index uint
}

type iluvatarDevice struct {
	devices []*device

	stopCheckHeal chan struct{}

	deviceCh chan *pluginapi.Device
}

func getDeviceCount() (uint, error) {
	count, err := ixml.GetDeviceCount()

	return count, err
}

func buildDevice(index uint, d ixml.Device) *device {
	var err error
	dev := device{}

	dev.name, err = d.DeviceGetName()
	if err != nil {
		glog.Errorf("Failed to get device name: %v", err)
	}

	dev.minor, err = d.DeviceGetMinorNumber()
	if err != nil {
		glog.Errorf("Failed to get device minor number: %v", err)
	}

	dev.minorslice, err = d.DeviceGetMinorSlice()
	if err != nil {
		glog.Errorf("Failed to get device minor number slice: %v", err)
	}

	dev.uuid, err = d.DeviceGetUUID()
	if err != nil {
		glog.Errorf("Failed to get device uuid: %v", err)
	}

	dev.uuidslice, err = d.DeviceGetUUIDSlice()
	if err != nil {
		glog.Errorf("Failed to get device uuid slice: %v", err)
	}

	dev.index = index
	dev.ID = dev.uuid
	dev.Health = pluginapi.Healthy

	glog.Infof("Detected device: %d, name: %s, uuid: %s", dev.index, dev.name, dev.uuid)

	return &dev
}

func newDevice() []*device {
	var devs []*device

	count, _ := getDeviceCount()

	for i := uint(0); i < count; i++ {
		dev, err := ixml.NewDeviceByIndex(i)
		if err != nil {
			glog.Errorf("Failed to get device-%d handle: %v", i, err)
			continue
		}

		devs = append(devs, buildDevice(i, dev))
	}

	return devs
}

func (d *iluvatarDevice) cachedDevices() []*pluginapi.Device {
	var devs []*pluginapi.Device

	for _, d := range d.devices {
		flag := false
		for _,v := range d.minorslice {
			if (v < d.minor) {
			   flag = true
			}
		}

		if flag == false {
			devs = append(devs, &d.Device)
		}
	}

	return devs
}

func (d *iluvatarDevice) deviceExist(id string) bool {
	for _, d := range d.cachedDevices() {
		if d.ID == id {
			return true
		}
	}
	return false
}

func (d *iluvatarDevice) checkHealth() {
	eventSet, _ := ixml.NewEventSet()
	defer eventSet.EventSetFree()

	glog.Infof("Start to GPU health checking.")

	for _, dev := range d.cachedDevices() {
		err := eventSet.RegisterEventsForDevice(dev.ID, ixml.XidCriticalError)
		if err != nil && strings.HasSuffix(err.Error(), "Not Supported") {
			glog.Warningf("%s is too old to support healthchecking: %s. Marking it unhealthy.", dev.ID, err)
			d.deviceCh <- dev

			continue
		}
	}

	for {
		select {
		case <-d.stopCheckHeal:
			glog.Info("Stoping GPU health checking")

			return
		default:
		}

		eventData, err := eventSet.WaitForEvent(5000)
		if err != nil && eventData.Type != ixml.XidCriticalError {
			continue
		}

		if eventData.Data == 31 || eventData.Data == 43 || eventData.Data == 45 {
			continue
		}

		for _, dev := range d.cachedDevices() {
			glog.Info("Error")
			d.deviceCh <- dev
		}
	}
}
