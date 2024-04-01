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

package ixml

import (
	"fmt"
	"sync"
)

// #cgo LDFLAGS: -ldl
// #include "ixml.h"
import "C"

const (
	szDriverVersion = C.NVML_SYSTEM_DRIVER_VERSION_BUFFER_SIZE
	szCudaVersion   = 256
	szName          = C.NVML_DEVICE_NAME_BUFFER_SIZE
	szUUID          = C.NVML_DEVICE_UUID_BUFFER_SIZE
)

type deviceHandle C.nvmlDevice_t

var cachedDevicesByIndex map[uint]deviceHandle
var cachedDevicesByUUID map[string]deviceHandle

func init() {
	var once sync.Once

	once.Do(func() {
		cachedDevicesByIndex = make(map[uint]deviceHandle)
		cachedDevicesByUUID = make(map[string]deviceHandle)
	})
}

func registerDevice(index uint, uuid string, device deviceHandle) {
	cachedDevicesByIndex[index] = device
	cachedDevicesByUUID[uuid] = device
}

func unregisterDevice() {
	for k := range cachedDevicesByIndex {
		delete(cachedDevicesByIndex, k)
	}

	for k := range cachedDevicesByUUID {
		delete(cachedDevicesByUUID, k)
	}
}

func deviceInit() error {
	ret := C.dl_init()
	if ret == C.NVML_ERROR_LIBRARY_NOT_FOUND {
		return fmt.Errorf("Library '%s' not found.", string(C.IXML_LIBRARY))
	} else if ret == C.NVML_ERROR_FUNCTION_NOT_FOUND {
		return fmt.Errorf("Symbol not found.")
	}

	ret = C.ixmlInit()
	if ret != C.NVML_SUCCESS {
		return fmt.Errorf("Failed to initialize ixml.")
	}

	return nil
}

func deviceShutdown() error {
	ret := C.ixmlShutdown()
	if ret != C.NVML_SUCCESS {
		return fmt.Errorf("Failed to shutdown ixml.")
	}

	ret = C.dl_close()
	if ret != C.NVML_SUCCESS {
		return fmt.Errorf("Failed to close handler of '%s'.", string(C.IXML_LIBRARY))
	}

	return nil
}

func getDeviceCount() (uint, error) {
	var num C.uint

	ret := C.ixmlDeviceGetCount(&num)
	if ret != C.NVML_SUCCESS {
		return 0, fmt.Errorf("Failed to get the count of gpu device.")
	}

	return uint(num), nil
}

func getDriverVersion() (string, error) {
	var version [szDriverVersion]C.char

	ret := C.ixmlSystemGetDriverVersion(&version[0], szDriverVersion)
	if ret != C.NVML_SUCCESS {
		return "", fmt.Errorf("Failed to get the driver version of gpu device.")
	}

	return C.GoString(&version[0]), nil
}

func getCudaVersion() (string, error) {
	var version C.int

	ret := C.ixmlSystemGetCudaDriverVersion(&version)
	if ret != C.NVML_SUCCESS {
		return "", fmt.Errorf("Failed to get the current CUDA version.")
	}

	major := uint(version / 1000)
	minor := uint(version % 1000 / 10)

	return fmt.Sprintf("%d", major) + "." + fmt.Sprintf("%d", minor), nil
}

func getDeviceByIndex(index uint) (*device, error) {
	if cachedDev, ok := cachedDevicesByIndex[index]; ok {
		return &device{handle: cachedDev}, nil
	}

	var dev C.nvmlDevice_t

	ret := C.ixmlDeviceGetHandleByIndex(C.uint(index), &dev)
	if ret != C.NVML_SUCCESS {
		return nil, fmt.Errorf("Failed to get device handle of gpu-%d", index)
	}

	d := &device{handle: deviceHandle(dev)}
	uuid, err := d.DeviceGetUUID()
	if err != nil {
		return nil, err
	}

	registerDevice(index, uuid, d.handle)

	return d, nil
}

