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
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"gitee.com/deep-spark/ix-device-plugin/pkg/config"
	"gitee.com/deep-spark/ix-device-plugin/pkg/gpuallocator"
	"golang.org/x/net/context"
	"k8s.io/klog/v2"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

// resourceName is the name to identify iluvatar device plugin
var ResourceName string = "iluvatar.com/gpu"

// iluvatarDevicePlugin is the implementation of iluvatar device plugin
type iluvatarDevicePlugin struct {
	iluvatarDevice

	name string

	stopList chan struct{}
}

// GetDevicePluginOptions returns the values of the optional settings for this plugin
func (p *iluvatarDevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{GetPreferredAllocationAvailable: true}, nil
}

// PreStartContainer allows kubelet to pass reinitialized devices to containers.
func (p *iluvatarDevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	klog.Info("PreStartContainer...")
	return &pluginapi.PreStartContainerResponse{}, nil
}

// ListAndWatch lists devices
func (p *iluvatarDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	devs := p.devSet.CachedDevices()

	klog.Info("Start to list and watch GPU.")

	for _, dev := range devs {
		klog.Infof("L->    %v\n", dev)
	}

	s.Send(&pluginapi.ListAndWatchResponse{Devices: devs})

	for {
		select {
		case <-p.stopList:
			klog.Info("Stoping list and watch GPU.")

			return nil
		case dev := <-p.deviceCh:
			devs := p.devSet.CachedDevices()
			if dev.Replicas == -1 {
				for _, dev := range devs {
					klog.Infof("L->    %v\n", dev)
				}
			} else {
				if dev.Exposed[0].Health == pluginapi.Unhealthy {
					klog.Infof("'%s' device marked unhealthy: %s", p.name, dev.UUID)
				} else {
					klog.Infof("'%s' device marked healthy: %s", p.name, dev.UUID)
				}
			}
			s.Send(&pluginapi.ListAndWatchResponse{Devices: devs})
		}
	}
}

// GetPreferredAllocation returns the preferred allocation from the set of devices specified in the request
func (p *iluvatarDevicePlugin) GetPreferredAllocation(ctx context.Context, r *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	klog.Info("Start to GetPreferred Allocation.")
	response := &pluginapi.PreferredAllocationResponse{}

	for _, req := range r.ContainerRequests {
		IDs, err := p.alignedAlloc(req.AvailableDeviceIDs, req.MustIncludeDeviceIDs, int(req.AllocationSize))
		if err != nil {
			klog.Infof("can't use prefered functionality:%v\n", err)
			return nil, err
		}
		resp := &pluginapi.ContainerPreferredAllocationResponse{DeviceIDs: IDs}
		response.ContainerResponses = append(response.ContainerResponses, resp)
	}
	return response, nil
}

// GetPreferredAllocation returns the preferred allocation from the set of devices specified in the request
func (p *iluvatarDevicePlugin) alignedAlloc(available, required []string, size int) ([]string, error) {
	var devices []string
	if p.devSet.Replicas > 0 {
		arg := gpuallocator.ReplicaPolicyArgs{Device: p.devSet.BuildReplicaMap(), Available: available, Required: required, Size: size}
		devices = gpuallocator.NewReplicaPolicy().Allocate(gpuallocator.PolicyArgs(arg))

	} else {
		availableDevices, err := p.devSet.Filter(available)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve list of available devices: %v", err)
		}

		requiredDevices, err := p.devSet.Filter(required)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve list of required devices: %v", err)
		}

		arg := gpuallocator.BestPolicyArgs{Available: availableDevices, Required: requiredDevices, Size: size}

		devices = gpuallocator.NewBestEffortPolicy().Allocate(gpuallocator.PolicyArgs(arg))

	}
	return devices, nil
}

