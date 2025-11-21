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
	response := &pluginapi.ContainerAllocateResponse{}

	klog.Infof("Allocate request: %v", reqs)

	var deviceIDs []string
	var replicaIDs []string
	var indexes []int
	for _, req := range reqs.ContainerRequests {

		if p.kubeclient != nil {
			volcanoDevices, isVolcano := p.UseVolcano(req.DevicesIDs)
			if isVolcano {
				req.DevicesIDs = volcanoDevices
			}
		}

		DeviceSpecList := make(map[string]bool)

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
				indexes = append(indexes, i)
			}

			deviceIDs = req.DevicesIDs
			replicaIDs = req.DevicesIDs
		} else {
			for _, id := range req.DevicesIDs {
				if !p.devSet.DeviceExist(id) {
					return nil, fmt.Errorf("Invalid allocation request for '%s': unknown device: %s", ResourceName, id)
				}
				prefix := gpuallocator.Alias(id).Prefix()
				if _, ok := DeviceSpecList[prefix]; !ok {
					DeviceSpecList[prefix] = true
					response.Devices = append(response.Devices, p.devSet.Devices[prefix].GenerateSpecList()...)
					deviceIDs = append(deviceIDs, p.devSet.Devices[prefix].GenerateIDS()...)
					indexes = append(indexes, p.devSet.Devices[prefix].GenerateIndexs()...)
				}
				replicaIDs = append(replicaIDs, id)
			}
		}
		responses.ContainerResponses = append(responses.ContainerResponses, response)
	}

	response.Envs = p.allocateEnvs("IX_VISIBLE_DEVICES", deviceIDs)
	response.Envs["IX_REPLICA_DEVICES"] = strings.Join(replicaIDs, ",")

	klog.Infof("Allocate response: %v", responses)

	p.resetGpusAndDeviceSet(indexes)
	return responses, nil
}

func (p *iluvatarDevicePlugin) allocateEnvs(envvar string, devices []string) map[string]string {
	return map[string]string{
		envvar: strings.Join(devices, ","),
	}
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
