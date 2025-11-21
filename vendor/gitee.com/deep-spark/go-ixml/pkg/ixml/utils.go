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
	"bytes"
	"fmt"
	"os"
	"strings"
)

func removeBytesSpaces(originalBytes []byte) string {
	lastNonZeroIndex := len(originalBytes) - 1
	for ; lastNonZeroIndex >= 0; lastNonZeroIndex-- {
		if originalBytes[lastNonZeroIndex] != 0 {
			break
		}
	}
	cleanedBytes := originalBytes[:lastNonZeroIndex+1]

	return string(cleanedBytes)
}

func getPidName(pid uint32) string {
	cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", pid)
	data, err := os.ReadFile(cmdlinePath)
	if err != nil {
		return ""
	}
	data = bytes.ReplaceAll(data, []byte{0}, []byte{' '})
	return strings.TrimSuffix(string(data), "\x00")
}