// Allocate returns list of devices.
func (p *iluvatarDevicePlugin) Allocate(ctx context.Context, reqs *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	responses := &pluginapi.AllocateResponse{}
	klog.Infof("Allocate request: %v", reqs)

	// reset before bind devices
	uuidResetMap := make(map[string]bool)
	for _, req := range reqs.ContainerRequests {
		if p.kubeclient != nil {
			volcanoDevices, isVolcano := p.UseVolcano(req.DevicesIDs)
			if isVolcano {
				req.DevicesIDs = volcanoDevices
			}
		}
		for _, id := range req.DevicesIDs {
			if !p.devSet.DeviceExist(id) {
				return nil, fmt.Errorf("Invalid allocation request for '%s': unknown device: %s", ResourceName, id)
			}
			// generateIDS: get all chip UUIDs of the device
			dev := p.devSet.Devices[gpuallocator.Alias(id).Prefix()]
			if dev == nil {
				return nil, fmt.Errorf("Invalid allocation request for '%s': device not found: %s", ResourceName, id)
			}
			deviceIDs := dev.GenerateIDS()
			for _, deviceID := range deviceIDs {
				uuidResetMap[deviceID] = true
			}

		}
	}

	var uuidResetList []string
	for uuid := range uuidResetMap {
		uuidResetList = append(uuidResetList, uuid)
	}

	p.resetGpusAndDeviceSet(uuidResetList)

	// After GPU reset, DeviceSet is rebuilt and device UUIDs may have changed.
	// Re-validate that all requested devices still exist in the new DeviceSet.
	for _, req := range reqs.ContainerRequests {
		for _, id := range req.DevicesIDs {
			if !p.devSet.DeviceExist(id) {
				return nil, fmt.Errorf("device '%s' no longer exists after GPU reset (UUID may have changed), please retry", id)
			}
		}
	}

	// bind devices
	for _, req := range reqs.ContainerRequests {
		response := &pluginapi.ContainerAllocateResponse{}
		var deviceIDs []string
		var replicaIDs []string

		// if all of the device is allocated to device plugin, keep container /dev/iluvatar[devMinor] same order with host
		if p.devSet.Replicas == 0 && len(req.DevicesIDs) == len(p.devSet.Devices) {
			var devMinors []int
			for _, device := range p.devSet.Devices {
				for _, chip := range device.Chips {
					devMinors = append(devMinors, int(chip.Minor))
				}
			}
			sort.Ints(devMinors)

			// generate device spec list by minor numbers
			for i, minor := range devMinors {
				klog.Infof("minor: %d, index: %d", minor, i)
				klog.Infof("HostPath: %s, ContainerPath: %s", config.HostPathPrefix+config.DeviceName+strconv.Itoa(minor), config.ContainerPathPrefix+config.DeviceName+strconv.Itoa(i))
				d := pluginapi.DeviceSpec{}
				d.HostPath = config.HostPathPrefix + config.DeviceName + strconv.Itoa(minor)
				d.ContainerPath = config.ContainerPathPrefix + config.DeviceName + strconv.Itoa(i) // start from 0 in container
				d.Permissions = "rw"
				response.Devices = append(response.Devices, &d)
			}

			deviceIDs = req.DevicesIDs
			replicaIDs = req.DevicesIDs
		} else {
			deviceSpecList := make(map[string]bool)
			for _, id := range req.DevicesIDs {
				deviceID := gpuallocator.Alias(id).Prefix()
				if _, ok := deviceSpecList[deviceID]; !ok {
					deviceSpecList[deviceID] = true
					dev := p.devSet.Devices[deviceID]
					if dev == nil {
						return nil, fmt.Errorf("device '%s' not found in DeviceSet after GPU reset", deviceID)
					}
					response.Devices = append(response.Devices, dev.GenerateSpecList()...)
					deviceIDs = append(deviceIDs, dev.GenerateIDS()...)
				}
				replicaIDs = append(replicaIDs, id)
			}
		}
		response.Devices = append(response.Devices, p.allocateCommonDeviceSpecs()...)
		response.Envs = p.allocateEnvs("IX_VISIBLE_DEVICES", deviceIDs)
		response.Envs["IX_REPLICA_DEVICES"] = strings.Join(replicaIDs, ",")

		responses.ContainerResponses = append(responses.ContainerResponses, response)
	}

	klog.Infof("Allocate response: %v", responses)
	return responses, nil
}

func (p *iluvatarDevicePlugin) allocateEnvs(envvar string, devices []string) map[string]string {
	return map[string]string{
		envvar: strings.Join(devices, ","),
	}
}

func (p *iluvatarDevicePlugin) allocateCommonDeviceSpecs() []*pluginapi.DeviceSpec {
	commonDevices := []string{
		"/dev/itrctl",
	}

	var specs []*pluginapi.DeviceSpec
	for _, dev := range commonDevices {
		if _, err := os.Stat(dev); err != nil {
			klog.Warningf("Control device %s not found on host, skipping mount", dev)
			continue
		}
		spec := &pluginapi.DeviceSpec{
			ContainerPath: dev,
			HostPath:      dev,
			Permissions:   "rw",
		}
		specs = append(specs, spec)
	}

	return specs
}

func (p *iluvatarDevicePlugin) allocateMountsByDeviceID(deviceID string) *pluginapi.Mount {
	var mount pluginapi.Mount

	for _, dev := range p.devSet.Devices {
		if deviceID == dev.UUID {
			// Mount for iluvatar pod
		}

	}

	return &mount
}

func (p *iluvatarDevicePlugin) allocateDevicesByDeviceID(hostminor uint, num int) *pluginapi.DeviceSpec {
	var device pluginapi.DeviceSpec

	hostPathPrefix := "/dev/"
	containerPathPrefix := "/dev/"

	// Expose the device node for iluvatar pod.
	device.HostPath = hostPathPrefix + deviceName + strconv.Itoa(int(hostminor))
	device.ContainerPath = containerPathPrefix + deviceName + strconv.Itoa(num)
	device.Permissions = "rw"

	return &device
}
