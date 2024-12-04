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
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"

	goixml "gitee.com/deep-spark/go-ixml/pkg/ixml"
)

var (
	HealthSYSHUBError      = fmt.Errorf("SYSHUBError")
	HealthMCError          = fmt.Errorf("MCError")
	HealthOverTempError    = fmt.Errorf("OverTempError")
	HealthOverVoltageError = fmt.Errorf("OverVoltageError")
	HealthECCError         = fmt.Errorf("ECCError")
	HealthMemoryError      = fmt.Errorf("MemoryError")
	HealthPCIEError        = fmt.Errorf("PCIEError")
)

func deviceInit() error {
	ret := goixml.Init()
	if ret != goixml.SUCCESS {
		return fmt.Errorf("Failed to init ixml.")
	}
	return nil
}

func deviceShutdown() error {
	ret := goixml.Shutdown()
	if ret != goixml.SUCCESS {
		return fmt.Errorf("Failed to shutdown ixml.")
	}

	return nil
}

func getDeviceCount() (uint, error) {
	num, ret := goixml.DeviceGetCount()
	if ret != goixml.SUCCESS {
		return 0, fmt.Errorf("Failed to get the count of gpu device.")
	}

	return uint(num), nil
}

func getDriverVersion() (string, error) {
	version, ret := goixml.SystemGetDriverVersion()
	if ret != goixml.SUCCESS {
		return "", fmt.Errorf("Failed to get the driver version of gpu device.")
	}

	return version, nil
}

func getCudaVersion() (string, error) {
	cudaversion, ret := goixml.SystemGetCudaDriverVersion()
	if ret != goixml.SUCCESS {
		return "", fmt.Errorf("Failed to get the current CUDA version.")
	}

	version, err := strconv.Atoi(cudaversion)
	if err != nil {
		return "", fmt.Errorf("Failed to get the current CUDA version.")
	}
	major := uint(version / 1000)
	minor := uint(version % 1000 / 10)

	return fmt.Sprintf("%d", major) + "." + fmt.Sprintf("%d", minor), nil
}

func getDeviceByIndex(index uint) (*device, error) {
	var dev goixml.Device
	ret := goixml.DeviceGetHandleByIndex(index, &dev)
	if ret != goixml.SUCCESS {
		return nil, fmt.Errorf("Failed to get device handle of gpu-%d", index)
	}

	d := &device{Device: dev}

	return d, nil
}

func getDeviceByUUID(uuid string) (*device, error) {
	dev, ret := goixml.GetHandleByUUID(uuid)
	if ret != goixml.SUCCESS {
		return nil, fmt.Errorf("Failed to get device handle of gpu-%s", uuid)
	}

	d := &device{Device: dev}

	return d, nil
}

func GetDeviceOnSameBoard(device1 Device, device2 Device) (error, bool) {
	isOnSameBoard := false
	dev1, ok := device1.(*device)
	if ok != true {
		return fmt.Errorf("Type Error"), isOnSameBoard
	}
	dev2, ok := device2.(*device)
	if ok != true {
		return fmt.Errorf("Type Error"), isOnSameBoard
	}

	onSameBoard, ret := goixml.GetOnSameBoard(dev1.Device, dev2.Device)
	if ret != goixml.SUCCESS {
		return fmt.Errorf("Failed to judge whether two devices on same board"), isOnSameBoard
	}

	if onSameBoard == 0 {
		isOnSameBoard = false
	} else {
		isOnSameBoard = true
	}

	return nil, isOnSameBoard
}

type device struct {
	goixml.Device
}

func (d *device) DeviceGetName() (string, error) {
	name, ret := d.GetName()
	if ret != goixml.SUCCESS {
		return "", fmt.Errorf("Failed to get device name of gpu")
	}

	return name, nil
}

func (d *device) DeviceGetUUID() (string, error) {
	uuid, ret := d.GetUUID()
	if ret != goixml.SUCCESS {
		return "", fmt.Errorf("Failed to get device UUID of gpu")
	}

	return uuid, nil
}

func (d *device) DeviceGetIndex() (uint, error) {
	index, ret := d.GetIndex()
	if ret != goixml.SUCCESS {
		return 0, fmt.Errorf("Failed to get device index of gpu")
	}

	return uint(index), nil
}

func (d *device) DeviceGetMinorNumber() (uint, error) {
	minor, ret := d.GetMinorNumber()
	if ret != goixml.SUCCESS {
		// FIXME: 100 is a pseudo value.
		return 100, fmt.Errorf("Failed to get device minor number of gpu")
	}

	return uint(minor), nil
}

func (d *device) DeviceGetFanSpeed() (uint, error) {
	speed, ret := d.GetFanSpeed()
	if ret != goixml.SUCCESS {
		return 0, fmt.Errorf("Failed to get fan speed of gpu")
	}
	return uint(speed), nil
}

func (d *device) DeviceGetMemoryInfo() (MemoryInfo, error) {
	mem, ret := d.GetMemoryInfo()
	if ret != goixml.SUCCESS {
		return MemoryInfo{}, fmt.Errorf("Failed to get memory information of gpu")
	}

	totalMem := uint64(mem.Total)
	usedMem := uint64(mem.Used)
	freeMem := uint64(mem.Free)

	// convert 'Byte' to 'MiB'
	return MemoryInfo{
		Total: totalMem / 1024 / 1024,
		Used:  usedMem / 1024 / 1024,
		Free:  freeMem / 1024 / 1024,
	}, nil
}

