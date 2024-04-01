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

#include <dlfcn.h>
#include <stdlib.h>
#include <stdio.h>
#include <stdarg.h>

#include "ixml.h"

void *ixmlHandle;

nvmlReturn_t (*ixmlInitFunc)();
nvmlReturn_t ixmlInit() {
    if (ixmlInitFunc == NULL) {
        return 1;
    }

    return ixmlInitFunc();
}

nvmlReturn_t (*ixmlShutdownFunc)();
nvmlReturn_t ixmlShutdown() {
    if (ixmlShutdownFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }
    return ixmlShutdownFunc();
}

nvmlReturn_t (*ixmlDeviceGetCountFunc)(unsigned int* deviceCount);
nvmlReturn_t ixmlDeviceGetCount(unsigned int* deviceCount) {
    if (ixmlDeviceGetCountFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }
    return ixmlDeviceGetCountFunc(deviceCount);
}

nvmlReturn_t (*ixmlSystemGetDriverVersionFunc)(char *version, unsigned int length);
nvmlReturn_t ixmlSystemGetDriverVersion(char *version, unsigned int length) {
    if (ixmlSystemGetDriverVersionFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }
    return ixmlSystemGetDriverVersionFunc(version, length);
}

nvmlReturn_t (*ixmlSystemGetCudaDriverVersionFunc)(int *version);
nvmlReturn_t ixmlSystemGetCudaDriverVersion(int *version) {
    if (ixmlSystemGetCudaDriverVersionFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    return ixmlSystemGetCudaDriverVersionFunc(version);
}

nvmlReturn_t (*ixmlDeviceGetHandleByIndexFunc)(unsigned int index, nvmlDevice_t* device);
nvmlReturn_t ixmlDeviceGetHandleByIndex(unsigned int index, nvmlDevice_t* device) {
    if (ixmlDeviceGetHandleByIndexFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    ixmlDeviceGetHandleByIndexFunc(index, device);
}

nvmlReturn_t (*ixmlDeviceGetHandleByUUIDFunc)(const char *uuid, nvmlDevice_t* device);
nvmlReturn_t ixmlDeviceGetHandleByUUID(const char *uuid, nvmlDevice_t* device) {
    if (ixmlDeviceGetHandleByUUIDFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    ixmlDeviceGetHandleByUUIDFunc(uuid, device);
}

nvmlReturn_t (*ixmlDeviceGetNameFunc)(nvmlDevice_t device, char* name, unsigned int length);
nvmlReturn_t ixmlDeviceGetName(nvmlDevice_t device, char* name, unsigned int length) {
    if (ixmlDeviceGetNameFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    return ixmlDeviceGetNameFunc(device, name, length);
}

nvmlReturn_t (*ixmlDeviceGetUUIDFunc)(nvmlDevice_t device, char* uuid, unsigned int length);
nvmlReturn_t ixmlDeviceGetUUID(nvmlDevice_t device, char* uuid, unsigned int length) {
    if (ixmlDeviceGetUUIDFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    return ixmlDeviceGetUUIDFunc(device, uuid, length);
}

nvmlReturn_t (*ixmlDeviceGetIndexFunc)(nvmlDevice_t device, unsigned int *index);
nvmlReturn_t ixmlDeviceGetIndex(nvmlDevice_t device, unsigned int *index) {
    if (ixmlDeviceGetIndexFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    return ixmlDeviceGetIndexFunc(device, index);
}

nvmlReturn_t (*ixmlDeviceGetFanSpeedFunc)(nvmlDevice_t device, unsigned int* speed);
nvmlReturn_t ixmlDeviceGetFanSpeed(nvmlDevice_t device, unsigned int* speed) {
    if (ixmlDeviceGetFanSpeedFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    ixmlDeviceGetFanSpeedFunc(device, speed);
}

nvmlReturn_t (*ixmlDeviceGetMemoryInfoFunc)(nvmlDevice_t device, nvmlMemory_t* memory);
nvmlReturn_t ixmlDeviceGetMemoryInfo(nvmlDevice_t device, nvmlMemory_t* memory) {
    if (ixmlDeviceGetMemoryInfoFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    return ixmlDeviceGetMemoryInfoFunc(device, memory);
}

nvmlReturn_t (*ixmlDeviceGetTemperatureFunc)(nvmlDevice_t device, nvmlTemperatureSensors_t sensorType, unsigned int* temp);
nvmlReturn_t ixmlDeviceGetTemperature(nvmlDevice_t device, nvmlTemperatureSensors_t sensorType, unsigned int* temp) {
    if (ixmlDeviceGetTemperatureFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    return ixmlDeviceGetTemperatureFunc(device, sensorType, temp);
}

nvmlReturn_t (*ixmlDeviceGetPciInfoFunc)(nvmlDevice_t device, nvmlPciInfo_t* pci);
nvmlReturn_t ixmlDeviceGetPciInfo(nvmlDevice_t device, nvmlPciInfo_t* pci) {
    if (ixmlDeviceGetPciInfoFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    return ixmlDeviceGetPciInfoFunc(device, pci);
}

nvmlReturn_t (*ixmlDeviceGetMinorNumberFunc)(nvmlDevice_t device, unsigned int* minorNumber);
nvmlReturn_t ixmlDeviceGetMinorNumber(nvmlDevice_t device, unsigned int* minorNumber) {
    if (ixmlDeviceGetMinorNumberFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    return ixmlDeviceGetMinorNumberFunc(device, minorNumber);
}

nvmlReturn_t (*ixmlDeviceGetPowerUsageFunc)(nvmlDevice_t device, unsigned int* power);
nvmlReturn_t ixmlDeviceGetPowerUsage(nvmlDevice_t device, unsigned int* power) {
    if (ixmlDeviceGetPowerUsageFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    return ixmlDeviceGetPowerUsageFunc(device, power);
}

nvmlReturn_t (*ixmlDeviceGetPowerManagementLimitConstraintsFunc)(nvmlDevice_t device, unsigned int* minLimit, unsigned int* maxLimit);
nvmlReturn_t ixmlDeviceGetPowerManagementLimitConstraints(nvmlDevice_t device, unsigned int* minLimit, unsigned int* maxLimit) {
    if (ixmlDeviceGetPowerManagementLimitConstraintsFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    return ixmlDeviceGetPowerManagementLimitConstraintsFunc(device, minLimit, maxLimit);
}

nvmlReturn_t (*ixmlDeviceGetClockInfoFunc)(nvmlDevice_t device, nvmlClockType_t type, unsigned int* clock);
nvmlReturn_t ixmlDeviceGetClockInfo(nvmlDevice_t device, nvmlClockType_t type, unsigned int* clock) {
    if (ixmlDeviceGetClockInfoFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    return ixmlDeviceGetClockInfoFunc(device, type, clock);
}

nvmlReturn_t (*ixmlDeviceGetUtilizationRatesFunc)(nvmlDevice_t device, nvmlUtilization_t* utilization);
nvmlReturn_t ixmlDeviceGetUtilizationRates(nvmlDevice_t device, nvmlUtilization_t* utilization) {
    if (ixmlDeviceGetUtilizationRatesFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    return ixmlDeviceGetUtilizationRatesFunc(device, utilization);
}

nvmlReturn_t (*ixmlEventSetCreateFunc)(nvmlEventSet_t *set);
nvmlReturn_t ixmlEventSetCreate(nvmlEventSet_t *set) {
    if (ixmlEventSetCreateFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    return ixmlEventSetCreateFunc(set);
}

nvmlReturn_t (*ixmlEventSetFreeFunc)(nvmlEventSet_t set);
nvmlReturn_t ixmlEventSetFree(nvmlEventSet_t set) {
    if (ixmlEventSetFreeFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    return ixmlEventSetFreeFunc(set);
}

nvmlReturn_t (*ixmlDeviceRegisterEventsFunc)(nvmlDevice_t device, unsigned long long eventTypes, nvmlEventSet_t set);
nvmlReturn_t ixmlDeviceRegisterEvents(nvmlDevice_t device, unsigned long long eventTypes, nvmlEventSet_t set) {
    if (ixmlDeviceRegisterEventsFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    return ixmlDeviceRegisterEventsFunc(device, eventTypes, set);
}

nvmlReturn_t (*ixmlEventSetWaitFunc)(nvmlEventSet_t set, nvmlEventData_t * data, unsigned int timeoutms);
nvmlReturn_t ixmlEventSetWait(nvmlEventSet_t set, nvmlEventData_t * data, unsigned int timeoutms) {
    if (ixmlEventSetWaitFunc == NULL) {
        return NVML_ERROR_UNKNOWN;
    }

    return ixmlEventSetWaitFunc(set, data, timeoutms);
}

nvmlReturn_t (*ixmlDeviceOnSameBoardFunc)(nvmlDevice_t device1, nvmlDevice_t device2, int* onSameBoard);
nvmlReturn_t ixmlDeviceOnSameBoard(nvmlDevice_t device1, nvmlDevice_t device2, int* onSameBoard) {
	if (ixmlDeviceOnSameBoardFunc == NULL){
	    return NVML_ERROR_UNKNOWN;
	}

	return ixmlDeviceOnSameBoardFunc(device1, device2, onSameBoard);
}

nvmlReturn_t dl_init() {
	ixmlHandle = dlopen(IXML_LIBRARY, RTLD_LAZY|RTLD_GLOBAL);
    if (ixmlHandle == NULL) {
        return NVML_ERROR_LIBRARY_NOT_FOUND;
    }

    ixmlInitFunc = dlsym(ixmlHandle, IXML_INIT);
    if (ixmlInitFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlShutdownFunc = dlsym(ixmlHandle, IXML_SHUTDOWN);
    if (ixmlShutdownFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlDeviceGetCountFunc = dlsym(ixmlHandle, IXML_DEVICE_GET_COUNT);
    if (ixmlDeviceGetCountFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlSystemGetDriverVersionFunc  = dlsym(ixmlHandle, IXML_GET_DRIVER_VERSION);
    if (ixmlSystemGetDriverVersionFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlSystemGetCudaDriverVersionFunc = dlsym(ixmlHandle, IXML_GET_CUDA_DRIVER_VERSION);
    if (ixmlSystemGetCudaDriverVersionFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlDeviceGetHandleByIndexFunc = dlsym(ixmlHandle, IXML_DEVICE_GET_HANDLE_BY_INDEX);
    if (ixmlDeviceGetHandleByIndexFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlDeviceGetHandleByUUIDFunc = dlsym(ixmlHandle, IXML_DEVICE_GET_HANDLE_BY_UUID);
    if (ixmlDeviceGetHandleByUUIDFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlDeviceGetNameFunc = dlsym(ixmlHandle, IXML_DEVICE_GET_NAME);
    if (ixmlDeviceGetNameFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlDeviceGetUUIDFunc = dlsym(ixmlHandle, IXML_DEVICE_GET_UUID);
    if (ixmlDeviceGetUUIDFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlDeviceGetIndexFunc = dlsym(ixmlHandle, IXML_DEVICE_GET_INDEX);
    if (ixmlDeviceGetIndexFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlDeviceGetFanSpeedFunc = dlsym(ixmlHandle, IXML_DEVICE_GET_FAN_SPEED);
    if (ixmlDeviceGetFanSpeedFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlDeviceGetMemoryInfoFunc = dlsym(ixmlHandle, IXML_DEVICE_GET_MEMORY_INFO);
    if (ixmlDeviceGetMemoryInfoFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlDeviceGetTemperatureFunc = dlsym(ixmlHandle, IXML_DEVICE_GET_TEMPERATURE);
    if (ixmlDeviceGetTemperatureFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlDeviceGetPciInfoFunc = dlsym(ixmlHandle, IXML_DEVICE_GET_PCI_INFO);
    if (ixmlDeviceGetPciInfoFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlDeviceGetMinorNumberFunc = dlsym(ixmlHandle, IXML_DEVICE_GET_MINOR_NUMBER);
    if (ixmlDeviceGetMinorNumberFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlDeviceGetPowerUsageFunc = dlsym(ixmlHandle, IXML_DEVICE_GET_POWER_USAGE);
    if (ixmlDeviceGetPowerUsageFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlDeviceGetPowerManagementLimitConstraintsFunc = dlsym(ixmlHandle, IXML_DEVICE_GET_POWER_CONSTRAINT);
    if (ixmlDeviceGetPowerManagementLimitConstraintsFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlDeviceGetClockInfoFunc = dlsym(ixmlHandle, IXML_DEVICE_GET_CLOCK_INFO);
    if (ixmlDeviceGetClockInfoFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlDeviceGetUtilizationRatesFunc = dlsym(ixmlHandle, IXML_DEVICE_GET_UTILIZATION_RATES);
    if (ixmlDeviceGetUtilizationRatesFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlEventSetCreateFunc = dlsym(ixmlHandle, IXML_EVENT_SET_CREATE);
    if (ixmlEventSetCreateFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlEventSetFreeFunc = dlsym(ixmlHandle, IXML_EVENT_SET_FREE);
    if (ixmlEventSetFreeFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlDeviceRegisterEventsFunc = dlsym(ixmlHandle, IXML_DEVICE_REGISTER_EVENTS);
    if (ixmlDeviceRegisterEventsFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlEventSetWaitFunc = dlsym(ixmlHandle, IXML_EVENT_SET_WAIT);
    if (ixmlEventSetWaitFunc == NULL) {
        return NVML_ERROR_FUNCTION_NOT_FOUND;
    }

    ixmlDeviceOnSameBoardFunc = dlsym(ixmlHandle, IXML_DEVICE_ON_SAME_BOARD);
    if (ixmlDeviceOnSameBoardFunc == NULL) {
    	return NVML_ERROR_FUNCTION_NOT_FOUND;
    }
}

nvmlReturn_t dl_close() {
    int ret;

    ret = dlclose(ixmlHandle);
    if (ret != 0) {
        return NVML_ERROR_UNKNOWN;
    }
}
