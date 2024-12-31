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
	"time"

	"gitee.com/deep-spark/ix-device-plugin/pkg/gpuallocator"
	"gitee.com/deep-spark/ix-device-plugin/pkg/ixml"
	"k8s.io/klog/v2"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const deviceName string = "iluvatar"

type iluvatarDevice struct {
	devSet *gpuallocator.DeviceSet

	stopCheckHeal chan struct{}

	deviceCh chan *gpuallocator.Device
}

func (d *iluvatarDevice) checkHealth() {
	klog.Infof("Start to GPU health checking.")

	d.devSet.Lk.Lock()
	LastCount := d.devSet.Count
	CurrentCount := d.devSet.Count
	d.devSet.Lk.Unlock()

	for {
		select {
		case <-d.stopCheckHeal:
			klog.Info("Stoping GPU health checking")

			return
		default:
		}
		time.Sleep(5 * time.Second)
		for _, dev := range d.devSet.Devices {
			for _, c := range dev.Chips {
				health, err := c.Operations.DeviceGetHealth()
				herr := ixml.CheckDeviceError(health)
				if err != nil {
					klog.Warningf("Unhealthy: dev:%v   err:%v\n", c.Device.ID, err)
					c.Health = pluginapi.Unhealthy
				} else if len(herr) > 0 {
					c.Health = pluginapi.Unhealthy
					klog.Warningf("Unhealthy Error Collection: dev:%v\n", c.Device.ID)
					for i, e := range herr {
						klog.Warningf("  Error(%d): %v\n", i, e)
					}
				} else {
					c.Health = pluginapi.Healthy
				}
			}
			if dev.UpdateHelath() {
				d.deviceCh <- dev
			}
		}

		d.devSet.Lk.Lock()
		CurrentCount = d.devSet.Count
		d.devSet.Lk.Unlock()
		if CurrentCount != LastCount {
			d.deviceCh <- &gpuallocator.Device{Replicas: -1}
			LastCount = CurrentCount
		}
	}
}