func getDeviceByUUID(uuid string) (*device, error) {
	if cachedDev, ok := cachedDevicesByUUID[uuid]; ok {
		return &device{handle: cachedDev}, nil
	}

	var dev C.nvmlDevice_t

	ret := C.ixmlDeviceGetHandleByUUID(C.CString(uuid), &dev)
	if ret != C.NVML_SUCCESS {
		return nil, fmt.Errorf("Failed to get device handle of gpu-%s", uuid)
	}

	d := &device{handle: deviceHandle(dev)}
	index, err := d.DeviceGetIndex()
	if err != nil {
		return nil, err
	}

	registerDevice(index, uuid, d.handle)

	return d, nil
}

func getDeviceOnSameBoard(device1 device, device2 device, onSameBoard *C.int) error {
	ret := C.ixmlDeviceOnSameBoard(device1.handle, device2.handle, onSameBoard)

	if ret != C.NVML_SUCCESS {
		return fmt.Errorf("Failed to judge whether two devices on same board")
	}

	return nil
}

type device struct {
	handle deviceHandle
}

func (d *device) DeviceGetName() (string, error) {
	var name [256]C.char

	ret := C.ixmlDeviceGetName(d.handle, &name[0], 256)
	if ret != C.NVML_SUCCESS {
		return "", fmt.Errorf("Failed to get device name of gpu")
	}

	return C.GoString(&name[0]), nil
}

func (d *device) DeviceGetUUID() (string, error) {
	var uuid [szUUID]C.char

	ret := C.ixmlDeviceGetUUID(d.handle, &uuid[0], szUUID)
	if ret != C.NVML_SUCCESS {
		return "", fmt.Errorf("Failed to get device UUID of gpu")
	}

	return C.GoString(&uuid[0]), nil
}

func (d *device) DeviceGetUUIDSlice() ([]string, error) {
	uuid, err := d.DeviceGetUUID()
	if err != nil {
		return nil, fmt.Errorf("Failed to get device UUID of gpu")
	}

	minor, err := d.DeviceGetMinorNumber()
	if err != nil {
		return nil, fmt.Errorf("Failed to get device minor number : %v", err)
	}

	var uuidslice []string
	count, _ := getDeviceCount()
	for i := uint(0); i < count; i++ {
		if i == minor {
			uuidslice = append(uuidslice, uuid)
			continue
		}

		device, err := getDeviceByIndex(i)
		if err != nil {
			return nil, fmt.Errorf("Failed to get device handle of gpu-%d", i)
		}

		var onSameBoard C.int
		getDeviceOnSameBoard(*d, *device, &onSameBoard)
		if onSameBoard == 1 {
			uuid, err := device.DeviceGetUUID()
			if err != nil {
				return nil, fmt.Errorf("Failed to get device UUID of gpu")
			}
			uuidslice = append(uuidslice, uuid)
		}

	}

	fmt.Println()

	return uuidslice, nil
}

func (d *device) DeviceGetIndex() (uint, error) {
	var index C.uint

	ret := C.ixmlDeviceGetIndex(d.handle, &index)
	if ret != C.NVML_SUCCESS {
		return 0, fmt.Errorf("Failed to get device index of gpu")
	}

	return uint(index), nil
}

func (d *device) DeviceGetMinorNumber() (uint, error) {
	var minor C.uint

	ret := C.ixmlDeviceGetMinorNumber(d.handle, &minor)
	if ret != C.NVML_SUCCESS {
		// FIXME: 100 is a pseudo value.
		return 100, fmt.Errorf("Failed to get device minor number of gpu")
	}

	return uint(minor), nil
}

func (d *device) DeviceGetMinorSlice() ([]uint, error) {
	minor, err := d.DeviceGetMinorNumber()
	if err != nil {
		return nil, fmt.Errorf("Failed to get device minor number : %v", err)
	}

	var minorslice []uint
	count, _ := getDeviceCount()
	for i := uint(0); i < count; i++ {
		if i == minor {
			minorslice = append(minorslice, minor)
			continue
		}

		var dev C.nvmlDevice_t
		ret := C.ixmlDeviceGetHandleByIndex(C.uint(i), &dev)
		if ret != C.NVML_SUCCESS {
			return nil, fmt.Errorf("Failed to get device handle of gpu-%d", i)
		}
		device := device{handle: deviceHandle(dev)}

		var onSameBoard C.int
		getDeviceOnSameBoard(*d, device, &onSameBoard)
		if onSameBoard == 1 {
			minorslice = append(minorslice, i)
		}
	}

	return minorslice, nil
}

