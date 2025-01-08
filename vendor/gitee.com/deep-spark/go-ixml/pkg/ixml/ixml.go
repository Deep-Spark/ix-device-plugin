// WARNING: This file has automatically been generated on Thu, 28 Nov 2024 17:23:16 CST.
// Code generated by https://git.io/c-for-go. DO NOT EDIT.

package ixml

/*
#cgo LDFLAGS: -Wl,--export-dynamic -Wl,--unresolved-symbols=ignore-in-object-files
#cgo CFLAGS: -DNVML_NO_UNVERSIONED_FUNC_DEFS=1
#include "api.h"
#include <stdlib.h>
#include "cgo_helpers.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// nvmlInit function as declared in ixml/api.h:470
func nvmlInit() Return {
	__ret := C.nvmlInit_v2()
	__v := (Return)(__ret)
	return __v
}

// nvmlShutdown function as declared in ixml/api.h:487
func nvmlShutdown() Return {
	__ret := C.nvmlShutdown()
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetCount function as declared in ixml/api.h:509
func nvmlDeviceGetCount(DeviceCount *uint32) Return {
	cDeviceCount, cDeviceCountAllocMap := (*C.uint)(unsafe.Pointer(DeviceCount)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetCount_v2(cDeviceCount)
	runtime.KeepAlive(cDeviceCountAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetHandleByIndex function as declared in ixml/api.h:557
func nvmlDeviceGetHandleByIndex(Index uint32, Device *Device) Return {
	cIndex, cIndexAllocMap := (C.uint)(Index), cgoAllocsUnknown
	cDevice, cDeviceAllocMap := (*C.nvmlDevice_t)(unsafe.Pointer(Device)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetHandleByIndex_v2(cIndex, cDevice)
	runtime.KeepAlive(cDeviceAllocMap)
	runtime.KeepAlive(cIndexAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetHandleByUUID function as declared in ixml/api.h:582
func nvmlDeviceGetHandleByUUID(Uuid string, Device *Device) Return {
	cUuid, cUuidAllocMap := unpackPCharString(Uuid)
	cDevice, cDeviceAllocMap := (*C.nvmlDevice_t)(unsafe.Pointer(Device)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetHandleByUUID(cUuid, cDevice)
	runtime.KeepAlive(cDeviceAllocMap)
	runtime.KeepAlive(cUuidAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetMinorNumber function as declared in ixml/api.h:601
func nvmlDeviceGetMinorNumber(Device Device, MinorNumber *uint32) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cMinorNumber, cMinorNumberAllocMap := (*C.uint)(unsafe.Pointer(MinorNumber)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetMinorNumber(cDevice, cMinorNumber)
	runtime.KeepAlive(cMinorNumberAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetUUID function as declared in ixml/api.h:629
func nvmlDeviceGetUUID(Device Device, Uuid *byte, Length uint32) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cUuid, cUuidAllocMap := (*C.char)(unsafe.Pointer(Uuid)), cgoAllocsUnknown
	cLength, cLengthAllocMap := (C.uint)(Length), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetUUID(cDevice, cUuid, cLength)
	runtime.KeepAlive(cLengthAllocMap)
	runtime.KeepAlive(cUuidAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetName function as declared in ixml/api.h:655
func nvmlDeviceGetName(Device Device, Name *byte, Length uint32) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cName, cNameAllocMap := (*C.char)(unsafe.Pointer(Name)), cgoAllocsUnknown
	cLength, cLengthAllocMap := (C.uint)(Length), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetName(cDevice, cName, cLength)
	runtime.KeepAlive(cLengthAllocMap)
	runtime.KeepAlive(cNameAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlSystemGetDriverVersion function as declared in ixml/api.h:674
func nvmlSystemGetDriverVersion(Version *byte, Length uint32) Return {
	cVersion, cVersionAllocMap := (*C.char)(unsafe.Pointer(Version)), cgoAllocsUnknown
	cLength, cLengthAllocMap := (C.uint)(Length), cgoAllocsUnknown
	__ret := C.nvmlSystemGetDriverVersion(cVersion, cLength)
	runtime.KeepAlive(cLengthAllocMap)
	runtime.KeepAlive(cVersionAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlSystemGetCudaDriverVersion function as declared in ixml/api.h:690
func nvmlSystemGetCudaDriverVersion(CudaDriverVersion *int32) Return {
	cCudaDriverVersion, cCudaDriverVersionAllocMap := (*C.int)(unsafe.Pointer(CudaDriverVersion)), cgoAllocsUnknown
	__ret := C.nvmlSystemGetCudaDriverVersion(cCudaDriverVersion)
	runtime.KeepAlive(cCudaDriverVersionAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetTemperature function as declared in ixml/api.h:711
func nvmlDeviceGetTemperature(Device Device, SensorType TemperatureSensors, Temp *uint32) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cSensorType, cSensorTypeAllocMap := (C.nvmlTemperatureSensors_t)(SensorType), cgoAllocsUnknown
	cTemp, cTempAllocMap := (*C.uint)(unsafe.Pointer(Temp)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetTemperature(cDevice, cSensorType, cTemp)
	runtime.KeepAlive(cTempAllocMap)
	runtime.KeepAlive(cSensorTypeAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlSystemGetCudaDriverVersion_v2 function as declared in ixml/api.h:728
func nvmlSystemGetCudaDriverVersion_v2(CudaDriverVersion *int32) Return {
	cCudaDriverVersion, cCudaDriverVersionAllocMap := (*C.int)(unsafe.Pointer(CudaDriverVersion)), cgoAllocsUnknown
	__ret := C.nvmlSystemGetCudaDriverVersion_v2(cCudaDriverVersion)
	runtime.KeepAlive(cCudaDriverVersionAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetFanSpeed function as declared in ixml/api.h:752
func nvmlDeviceGetFanSpeed(Device Device, Speed *uint32) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cSpeed, cSpeedAllocMap := (*C.uint)(unsafe.Pointer(Speed)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetFanSpeed(cDevice, cSpeed)
	runtime.KeepAlive(cSpeedAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetClockInfo function as declared in ixml/api.h:773
func nvmlDeviceGetClockInfo(Device Device, _type ClockType, Clock *uint32) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	c_type, c_typeAllocMap := (C.nvmlClockType_t)(_type), cgoAllocsUnknown
	cClock, cClockAllocMap := (*C.uint)(unsafe.Pointer(Clock)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetClockInfo(cDevice, c_type, cClock)
	runtime.KeepAlive(cClockAllocMap)
	runtime.KeepAlive(c_typeAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetMemoryInfo function as declared in ixml/api.h:806
func nvmlDeviceGetMemoryInfo(Device Device, Memory *Memory) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cMemory, cMemoryAllocMap := (*C.nvmlMemory_t)(unsafe.Pointer(Memory)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetMemoryInfo(cDevice, cMemory)
	runtime.KeepAlive(cMemoryAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetMemoryInfo_v2 function as declared in ixml/api.h:807
func nvmlDeviceGetMemoryInfo_v2(Device Device, Memory *Memory_v2) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cMemory, cMemoryAllocMap := (*C.nvmlMemory_v2_t)(unsafe.Pointer(Memory)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetMemoryInfo_v2(cDevice, cMemory)
	runtime.KeepAlive(cMemoryAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetFanSpeed_v2 function as declared in ixml/api.h:832
func nvmlDeviceGetFanSpeed_v2(Device Device, Fan uint32, Speed *uint32) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cFan, cFanAllocMap := (C.uint)(Fan), cgoAllocsUnknown
	cSpeed, cSpeedAllocMap := (*C.uint)(unsafe.Pointer(Speed)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetFanSpeed_v2(cDevice, cFan, cSpeed)
	runtime.KeepAlive(cSpeedAllocMap)
	runtime.KeepAlive(cFanAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetUtilizationRates function as declared in ixml/api.h:857
func nvmlDeviceGetUtilizationRates(Device Device, Utilization *Utilization) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cUtilization, cUtilizationAllocMap := (*C.nvmlUtilization_t)(unsafe.Pointer(Utilization)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetUtilizationRates(cDevice, cUtilization)
	runtime.KeepAlive(cUtilizationAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetPciInfo function as declared in ixml/api.h:876
func nvmlDeviceGetPciInfo(Device Device, Pci *PciInfo) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cPci, cPciAllocMap := (*C.nvmlPciInfo_t)(unsafe.Pointer(Pci)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetPciInfo_v3(cDevice, cPci)
	runtime.KeepAlive(cPciAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetIndex function as declared in ixml/api.h:910
func nvmlDeviceGetIndex(Device Device, Index *uint32) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cIndex, cIndexAllocMap := (*C.uint)(unsafe.Pointer(Index)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetIndex(cDevice, cIndex)
	runtime.KeepAlive(cIndexAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetPowerUsage function as declared in ixml/api.h:932
func nvmlDeviceGetPowerUsage(Device Device, Power *uint32) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cPower, cPowerAllocMap := (*C.uint)(unsafe.Pointer(Power)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetPowerUsage(cDevice, cPower)
	runtime.KeepAlive(cPowerAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceOnSameBoard function as declared in ixml/api.h:952
func nvmlDeviceOnSameBoard(Device1 Device, Device2 Device, OnSameBoard *int32) Return {
	cDevice1, cDevice1AllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device1)), cgoAllocsUnknown
	cDevice2, cDevice2AllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device2)), cgoAllocsUnknown
	cOnSameBoard, cOnSameBoardAllocMap := (*C.int)(unsafe.Pointer(OnSameBoard)), cgoAllocsUnknown
	__ret := C.nvmlDeviceOnSameBoard(cDevice1, cDevice2, cOnSameBoard)
	runtime.KeepAlive(cOnSameBoardAllocMap)
	runtime.KeepAlive(cDevice2AllocMap)
	runtime.KeepAlive(cDevice1AllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetComputeRunningProcesses function as declared in ixml/api.h:995
func nvmlDeviceGetComputeRunningProcesses(Device Device, InfoCount *uint32, Infos *ProcessInfo_v1) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cInfoCount, cInfoCountAllocMap := (*C.uint)(unsafe.Pointer(InfoCount)), cgoAllocsUnknown
	cInfos, cInfosAllocMap := (*C.nvmlProcessInfo_v1_t)(unsafe.Pointer(Infos)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetComputeRunningProcesses(cDevice, cInfoCount, cInfos)
	runtime.KeepAlive(cInfosAllocMap)
	runtime.KeepAlive(cInfoCountAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetPcieReplayCounter function as declared in ixml/api.h:1017
func nvmlDeviceGetPcieReplayCounter(Device Device, Value *uint32) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cValue, cValueAllocMap := (*C.uint)(unsafe.Pointer(Value)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetPcieReplayCounter(cDevice, cValue)
	runtime.KeepAlive(cValueAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlGpmMetricsGet function as declared in ixml/api.h:1038
func nvmlGpmMetricsGet(MetricsGet *nvmlGpmMetricsGetType) Return {
	cMetricsGet, cMetricsGetAllocMap := (*C.nvmlGpmMetricsGet_t)(unsafe.Pointer(MetricsGet)), cgoAllocsUnknown
	__ret := C.nvmlGpmMetricsGet(cMetricsGet)
	runtime.KeepAlive(cMetricsGetAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlGpmQueryDeviceSupport function as declared in ixml/api.h:1052
func nvmlGpmQueryDeviceSupport(Device Device, GpmSupport *GpmSupport) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cGpmSupport, cGpmSupportAllocMap := (*C.nvmlGpmSupport_t)(unsafe.Pointer(GpmSupport)), cgoAllocsUnknown
	__ret := C.nvmlGpmQueryDeviceSupport(cDevice, cGpmSupport)
	runtime.KeepAlive(cGpmSupportAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlGpmSampleFree function as declared in ixml/api.h:1065
func nvmlGpmSampleFree(GpmSample GpmSample) Return {
	cGpmSample, cGpmSampleAllocMap := *(*C.nvmlGpmSample_t)(unsafe.Pointer(&GpmSample)), cgoAllocsUnknown
	__ret := C.nvmlGpmSampleFree(cGpmSample)
	runtime.KeepAlive(cGpmSampleAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlGpmSampleAlloc function as declared in ixml/api.h:1080
func nvmlGpmSampleAlloc(GpmSample *GpmSample) Return {
	cGpmSample, cGpmSampleAllocMap := (*C.nvmlGpmSample_t)(unsafe.Pointer(GpmSample)), cgoAllocsUnknown
	__ret := C.nvmlGpmSampleAlloc(cGpmSample)
	runtime.KeepAlive(cGpmSampleAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlGpmSampleGet function as declared in ixml/api.h:1096
func nvmlGpmSampleGet(Device Device, GpmSample GpmSample) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cGpmSample, cGpmSampleAllocMap := *(*C.nvmlGpmSample_t)(unsafe.Pointer(&GpmSample)), cgoAllocsUnknown
	__ret := C.nvmlGpmSampleGet(cDevice, cGpmSample)
	runtime.KeepAlive(cGpmSampleAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetPowerManagementLimitConstraints function as declared in ixml/api.h:1117
func nvmlDeviceGetPowerManagementLimitConstraints(Device Device, MinLimit *uint32, MaxLimit *uint32) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cMinLimit, cMinLimitAllocMap := (*C.uint)(unsafe.Pointer(MinLimit)), cgoAllocsUnknown
	cMaxLimit, cMaxLimitAllocMap := (*C.uint)(unsafe.Pointer(MaxLimit)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetPowerManagementLimitConstraints(cDevice, cMinLimit, cMaxLimit)
	runtime.KeepAlive(cMaxLimitAllocMap)
	runtime.KeepAlive(cMinLimitAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetCurrentClocksThrottleReasons function as declared in ixml/api.h:1119
func nvmlDeviceGetCurrentClocksThrottleReasons(Device Device, ClocksThrottleReasons *uint64) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cClocksThrottleReasons, cClocksThrottleReasonsAllocMap := (*C.ulonglong)(unsafe.Pointer(ClocksThrottleReasons)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetCurrentClocksThrottleReasons(cDevice, cClocksThrottleReasons)
	runtime.KeepAlive(cClocksThrottleReasonsAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// nvmlDeviceGetTopologyCommonAncestor function as declared in ixml/api.h:1138
func nvmlDeviceGetTopologyCommonAncestor(Device1 Device, Device2 Device, PathInfo *GpuTopologyLevel) Return {
	cDevice1, cDevice1AllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device1)), cgoAllocsUnknown
	cDevice2, cDevice2AllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device2)), cgoAllocsUnknown
	cPathInfo, cPathInfoAllocMap := (*C.nvmlGpuTopologyLevel_t)(unsafe.Pointer(PathInfo)), cgoAllocsUnknown
	__ret := C.nvmlDeviceGetTopologyCommonAncestor(cDevice1, cDevice2, cPathInfo)
	runtime.KeepAlive(cPathInfoAllocMap)
	runtime.KeepAlive(cDevice2AllocMap)
	runtime.KeepAlive(cDevice1AllocMap)
	__v := (Return)(__ret)
	return __v
}

// ixmlDeviceGetBoardPosition function as declared in ixml/api.h:1140
func ixmlDeviceGetBoardPosition(Device Device, Position *uint32) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cPosition, cPositionAllocMap := (*C.uint)(unsafe.Pointer(Position)), cgoAllocsUnknown
	__ret := C.ixmlDeviceGetBoardPosition(cDevice, cPosition)
	runtime.KeepAlive(cPositionAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// ixmlDeviceGetGPUVoltage function as declared in ixml/api.h:1142
func ixmlDeviceGetGPUVoltage(Device Device, Integer *uint32, Decimal *uint32) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cInteger, cIntegerAllocMap := (*C.uint)(unsafe.Pointer(Integer)), cgoAllocsUnknown
	cDecimal, cDecimalAllocMap := (*C.uint)(unsafe.Pointer(Decimal)), cgoAllocsUnknown
	__ret := C.ixmlDeviceGetGPUVoltage(cDevice, cInteger, cDecimal)
	runtime.KeepAlive(cDecimalAllocMap)
	runtime.KeepAlive(cIntegerAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// ixmlDeviceGetEccErros function as declared in ixml/api.h:1144
func ixmlDeviceGetEccErros(Device Device, Single_error *uint32, Double_error *uint32) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cSingle_error, cSingle_errorAllocMap := (*C.uint)(unsafe.Pointer(Single_error)), cgoAllocsUnknown
	cDouble_error, cDouble_errorAllocMap := (*C.uint)(unsafe.Pointer(Double_error)), cgoAllocsUnknown
	__ret := C.ixmlDeviceGetEccErros(cDevice, cSingle_error, cDouble_error)
	runtime.KeepAlive(cDouble_errorAllocMap)
	runtime.KeepAlive(cSingle_errorAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}

// ixmlDeviceGetHealth function as declared in ixml/api.h:1146
func ixmlDeviceGetHealth(Device Device, Health *uint64) Return {
	cDevice, cDeviceAllocMap := *(*C.nvmlDevice_t)(unsafe.Pointer(&Device)), cgoAllocsUnknown
	cHealth, cHealthAllocMap := (*C.ulonglong)(unsafe.Pointer(Health)), cgoAllocsUnknown
	__ret := C.ixmlDeviceGetHealth(cDevice, cHealth)
	runtime.KeepAlive(cHealthAllocMap)
	runtime.KeepAlive(cDeviceAllocMap)
	__v := (Return)(__ret)
	return __v
}
