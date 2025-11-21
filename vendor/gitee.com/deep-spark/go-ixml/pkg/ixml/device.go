/*
Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ixml

// ixml.DeviceGetCount()
func DeviceGetCount() (uint, Return) {
	var DeviceCount uint32
	ret := nvmlDeviceGetCount(&DeviceCount)
	return uint(DeviceCount), ret
}

// ixml.DeviceGetHandleByIndex()
func DeviceGetHandleByIndex(index uint, device *Device) Return {
	ret := nvmlDeviceGetHandleByIndex(uint32(index), device)
	return ret
}

// ixml.DeviceGetHandleBySerial()
func DeviceGetHandleBySerial(serial string) (Device, Return) {
	var device Device
	ret := nvmlDeviceGetHandleBySerial(serial+string(rune(0)), &device)
	return device, ret
}

// ixml.DeviceGetHandleByUUID()
func GetHandleByUUID(Uuid string) (Device, Return) {
	var Device Device
	ret := nvmlDeviceGetHandleByUUID(Uuid, &Device)
	return Device, ret
}

// ixml.DeviceGetUUID()
func DeviceGetUUID(device Device) (string, Return) {
	return device.GetUUID()
}

func (device Device) GetUUID() (string, Return) {
	uuid := make([]byte, DEVICE_UUID_BUFFER_SIZE)
	ret := nvmlDeviceGetUUID(device, &uuid[0], DEVICE_UUID_BUFFER_SIZE)
	return removeBytesSpaces(uuid), ret
}

// ixml.DeviceGetHandleByPciBusId()
// The format of pciBusId is "domain:bus:device.function", e.g., "00000000:1F:00.0".
func DeviceGetHandleByPciBusId(pciBusId string) (Device, Return) {
	var device Device
	ret := nvmlDeviceGetHandleByPciBusId_v2(pciBusId+string(rune(0)), &device)
	return device, ret
}

// ixml.DeviceGetMinorNumber()
func DeviceGetMinorNumber(device Device) (int, Return) {
	return device.GetMinorNumber()
}

func (device Device) GetMinorNumber() (int, Return) {
	var minorNumber uint32
	ret := nvmlDeviceGetMinorNumber(device, &minorNumber)
	return int(minorNumber), ret
}

// ixml.DeviceGetName()
func DeviceGetName(device Device) (string, Return) {
	return device.GetName()
}

func (device Device) GetName() (string, Return) {
	name := make([]byte, DEVICE_NAME_BUFFER_SIZE)
	ret := nvmlDeviceGetName(device, &name[0], DEVICE_NAME_BUFFER_SIZE)
	removeBytesSpaces(name)
	return removeBytesSpaces(name), ret
}

// ixml.DeviceGetTemperature()
func DeviceGetTemperature(device Device) (uint32, Return) {
	return device.GetTemperature()
}

func (device Device) GetTemperature() (uint32, Return) {
	var sensorType TemperatureSensors
	var temp uint32
	ret := nvmlDeviceGetTemperature(device, sensorType, &temp)
	return temp, ret
}

func (device Device) GetFanSpeed() (uint32, Return) {
	return deviceGetFanSpeed(device)
}

func deviceGetFanSpeed(Device Device) (uint32, Return) {
	var Speed uint32
	ret := nvmlDeviceGetFanSpeed(Device, &Speed)
	return Speed, ret
}

func (device Device) GetFanSpeed_v2(Fan int) (uint32, Return) {
	return deviceGetFanSpeed_v2(device, Fan)
}

func deviceGetFanSpeed_v2(Device Device, Fan int) (uint32, Return) {
	var Speed uint32
	ret := nvmlDeviceGetFanSpeed_v2(Device, uint32(Fan), &Speed)
	return Speed, ret
}

type ClockInfo struct {
	Sm  uint32
	Mem uint32
}

func (device Device) GetClockInfo() (ClockInfo, Return) {
	return deviceGetClockInfo(device)
}

func deviceGetClockInfo(Device Device) (ClockInfo, Return) {
	var sm, mem uint32
	_ = nvmlDeviceGetClockInfo(Device, CLOCK_SM, &sm)
	ret := nvmlDeviceGetClockInfo(Device, CLOCK_MEM, &mem)
	return ClockInfo{
		Sm:  sm,
		Mem: mem,
	}, ret
}

func (device Device) GetMemoryInfo() (Memory, Return) {
	mem, ret := deviceGetMemoryInfo(device)
	return Memory{
		Total: mem.Total / 1024 / 1024,
		Free:  mem.Free / 1024 / 1024,
		Used:  mem.Used / 1024 / 1024, //to MB
	}, ret
}

func deviceGetMemoryInfo(Device Device) (Memory, Return) {
	var Memory Memory
	ret := nvmlDeviceGetMemoryInfo(Device, &Memory)
	return Memory, ret
}

func (device Device) GetMemoryInfo_v2() (Memory, Return) {
	return deviceGetMemoryInfo(device)
}

func deviceGetMemoryInfo_v2(Device Device) (Memory_v2, Return) {
	var Memory Memory_v2
	ret := nvmlDeviceGetMemoryInfo_v2(Device, &Memory)
	return Memory, ret
}

func (device Device) GetUtilizationRates() (Utilization, Return) {
	return deviceGetUtilizationRates(device)
}

func deviceGetUtilizationRates(Device Device) (Utilization, Return) {
	var Utilization Utilization
	ret := nvmlDeviceGetUtilizationRates(Device, &Utilization)
	return Utilization, ret
}

// ixml.DeviceGetComputeMode()
func DeviceGetComputeMode(device Device) (ComputeMode, Return) {
	return device.GetComputeMode()
}

func (device Device) GetComputeMode() (ComputeMode, Return) {
	var mode ComputeMode
	ret := nvmlDeviceGetComputeMode(device, &mode)
	return mode, ret
}

// ixml.DeviceGetCudaComputeCapability()
func DeviceGetCudaComputeCapability(device Device) (int, int, Return) {
	return device.GetCudaComputeCapability()
}

func (device Device) GetCudaComputeCapability() (int, int, Return) {
	var major, minor int32
	ret := nvmlDeviceGetCudaComputeCapability(device, &major, &minor)
	return int(major), int(minor), ret
}

// ixml.DeviceGetEccMode()
func DeviceGetEccMode(device Device) (EnableState, EnableState, Return) {
	return device.GetEccMode()
}

func (device Device) GetEccMode() (EnableState, EnableState, Return) {
	var current, pending EnableState
	ret := nvmlDeviceGetEccMode(device, &current, &pending)
	return current, pending, ret
}

// ixml.DeviceGetBoardId()
func DeviceGetBoardId(device Device) (uint32, Return) {
	return device.GetBoardId()
}

func (device Device) GetBoardId() (uint32, Return) {
	var boardId uint32
	ret := nvmlDeviceGetBoardId(device, &boardId)
	return boardId, ret
}

// ixml.DeviceGetPciInfo()
func DeviceGetPciInfo(device Device) (PciInfo, Return) {
	return device.GetPciInfo()
}

func (device Device) GetPciInfo() (PciInfo, Return) {
	var PciInfo PciInfo
	ret := nvmlDeviceGetPciInfo(device, &PciInfo)
	return PciInfo, ret
}

// ixml.DeviceGetIndex()
func DeviceGetIndex(device Device) (int, Return) {
	return device.GetIndex()
}

func (device Device) GetIndex() (int, Return) {
	var Index uint32
	ret := nvmlDeviceGetIndex(device, &Index)
	return int(Index), ret
}

// ixml.DeviceGetSerial()
func DeviceGetSerial(device Device) (string, Return) {
	return device.GetSerial()
}

func (device Device) GetSerial() (string, Return) {
	serial := make([]byte, DEVICE_SERIAL_BUFFER_SIZE)
	ret := nvmlDeviceGetSerial(device, &serial[0], DEVICE_SERIAL_BUFFER_SIZE)
	return string(serial[:clen(serial)]), ret
}

// ixml.DeviceGetPowerUsage()
func DeviceGetPowerUsage(device Device) (uint32, Return) {
	return device.GetPowerUsage()
}

func (device Device) GetPowerUsage() (uint32, Return) {
	var Power uint32
	ret := nvmlDeviceGetPowerUsage(device, &Power)
	return Power, ret
}

// ixml.DeviceGetOnSameBoard()
func GetOnSameBoard(device1, device2 Device) (int, Return) {
	var OnSameBoard int32
	ret := nvmlDeviceOnSameBoard(device1, device2, &OnSameBoard)
	return int(OnSameBoard), ret
}

// ixml.DeviceGetBoardPosition()
func DeviceGetBoardPosition(device Device) (uint32, Return) {
	return device.GetBoardPosition()
}

func (device Device) GetBoardPosition() (uint32, Return) {
	var pos uint32
	ret := ixmlDeviceGetBoardPosition(device, &pos)
	return pos, ret
}

// ixml.DeviceGetGPUVoltage()
func DeviceGetGPUVoltage(device Device) (uint32, uint32, Return) {
	return device.GetGPUVoltage()
}

func (device Device) GetGPUVoltage() (uint32, uint32, Return) {
	var integer, decimal uint32
	ret := ixmlDeviceGetGPUVoltage(device, &integer, &decimal)
	return integer, decimal, ret
}

type Info struct {
	Pid           uint32
	Name          string
	UsedGpuMemory uint64
}

func (device Device) GetComputeRunningProcesses() ([]Info, Return) {
	processInfos, ret := deviceGetComputeRunningProcesses(device)
	if ret != SUCCESS {
		return nil, ret
	}

	Infos := make([]Info, len(processInfos))
	for i, processInfo := range processInfos {
		Infos[i].Pid = processInfo.Pid
		Infos[i].Name = getPidName(processInfo.Pid)
		Infos[i].UsedGpuMemory = processInfo.UsedGpuMemory / 1024 / 1024
	}
	return Infos, ret
}

func deviceGetComputeRunningProcesses(device Device) ([]ProcessInfo_v1, Return) {
	var InfoCount uint32 = 1
	for {
		infos := make([]ProcessInfo_v1, InfoCount)
		ret := nvmlDeviceGetComputeRunningProcesses(device, &InfoCount, &infos[0])
		if ret == SUCCESS {
			return infos[:InfoCount], ret
		}
		if ret != ERROR_INSUFFICIENT_SIZE {
			return nil, ret
		}
		InfoCount *= 2
	}
}

// ixml.DeviceGetCurrentClocksThrottleReasons()
func DeviceGetCurrentClocksThrottleReasons(device Device) (uint64, Return) {
	return device.GetCurrentClocksThrottleReasons()
}

func (device Device) GetCurrentClocksThrottleReasons() (uint64, Return) {
	var clocksThrottleReasons uint64
	ret := nvmlDeviceGetCurrentClocksThrottleReasons(device, &clocksThrottleReasons)
	return clocksThrottleReasons, ret
}

// ixml.DeviceGetMaxPcieLinkGeneration()
func DeviceGetMaxPcieLinkGeneration(device Device) (int, Return) {
	return device.GetMaxPcieLinkGeneration()
}

func (device Device) GetMaxPcieLinkGeneration() (int, Return) {
	var maxLinkGen uint32
	ret := nvmlDeviceGetMaxPcieLinkGeneration(device, &maxLinkGen)
	return int(maxLinkGen), ret
}

// ixml.DeviceGetMaxPcieLinkWidth()
func DeviceGetMaxPcieLinkWidth(device Device) (int, Return) {
	return device.GetMaxPcieLinkWidth()
}

func (device Device) GetMaxPcieLinkWidth() (int, Return) {
	var maxLinkWidth uint32
	ret := nvmlDeviceGetMaxPcieLinkWidth(device, &maxLinkWidth)
	return int(maxLinkWidth), ret
}

// ixml.DeviceGetCurrPcieLinkGeneration()
func DeviceGetCurrPcieLinkGeneration(device Device) (int, Return) {
	return device.GetCurrPcieLinkGeneration()
}

func (device Device) GetCurrPcieLinkGeneration() (int, Return) {
	var currLinkGen uint32
	ret := nvmlDeviceGetCurrPcieLinkGeneration(device, &currLinkGen)
	return int(currLinkGen), ret
}

// ixml.DeviceGetCurrPcieLinkWidth()
func DeviceGetCurrPcieLinkWidth(device Device) (int, Return) {
	return device.GetCurrPcieLinkWidth()
}

func (device Device) GetCurrPcieLinkWidth() (int, Return) {
	var currLinkWidth uint32
	ret := nvmlDeviceGetCurrPcieLinkWidth(device, &currLinkWidth)
	return int(currLinkWidth), ret
}

// ixml.DeviceGetPcieThroughput()
func DeviceGetPcieThroughput(device Device, counter PcieUtilCounter) (uint32, Return) {
	return device.GetPcieThroughput(counter)
}

func (device Device) GetPcieThroughput(counter PcieUtilCounter) (uint32, Return) {
	var value uint32
	ret := nvmlDeviceGetPcieThroughput(device, counter, &value)
	return value, ret
}

// ixml.DeviceGetPcieReplayCounter()
func DeviceGetPcieReplayCounter(device Device) (uint32, Return) {
	return device.GetPcieReplayCounter()
}

func (device Device) GetPcieReplayCounter() (uint32, Return) {
	var value uint32
	ret := nvmlDeviceGetPcieReplayCounter(device, &value)
	return value, ret
}

// ixml.DeviceGetEccErros()
func DeviceGetEccErros(device Device) (uint32, uint32, Return) {
	return device.GetEccErros()
}

func (device Device) GetEccErros() (uint32, uint32, Return) {
	var singleErr, doubleErr uint32
	ret := ixmlDeviceGetEccErros(device, &singleErr, &doubleErr)
	return singleErr, doubleErr, ret
}

// ixml.DeviceGetHealth()
func DeviceGetHealth(device Device) (uint64, Return) {
	return device.GetHealth()
}

func (device Device) GetHealth() (uint64, Return) {
	var health uint64
	ret := ixmlDeviceGetHealth(device, &health)
	return health, ret
}

// ixml.DeviceGetTopology()
func DeviceGetTopology(device1, device2 Device) (GpuTopologyLevel, Return) {
	return device1.GetTopology(device2)
}

func (device Device) GetTopology(device2 Device) (GpuTopologyLevel, Return) {
	var pathInfo GpuTopologyLevel
	ret := nvmlDeviceGetTopologyCommonAncestor(device, device2, &pathInfo)
	return pathInfo, ret
}

// ixml.DeviceGetPowerManagementLimit()
func DeviceGetPowerManagementLimit(device Device) (uint32, Return) {
	return device.GetPowerManagementLimit()
}

func (device Device) GetPowerManagementLimit() (uint32, Return) {
	var limit uint32
	ret := nvmlDeviceGetPowerManagementLimit(device, &limit)
	return limit, ret
}

// ixml.DeviceGetPowerManagementLimitConstraints()
func DeviceGetPowerManagementLimitConstraints(device Device) (uint32, uint32, Return) {
	return device.GetPowerManagementLimitConstraints()
}

func (device Device) GetPowerManagementLimitConstraints() (uint32, uint32, Return) {
	var minLimit, maxLimit uint32
	ret := nvmlDeviceGetPowerManagementLimitConstraints(device, &minLimit, &maxLimit)
	return minLimit, maxLimit, ret
}

// ixml.DeviceGetPowerManagementDefaultLimit()
func DeviceGetPowerManagementDefaultLimit(device Device) (uint32, Return) {
	return device.GetPowerManagementDefaultLimit()
}

func (device Device) GetPowerManagementDefaultLimit() (uint32, Return) {
	var defaultLimit uint32
	ret := nvmlDeviceGetPowerManagementDefaultLimit(device, &defaultLimit)
	return defaultLimit, ret
}

// ixml.DeviceGetTemperatureThreshold()
func DeviceGetTemperatureThreshold(device Device, thresholdType TemperatureThresholds) (uint32, Return) {
	return device.GetTemperatureThreshold(thresholdType)
}

func (device Device) GetTemperatureThreshold(thresholdType TemperatureThresholds) (uint32, Return) {
	var temp uint32
	ret := nvmlDeviceGetTemperatureThreshold(device, thresholdType, &temp)
	return temp, ret
}

// ixml.DeviceRegisterEvents()
func DeviceRegisterEvents(device Device, eventTypes uint64, set EventSet) Return {
	return device.RegisterEvents(eventTypes, set)
}

func (device Device) RegisterEvents(eventTypes uint64, set EventSet) Return {
	return nvmlDeviceRegisterEvents(device, eventTypes, set.(nvmlEventSet))
}

// ixml.DeviceGetSupportedEventTypes()
func DeviceGetSupportedEventTypes(device Device) (uint64, Return) {
	return device.GetSupportedEventTypes()
}

func (device Device) GetSupportedEventTypes() (uint64, Return) {
	var eventTypes uint64
	ret := nvmlDeviceGetSupportedEventTypes(device, &eventTypes)
	return eventTypes, ret
}

// ixml.DeviceGetBoardPartNumber()
func DeviceGetBoardPartNumber(device Device) (string, Return) {
	return device.GetBoardPartNumber()
}

func (device Device) GetBoardPartNumber() (string, Return) {
	partNumber := make([]byte, DEVICE_PART_NUMBER_BUFFER_SIZE)
	ret := nvmlDeviceGetBoardPartNumber(device, &partNumber[0], DEVICE_PART_NUMBER_BUFFER_SIZE)
	return string(partNumber[:clen(partNumber)]), ret
}
