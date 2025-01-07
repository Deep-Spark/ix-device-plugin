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

func DeviceGetCount() (uint, Return) {
	var DeviceCount uint32
	ret := nvmlDeviceGetCount(&DeviceCount)
	return uint(DeviceCount), ret
}

func DeviceGetHandleByIndex(Index uint, device *Device) Return {
	ret := nvmlDeviceGetHandleByIndex(uint32(Index), device)
	return ret
}

func GetHandleByUUID(Uuid string) (Device, Return) {
	var Device Device
	ret := nvmlDeviceGetHandleByUUID(Uuid, &Device)
	return Device, ret
}

func (device Device) GetUUID() (string, Return) {
	return deviceGetUUID(device)
}

func deviceGetUUID(Device Device) (string, Return) {
	Uuid := make([]byte, DEVICE_UUID_BUFFER_SIZE)
	ret := nvmlDeviceGetUUID(Device, &Uuid[0], DEVICE_UUID_BUFFER_SIZE)
	return removeBytesSpaces(Uuid), ret
}

func (device Device) GetMinorNumber() (int, Return) {
	return deviceGetMinorNumber(device)
}

func deviceGetMinorNumber(Device Device) (int, Return) {
	var minorNumber uint32
	ret := nvmlDeviceGetMinorNumber(Device, &minorNumber)
	return int(minorNumber), ret
}

func (device Device) GetName() (string, Return) {
	return deviceGetName(device)
}

func deviceGetName(device Device) (string, Return) {
	Name := make([]byte, DEVICE_NAME_BUFFER_SIZE)
	ret := nvmlDeviceGetName(device, &Name[0], DEVICE_NAME_BUFFER_SIZE)
	removeBytesSpaces(Name)
	return removeBytesSpaces(Name), ret
}

func (device Device) GetTemperature() (uint32, Return) {
	var SensorType TemperatureSensors
	return deviceGetTemperature(device, SensorType)
}

func deviceGetTemperature(Device Device, SensorType TemperatureSensors) (uint32, Return) {
	var Temperature uint32
	ret := nvmlDeviceGetTemperature(Device, SensorType, &Temperature)
	return Temperature, ret
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
	ret := nvmlDeviceGetClockInfo(Device, CLOCK_SM, &sm)
	ret = nvmlDeviceGetClockInfo(Device, CLOCK_MEM, &mem)
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

func (device Device) GetPciInfo() (PciInfo, Return) {
	return deviceGetPciInfo(device)
}

func deviceGetPciInfo(Device Device) (PciInfo, Return) {
	var PciInfo PciInfo
	ret := nvmlDeviceGetPciInfo(Device, &PciInfo)
	return PciInfo, ret
}

func (device Device) GetIndex() (int, Return) {
	return deviceGetIndex(device)
}

func deviceGetIndex(device Device) (int, Return) {
	var Index uint32
	ret := nvmlDeviceGetIndex(device, &Index)
	return int(Index), ret
}

func (device Device) GetPowerUsage() (uint32, Return) {
	return deviceGetPowerUsage(device)
}

func deviceGetPowerUsage(Device Device) (uint32, Return) {
	var Power uint32
	ret := nvmlDeviceGetPowerUsage(Device, &Power)
	return Power, ret
}

func GetOnSameBoard(device1, device2 Device) (int, Return) {
	var OnSameBoard int32
	ret := nvmlDeviceOnSameBoard(device1, device2, &OnSameBoard)
	return int(OnSameBoard), ret
}

func (device Device) GetBoardPosition() (uint32, Return) {
	return deviceGetBoardPosition(device)
}

func deviceGetBoardPosition(device Device) (uint32, Return) {
	var pos uint32
	ret := ixmlDeviceGetBoardPosition(device, &pos)
	return pos, ret
}

func (device Device) GetGPUVoltage() (uint32, uint32, Return) {
	return deviceGetGPUVoltage(device)
}

func deviceGetGPUVoltage(device Device) (uint32, uint32, Return) {
	var Integer, Decimal uint32
	ret := ixmlDeviceGetGPUVoltage(device, &Integer, &Decimal)
	return Integer, Decimal, ret
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

func (device Device) GetCurrentClocksThrottleReasons() (uint64, Return) {
	return deviceGetCurrentClocksThrottleReasons(device)
}

func deviceGetCurrentClocksThrottleReasons(device Device) (uint64, Return) {
	var clocksThrottleReasons uint64
	ret := nvmlDeviceGetCurrentClocksThrottleReasons(device, &clocksThrottleReasons)
	return clocksThrottleReasons, ret
}

func (device Device) GetPcieReplayCounter() (uint32, Return) {
	return deviceGetPcieReplayCounter(device)
}

func deviceGetPcieReplayCounter(device Device) (uint32, Return) {
	var value uint32
	ret := nvmlDeviceGetPcieReplayCounter(device, &value)
	return value, ret
}

func (device Device) GetEccErros() (uint32, uint32, Return) {
	return deviceGetEccErros(device)
}

func deviceGetEccErros(device Device) (uint32, uint32, Return) {
	var singleErr, doubleErr uint32
	ret := ixmlDeviceGetEccErros(device, &singleErr, &doubleErr)
	return singleErr, doubleErr, ret
}

func (device Device) GetHealth() (uint64, Return) {
	return deviceGetHealth(device)
}

func deviceGetHealth(device Device) (uint64, Return) {
	var health uint64
	ret := ixmlDeviceGetHealth(device, &health)
	return health, ret
}

func (device Device) GetTopology(device2 Device) (GpuTopologyLevel, Return) {
	return deviceGetTopology(device, device2)
}

func deviceGetTopology(device, device2 Device) (GpuTopologyLevel, Return) {
	var pathInfo GpuTopologyLevel
	ret := nvmlDeviceGetTopologyCommonAncestor(device, device2, &pathInfo)
	return pathInfo, ret
}

func (device Device) GetPowerManagementLimitConstraints() (uint32, uint32, Return) {
	return deviceGetPowerManagementLimitConstraints(device)
}

func deviceGetPowerManagementLimitConstraints(device Device) (uint32, uint32, Return) {
	var minLimit, maxLimit uint32
	ret := nvmlDeviceGetPowerManagementLimitConstraints(device, &minLimit, &maxLimit)
	return minLimit, maxLimit, ret
}
