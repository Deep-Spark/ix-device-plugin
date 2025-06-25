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
	"math"
	"strconv"
	"strings"
	"time"

	"gitee.com/deep-spark/ix-device-plugin/pkg/gpuallocator"
	"gitee.com/deep-spark/ix-device-plugin/pkg/ixml"
	"gitee.com/deep-spark/ix-device-plugin/pkg/kube"
	"k8s.io/klog/v2"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const deviceName string = "iluvatar"
const updatePeriod = 5

type iluvatarDevice struct {
	devSet *gpuallocator.DeviceSet

	stopCheckHeal chan struct{}

	deviceCh chan *gpuallocator.Device

	updateCh chan struct{}

	kubeclient *kube.KubeClient
}

func (d *iluvatarDevice) notifyUpdate() {
	if d.kubeclient != nil {
		d.updateCh <- struct{}{}
	}
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
				d.notifyUpdate()
			}
		}

		d.devSet.Lk.Lock()
		CurrentCount = d.devSet.Count
		d.devSet.Lk.Unlock()
		if CurrentCount != LastCount {
			d.deviceCh <- &gpuallocator.Device{Replicas: -1}
			d.notifyUpdate()
			LastCount = CurrentCount
		}
	}
}

func (d *iluvatarDevice) updateDeviceinfo() {
	klog.Infof("Start to update deviceinfo.")

	ticker := time.NewTicker(updatePeriod * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-d.stopCheckHeal:
			klog.Info("Stoping update deviceinfo")

			return
		case <-ticker.C:
			d.updatingDevicelist(false)
		case <-d.updateCh:
			d.updatingDeviceinfo(true)
		}
	}
}

func (d *iluvatarDevice) updatePodAnnotation(verbose bool) {
	podList := d.kubeclient.GetActivePodListCache()

	podDeviceInfo, err := kube.GetKltAndRealAllocateDev(podList)
	if err != nil {
		klog.Errorf("get pod device info failed, %v", err)
		return
	}

	for _, deviceInfo := range podDeviceInfo {
		if verbose {
			klog.Infof("pods: %s, %s, %s", deviceInfo.Pod.Name, deviceInfo.Pod.Status.Phase, deviceInfo.Pod.UID)
		}

		_, existRealAlloc := deviceInfo.Pod.Annotations[kube.ResourceNamePrefix+kube.PodDevRealAlloc]
		if existRealAlloc {
			continue
		}
		if len(deviceInfo.KltDevice) == 0 || len(deviceInfo.RealDevice) == 0 {
			klog.Warningf("%s %s klt device or real device is empty", deviceInfo.Pod.Namespace,
				deviceInfo.Pod.Name)
			continue
		}

		annotations := map[string]string{
			kube.ResourceNamePrefix + kube.PodDevKubelet:   strings.Join(deviceInfo.KltDevice, ","),
			kube.ResourceNamePrefix + kube.PodDevRealAlloc: strings.Join(deviceInfo.RealDevice, ","),
		}

		klog.Infof("Update annotation for %s, %v", deviceInfo.Pod.Name, annotations)

		d.kubeclient.TryUpdatePodCacheAnnotation(&deviceInfo.Pod, annotations)
	}

}

func (d *iluvatarDevice) updatingDevicelist(verbose bool) {

	d.updatePodAnnotation(verbose)

	devices := []string{}
	allocated := d.GetAllocatedDevicesFromPodCache(verbose)

	for _, dev := range d.devSet.Devices {
		for _, rdev := range dev.Exposed {
			if !allocated[rdev.ID] && rdev.Health == pluginapi.Healthy {
				devices = append(devices, rdev.ID)
			}
		}
	}
	err := d.kubeclient.WriteDevicelisttoCM(devices)
	if err != nil {
		klog.Errorf("failed to write devicelist to configmap: %v", err)
	}
}

func (d *iluvatarDevice) updatingDeviceinfo(verbose bool) {

	d.updatePodAnnotation(verbose)

	deviceinfomap := map[string]kube.DeviceInfo{}
	devices := []string{}
	allocated := d.GetAllocatedDevicesFromPodCache(verbose)

	for _, dev := range d.devSet.Devices {
		for _, rdev := range dev.Exposed {
			if !allocated[rdev.ID] && rdev.Health == pluginapi.Healthy {
				devices = append(devices, rdev.ID)
			}
		}

		deviceinfo := kube.DeviceInfo{
			Name:  dev.Name,
			UUID:  dev.UUID,
			Links: map[string][]kube.P2PLink{},
		}

		for uuid, links := range dev.Links {
			deviceinfo.Links[uuid] = []kube.P2PLink{}
			for _, link := range links {
				deviceinfo.Links[uuid] = append(deviceinfo.Links[uuid], kube.P2PLink{
					TypeName:  gpuallocator.P2PLinkTypeToString(link.Type),
					TypeIndex: link.Type,
				})
			}
		}
		deviceinfomap[dev.UUID] = deviceinfo
	}
	err := d.kubeclient.WriteDeviceInfoDataIntoCM(devices, deviceinfomap, verbose)
	if err != nil {
		klog.Errorf("failed to write deviceinfo to configmap: %v", err)
	}
}

func (d *iluvatarDevice) GetAllocatedDevicesFromPodCache(verbose bool) map[string]bool {
	allocated := map[string]bool{}
	devices := []string{}
	podlist := d.kubeclient.GetActivePodListCache()
	for _, pod := range podlist {
		devStr, exist := pod.Annotations[kube.ResourceNamePrefix+kube.PodDevRealAlloc]
		if !exist {
			continue
		}

		devs := strings.Split(devStr, kube.CommaSepDev)

		for _, dev := range devs {
			allocated[dev] = true
			devices = append(devices, dev)
		}

		if verbose {
			klog.Infof("Pod %s real allocated devices: %v", pod.Name, devs)
		}
	}
	if verbose {
		klog.Infof("Allocated devices: %v", devices)
	}
	return allocated
}

func (d *iluvatarDevice) UseVolcano(devices []string) ([]string, bool) {
	pod, err := d.kubeclient.GetMatchedPod(devices)
	if err != nil {
		klog.Errorf("Failed to get matched pod: %v", err)
		return devices, false
	} else {
		annotations := map[string]string{
			kube.PodPredicateTime: strconv.FormatUint(math.MaxUint64, 10),
		}
		d.kubeclient.TryUpdatePodCacheAnnotation(pod, annotations)

		vDevs, ok := pod.Annotations[kube.ResourceNamePrefix+kube.PodDevVolcano]
		if ok {
			return strings.Split(vDevs, kube.CommaSepDev), true
		} else {
			klog.Errorf("Failed to get volcano devices from pod %s", pod.Name)
			return devices, false
		}
	}
}