func (d *device) DeviceGetTemperature() (uint, error) {
	temp, ret := d.GetTemperature()
	if ret != goixml.SUCCESS {
		return 0, fmt.Errorf("Failed to get the current temperature of gpu")
	}

	return uint(temp), nil
}

func (d *device) DeviceGetPciInfo() (PciInfo, error) {
	pci, ret := d.GetPciInfo()
	if ret != goixml.SUCCESS {
		return PciInfo{}, fmt.Errorf("Failed to get pci information of gpu")
	}

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, &pci.BusIdLegacy)
	busidlegacy := bytesBuffer.String()

	bytesBuffer = bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, &pci.BusId)
	busid := bytesBuffer.String()

	return PciInfo{
		Bus:            uint(pci.Bus),
		BusId:          busid,
		BusIdLegacy:    busidlegacy,
		Device:         uint(pci.Device),
		Domain:         uint(pci.Domain),
		PciDeviceId:    uint(pci.PciDeviceId),
		PciSubSystemId: uint(pci.PciSubSystemId),
	}, nil
}

func (d *device) DeviceGetPowerUsage() (uint, error) {
	usage, ret := d.GetPowerUsage()
	if ret != goixml.SUCCESS {
		return 0, fmt.Errorf("Failed to get power usage of gpu")
	}

	return uint(usage), nil
}

func (d *device) DeviceGetPowerLimitConstraints() (PowerLimitConstraints, error) {
	min, max, ret := d.GetPowerManagementLimitConstraints()
	if ret != goixml.SUCCESS {
		return PowerLimitConstraints{}, fmt.Errorf("Failed to get power limitation of gpu")
	}

	return PowerLimitConstraints{
		MaxLimit: uint(max),
		MinLimit: uint(min),
	}, nil
}

func (d *device) DeviceGetClockInfo() (ClockInfo, error) {
	clockinfo, ret := d.GetClockInfo()
	if ret != goixml.SUCCESS {
		return ClockInfo{}, fmt.Errorf("Failed to get SM clock of gpu")
	}

	return ClockInfo{
		Sm:  uint(clockinfo.Sm),
		Mem: uint(clockinfo.Mem),
	}, nil
}

func (d *device) DeviceGetUtilization() (Utilization, error) {
	utilization, ret := d.GetUtilizationRates()
	if ret != goixml.SUCCESS {
		return Utilization{}, fmt.Errorf("Failed to get utilization rates of gpu")
	}

	return Utilization{
		GPU: uint(utilization.Gpu),
		Mem: uint(utilization.Memory),
	}, nil
}

func CheckDeviceError(health Health) []error {
	errs := []error{}
	if (health & Health(goixml.HealthSYSHUBError)) > 0 {
		errs = append(errs, HealthSYSHUBError)
	}
	if (health & Health(goixml.HealthMCError)) > 0 {
		errs = append(errs, HealthMCError)
	}
	if (health & Health(goixml.HealthOverTempError)) > 0 {
		errs = append(errs, HealthOverTempError)
	}
	if (health & Health(goixml.HealthOverVoltageError)) > 0 {
		errs = append(errs, HealthOverVoltageError)
	}
	if (health & Health(goixml.HealthECCError)) > 0 {
		errs = append(errs, HealthECCError)
	}
	if (health & Health(goixml.HealthMemoryError)) > 0 {
		errs = append(errs, HealthMemoryError)
	}
	if (health & Health(goixml.HealthPCIEError)) > 0 {
		errs = append(errs, HealthPCIEError)
	}

	return errs
}

func (d *device) DeviceGetHealth() (Health, error) {
	health, ret := d.GetHealth()
	if ret != goixml.SUCCESS {
		return Health(health), fmt.Errorf("Failed to get Health status of GPU: %v", ret)
	}

	return Health(health), nil
}

func (d *device) DeviceGetNumaNode() (bool, int, error) {
	info, err := d.DeviceGetPciInfo()
	if err != nil {
		return false, 0, fmt.Errorf("error getting PCI Bus Info of device: %v", err)
	}

	busID := strings.ToLower(info.BusIdLegacy)
	b, err := os.ReadFile(fmt.Sprintf("/sys/bus/pci/devices/%s/numa_node", busID))
	if err != nil {
		return false, 0, nil
	}

	node, err := strconv.Atoi(string(bytes.TrimSpace(b)))
	if err != nil {
		return false, 0, fmt.Errorf("eror parsing value for NUMA node: %v", err)
	}

	if node < 0 {
		return false, 0, nil
	}

	return true, node, nil
}

func (d *device) DeviceGetTopology(device2 *Device) (goixml.GpuTopologyLevel, error) {
	dev2, ok := (*device2).(*device)
	if ok != true {
		return goixml.GpuTopologyLevel(0), fmt.Errorf("unkown topology")
	}

	pathinfo, ret := d.GetTopology(dev2.Device)
	if ret != goixml.SUCCESS {
		return goixml.GpuTopologyLevel(0), fmt.Errorf("unkown topology %v", ret)
	} else {
		return pathinfo, nil
	}
}

func (d *device) DeviceGetBoardPosition() (bool, int) {
	pos, ret := d.GetBoardPosition()
	if ret == goixml.SUCCESS {
		return true, int(pos)
	} else {
		return false, 0
	}
}

func (d *device) GetSelf() *device {
	return d
}