func (d *device) DeviceGetFanSpeed() (uint, error) {
	var speed C.uint

	ret := C.ixmlDeviceGetFanSpeed(d.handle, &speed)
	if ret != C.NVML_SUCCESS {
		return 0, fmt.Errorf("Failed to get fan speed of gpu")
	}
	return uint(speed), nil
}

func (d *device) DeviceGetMemoryInfo() (MemoryInfo, error) {
	var mem C.nvmlMemory_t

	ret := C.ixmlDeviceGetMemoryInfo(d.handle, &mem)
	if ret != C.NVML_SUCCESS {
		return MemoryInfo{}, fmt.Errorf("Failed to get memory information of gpu")
	}

	totalMem := uint64(mem.total)
	usedMem := uint64(mem.used)
	freeMem := uint64(mem.free)

	// convert 'Byte' to 'MiB'
	return MemoryInfo{
		Total: totalMem / 1024 / 1024,
		Used:  usedMem / 1024 / 1024,
		Free:  freeMem / 1024 / 1024,
	}, nil
}

func (d *device) DeviceGetTemperature() (uint, error) {
	var temp C.uint

	ret := C.ixmlDeviceGetTemperature(d.handle, C.NVML_TEMPERATURE_GPU, &temp)
	if ret != C.NVML_SUCCESS {
		return 0, fmt.Errorf("Failed to get the current temperature of gpu")
	}

	return uint(temp), nil
}

func (d *device) DeviceGetPciInfo() (PciInfo, error) {
	var pci C.nvmlPciInfo_t

	ret := C.ixmlDeviceGetPciInfo(d.handle, &pci)
	if ret != C.NVML_SUCCESS {
		return PciInfo{}, fmt.Errorf("Failed to get pci information of gpu")
	}

	return PciInfo{
		Bus:            uint(pci.bus),
		BusId:          C.GoString(&pci.busId[0]),
		BusIdLegacy:    C.GoString(&pci.busIdLegacy[0]),
		Device:         uint(pci.device),
		Domain:         uint(pci.domain),
		PciDeviceId:    uint(pci.pciDeviceId),
		PciSubSystemId: uint(pci.pciSubSystemId),
	}, nil
}

func (d *device) DeviceGetPowerUsage() (uint, error) {
	var usage C.uint

	ret := C.ixmlDeviceGetPowerUsage(d.handle, &usage)
	if ret != C.NVML_SUCCESS {
		return 0, fmt.Errorf("Failed to get power usage of gpu")
	}

	return uint(usage), nil
}

func (d *device) DeviceGetPowerLimitConstraints() (PowerLimitConstraints, error) {
	var max, min C.uint

	ret := C.ixmlDeviceGetPowerManagementLimitConstraints(d.handle, &max, &min)
	if ret != C.NVML_SUCCESS {
		return PowerLimitConstraints{}, fmt.Errorf("Failed to get power limitation of gpu")
	}

	return PowerLimitConstraints{
		MaxLimit: uint(max),
		MinLimit: uint(min),
	}, nil
}

func (d *device) DeviceGetClockInfo() (ClockInfo, error) {
	var sm, mem C.uint

	ret := C.ixmlDeviceGetClockInfo(d.handle, C.NVML_CLOCK_SM, &sm)
	if ret != C.NVML_SUCCESS {
		return ClockInfo{}, fmt.Errorf("Failed to get SM clock of gpu")
	}

	ret = C.ixmlDeviceGetClockInfo(d.handle, C.NVML_CLOCK_MEM, &mem)
	if ret != C.NVML_SUCCESS {
		return ClockInfo{}, fmt.Errorf("Failed to get MEM clock of gpu")
	}

	return ClockInfo{
		Sm:  uint(sm),
		Mem: uint(mem),
	}, nil
}

func (d *device) DeviceGetUtilization() (Utilization, error) {
	var utilization C.nvmlUtilization_t

	ret := C.ixmlDeviceGetUtilizationRates(d.handle, &utilization)
	if ret != C.NVML_SUCCESS {
		return Utilization{}, fmt.Errorf("Failed to get utilization rates of gpu")
	}

	return Utilization{
		GPU: uint(utilization.gpu),
		Mem: uint(utilization.memory),
	}, nil
}
