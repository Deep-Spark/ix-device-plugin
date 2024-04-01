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
	"strconv"
	"strings"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

// resourceName is the name to identify iluvatar device plugin
const resourceName string = "iluvatar.ai/gpu"

// iluvatarDevicePlugin is the implementation of iluvatar device plugin
type iluvatarDevicePlugin struct {
	iluvatarDevice

	name string

	stopList chan struct{}

	deviceCh chan pluginapi.Device
}

// GetDevicePluginOptions returns the values of the optional settings for this plugin
func (p *iluvatarDevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{}, nil
}

// PreStartContainer allows kubelet to pass reinitialized devices to containers.
func (p *iluvatarDevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	glog.Info("PreStartContainer...")
	return &pluginapi.PreStartContainerResponse{}, nil
}

// ListAndWatch lists devices
func (p *iluvatarDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	devs := p.cachedDevices()

	glog.Info("Start to list and watch GPU.")

	s.Send(&pluginapi.ListAndWatchResponse{Devices: devs})

	for {
		select {
		case <-p.stopList:
			glog.Info("Stoping list and watch GPU.")

			return nil
		case dev := <-p.deviceCh:
			for _, d := range devs {
				if dev.ID == d.ID {
					d.Health = pluginapi.Unhealthy
					glog.Infof("'%s' device marked unhealthy: %s", p.name, dev.ID)
				}
			}

			s.Send(&pluginapi.ListAndWatchResponse{Devices: devs})
		}
	}
}

// GetPreferredAllocation returns the preferred allocation from the set of devices specified in the request
func (p *iluvatarDevicePlugin) GetPreferredAllocation(ctx context.Context, r *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	response := &pluginapi.PreferredAllocationResponse{}

	return response, nil
}

// Allocate returns list of devices.
func (p *iluvatarDevicePlugin) Allocate(ctx context.Context, reqs *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	responses := &pluginapi.AllocateResponse{}
	response := &pluginapi.ContainerAllocateResponse{}

	glog.Infof("Allocate request: %v", reqs)

	for _, req := range reqs.ContainerRequests {
		var deviceIDs []string
		minorInContainer := 0
		for _, id := range req.DevicesIDs {
			if !p.deviceExist(id) {
				return nil, fmt.Errorf("Invalid allocation request for '%s': unknown device: %s", resourceName, id)
			}

			for _, dev := range p.devices {
				if id == dev.uuid {
					for k, v := range dev.minorslice {
						deviceIDs = append(deviceIDs, dev.uuidslice[k])
						device := p.allocateDevicesByDeviceID(v, minorInContainer)
						minorInContainer++
						response.Devices = append(response.Devices, device)
					}
				}
			}

		}
		response.Envs = p.allocateEnvs("ILUVATAR_COREX_VISIBLE_DEVICES", deviceIDs)
		responses.ContainerResponses = append(responses.ContainerResponses, response)
	}

	glog.Infof("Allocate response: %v", responses)

	return responses, nil
}

func (p *iluvatarDevicePlugin) allocateEnvs(envvar string, devices []string) map[string]string {
	return map[string]string{
		envvar: strings.Join(devices, ","),
	}
}

func (p *iluvatarDevicePlugin) allocateMountsByDeviceID(deviceID string) *pluginapi.Mount {
	var mount pluginapi.Mount

	for _, dev := range p.devices {
		if deviceID == dev.ID {
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
