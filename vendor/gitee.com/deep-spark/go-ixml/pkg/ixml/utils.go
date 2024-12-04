package ixml

import (
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
	return strings.TrimSuffix(string(data), "\x00")
}
