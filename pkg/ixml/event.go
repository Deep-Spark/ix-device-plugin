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

// #cgo LDFLAGS: -ldl
// #include "ixml.h"
import "C"
import (
	"fmt"
)

type EventType uint64
type eventSetHandle C.nvmlEventSet_t

type eventSet struct {
	set eventSetHandle
}

const (
	eventTypeXidCriticalError EventType = C.ixmlEventTypeXidCriticalError
)

func (e *eventSet) RegisterEventsForDevice(uuid string, eventType EventType) error {
	// get device handle from cache list.
	dev, ok := cachedDevicesByUUID[uuid]
	if !ok {
		return fmt.Errorf("Cannot find the device by UUID: %v", uuid)
	}

	ret := C.ixmlDeviceRegisterEvents(dev, C.ulonglong(eventType), e.set)
	if ret != C.NVML_SUCCESS {
		return fmt.Errorf("Failed to register events for device: %v", uuid)
	}

	return nil
}

func (e *eventSet) WaitForEvent(timeout uint) (EventData, error) {
	var data C.ixmlEventData_t

	ret := C.ixmlEventSetWait(e.set, &data, C.uint(timeout))
	if ret != C.NVML_SUCCESS {
		return EventData{}, fmt.Errorf("Failed to wait for event")
	}

	return EventData{
		Type: EventType(data.eventType),
		Data: uint64(data.eventData),
	}, nil
}

func (e *eventSet) EventSetFree() error {
	ret := C.ixmlEventSetFree(e.set)
	if ret != C.NVML_SUCCESS {
		return fmt.Errorf("Failed to free event set.")
	}

	return nil
}

func newEventSet() (eventSetHandle, error) {
	var set C.nvmlEventSet_t

	ret := C.ixmlEventSetCreate(&set)
	if ret != C.NVML_SUCCESS {
		return nil, fmt.Errorf("Failed to create event set.")
	}

	return eventSetHandle(set), nil
}
