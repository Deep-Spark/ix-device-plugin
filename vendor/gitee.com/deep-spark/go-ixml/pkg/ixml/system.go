package ixml

import (
	"fmt"
)

func SystemGetDriverVersion() (string, Return) {
	version := make([]byte, SYSTEM_DRIVER_VERSION_BUFFER_SIZE)
	ret := nvmlSystemGetDriverVersion(&version[0], SYSTEM_DRIVER_VERSION_BUFFER_SIZE)
	return removeBytesSpaces(version), ret
}

func SystemGetCudaDriverVersion() (string, Return) {
	var CudaDriverVersion int32
	ret := nvmlSystemGetCudaDriverVersion(&CudaDriverVersion)
	return fmt.Sprintf("%d", CudaDriverVersion), ret
}

func SystemGetCudaDriverVersion_v2() (string, Return) {
	var CudaDriverVersion int32
	ret := nvmlSystemGetCudaDriverVersion_v2(&CudaDriverVersion)
	return fmt.Sprintf("%d", CudaDriverVersion), ret
}
