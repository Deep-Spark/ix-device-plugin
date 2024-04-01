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

#include <nvml.h>

#define IXML_LIBRARY "libixml.so"

#define IXML_INIT                         "nvmlInit"
#define IXML_SHUTDOWN                     "nvmlShutdown"
#define IXML_DEVICE_GET_COUNT             "nvmlDeviceGetCount"
#define IXML_GET_DRIVER_VERSION           "nvmlSystemGetDriverVersion"
#define IXML_GET_CUDA_DRIVER_VERSION      "nvmlSystemGetCudaDriverVersion"
#define IXML_DEVICE_GET_HANDLE_BY_INDEX   "nvmlDeviceGetHandleByIndex"
#define IXML_DEVICE_GET_HANDLE_BY_UUID    "nvmlDeviceGetHandleByUUID"
#define IXML_DEVICE_GET_NAME              "nvmlDeviceGetName"
#define IXML_DEVICE_GET_UUID              "nvmlDeviceGetUUID"
#define IXML_DEVICE_GET_INDEX             "nvmlDeviceGetIndex"
#define IXML_DEVICE_GET_FAN_SPEED         "nvmlDeviceGetFanSpeed"
#define IXML_DEVICE_GET_MEMORY_INFO       "nvmlDeviceGetMemoryInfo"
#define IXML_DEVICE_GET_TEMPERATURE       "nvmlDeviceGetTemperature"
#define IXML_DEVICE_GET_PCI_INFO          "nvmlDeviceGetPciInfo"
#define IXML_DEVICE_GET_MINOR_NUMBER      "nvmlDeviceGetMinorNumber"
#define IXML_DEVICE_GET_POWER_USAGE       "nvmlDeviceGetPowerUsage"
#define IXML_DEVICE_GET_POWER_CONSTRAINT  "nvmlDeviceGetPowerManagementLimitConstraints"
#define IXML_DEVICE_GET_CLOCK_INFO        "nvmlDeviceGetClockInfo"
#define IXML_DEVICE_GET_UTILIZATION_RATES "nvmlDeviceGetUtilizationRates"
#define IXML_EVENT_SET_CREATE             "nvmlEventSetCreate"
#define IXML_EVENT_SET_FREE               "nvmlEventSetFree"
#define IXML_DEVICE_REGISTER_EVENTS       "nvmlDeviceRegisterEvents"
#define IXML_EVENT_SET_WAIT               "nvmlEventSetWait"
#define IXML_DEVICE_ON_SAME_BOARD         "nvmlDeviceOnSameBoard"

typedef enum{
    ixmlEventTypeXidCriticalError = nvmlEventTypeXidCriticalError
} ixmlEventTypes;

typedef nvmlEventData_t ixmlEventData_t;

nvmlReturn_t dl_init();
nvmlReturn_t dl_close();
nvmlReturn_t ixmlInit();
nvmlReturn_t ixmlShutdown();
nvmlReturn_t ixmlDeviceGetCount(unsigned int* deviceCount);
nvmlReturn_t ixmlSystemGetDriverVersion(char *version, unsigned int length);
nvmlReturn_t ixmlSystemGetCudaDriverVersion(int *version);
nvmlReturn_t ixmlDeviceGetHandleByIndex(unsigned int index, nvmlDevice_t* device);
nvmlReturn_t ixmlDeviceGetHandleByUUID(const char *uuid, nvmlDevice_t* device);
nvmlReturn_t ixmlDeviceGetName(nvmlDevice_t device, char* name, unsigned int length);
nvmlReturn_t ixmlDeviceGetIndex(nvmlDevice_t device, unsigned int *index);
nvmlReturn_t ixmlDeviceGetUUID(nvmlDevice_t device, char* uuid, unsigned int length);
nvmlReturn_t ixmlDeviceGetFanSpeed(nvmlDevice_t device, unsigned int* speed);
nvmlReturn_t ixmlDeviceGetMemoryInfo(nvmlDevice_t device, nvmlMemory_t* memory);
nvmlReturn_t ixmlDeviceGetTemperature(nvmlDevice_t device, nvmlTemperatureSensors_t sensorType, unsigned int* temp);
nvmlReturn_t ixmlDeviceGetPciInfo(nvmlDevice_t device, nvmlPciInfo_t* pci);
nvmlReturn_t ixmlDeviceGetMinorNumber(nvmlDevice_t device, unsigned int* minorNumber);
nvmlReturn_t ixmlDeviceGetPowerUsage(nvmlDevice_t device, unsigned int* power);
nvmlReturn_t ixmlDeviceGetPowerManagementLimitConstraints(nvmlDevice_t device, unsigned int* minLimit, unsigned int* maxLimit);
nvmlReturn_t ixmlDeviceGetClockInfo(nvmlDevice_t device, nvmlClockType_t type, unsigned int* clock);
nvmlReturn_t ixmlDeviceGetUtilizationRates(nvmlDevice_t device, nvmlUtilization_t* utilization);
nvmlReturn_t ixmlEventSetCreate(nvmlEventSet_t *set);
nvmlReturn_t ixmlEventSetFree(nvmlEventSet_t set);
nvmlReturn_t ixmlDeviceRegisterEvents(nvmlDevice_t device, unsigned long long eventTypes, nvmlEventSet_t set);
nvmlReturn_t ixmlEventSetWait(nvmlEventSet_t set, ixmlEventData_t * data, unsigned int timeoutms);
nvmlReturn_t ixmlDeviceOnSameBoard(nvmlDevice_t device1, nvmlDevice_t device2, int* onSameBoard);
