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
