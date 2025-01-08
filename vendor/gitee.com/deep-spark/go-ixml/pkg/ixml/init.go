package ixml

import (
	"fmt"

	"gitee.com/deep-spark/go-ixml/pkg/dl"
)

import "C"

const (
	ixmlLibraryName      = "libixml.so"
	ixmlLibraryLoadFlags = dl.RTLD_LAZY | dl.RTLD_GLOBAL
)

var ixml *dl.DynamicLibrary

// ixml_init()
func Init() Return {
	lib := dl.New(ixmlLibraryName, ixmlLibraryLoadFlags)
	err := lib.Open()
	if err != nil {
		return ERROR_LIBRARY_NOT_FOUND
	}
	ixml = lib

	return nvmlInit()
}

func AbsInit(path string) Return {
	lib := dl.New(path, ixmlLibraryLoadFlags)
	err := lib.Open()
	if err != nil {
		return ERROR_LIBRARY_NOT_FOUND
	}
	ixml = lib

	return nvmlInit()
}

// ixml_shutdown()
func Shutdown() Return {
	ret := nvmlShutdown()
	if ret != SUCCESS {
		return ret
	}

	err := ixml.Close()
	if err != nil {
		panic(fmt.Sprintf("error closing %s: %v", ixmlLibraryName, err))
	}

	return ret
}
