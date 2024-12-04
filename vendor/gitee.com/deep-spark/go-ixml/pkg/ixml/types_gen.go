// Code generated by cmd/cgo -godefs; DO NOT EDIT.
// cgo -godefs types.go

package ixml

type Device struct {
	Handle *_Ctype_struct_nvmlDevice_st
}

type EventSet struct {
	Handle *_Ctype_struct_nvmlEventSet_st
}

type Memory struct {
	Total uint64
	Free  uint64
	Used  uint64
}

type Memory_v2 struct {
	Version  uint32
	Total    uint64
	Reserved uint64
	Free     uint64
	Used     uint64
}

type Utilization struct {
	Gpu    uint32
	Memory uint32
}

type ProcessInfo struct {
	Pid                      uint32
	UsedGpuMemory            uint64
	GpuInstanceId            uint32
	ComputeInstanceId        uint32
	UsedGpuCcProtectedMemory uint64
}

type ProcessInfo_v1 struct {
	Pid           uint32
	UsedGpuMemory uint64
}

type GpmSample struct {
	Handle *_Ctype_struct_nvmlGpmSample_st
}

type GpmMetric struct {
	MetricId   uint32
	NvmlReturn uint32
	Value      float64
	MetricInfo _Ctype_struct___3
}

type nvmlGpmMetricsGetType struct {
	Version    uint32
	NumMetrics uint32
	Sample1    GpmSample
	Sample2    GpmSample
	Metrics    [98]GpmMetric
}

type GpmSupport struct {
	Version           uint32
	IsSupportedDevice uint32
}

type PciInfo struct {
	BusIdLegacy    [16]int8
	Domain         uint32
	Bus            uint32
	Device         uint32
	PciDeviceId    uint32
	PciSubSystemId uint32
	BusId          [32]int8
}