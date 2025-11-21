/*
 * Copyright 1993-2022 NVIDIA Corporation.  All rights reserved.
 *
 * NOTICE TO USER:
 *
 * This source code is subject to NVIDIA ownership rights under U.S. and
 * international Copyright laws.  Users and possessors of this source code
 * are hereby granted a nonexclusive, royalty-free license to use this code
 * in individual and commercial software.
 *
 * NVIDIA MAKES NO REPRESENTATION ABOUT THE SUITABILITY OF THIS SOURCE
 * CODE FOR ANY PURPOSE.  IT IS PROVIDED "AS IS" WITHOUT EXPRESS OR
 * IMPLIED WARRANTY OF ANY KIND.  NVIDIA DISCLAIMS ALL WARRANTIES WITH
 * REGARD TO THIS SOURCE CODE, INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY, NONINFRINGEMENT, AND FITNESS FOR A PARTICULAR PURPOSE.
 * IN NO EVENT SHALL NVIDIA BE LIABLE FOR ANY SPECIAL, INDIRECT, INCIDENTAL,
 * OR CONSEQUENTIAL DAMAGES, OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS
 * OF USE, DATA OR PROFITS,  WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE
 * OR OTHER TORTIOUS ACTION,  ARISING OUT OF OR IN CONNECTION WITH THE USE
 * OR PERFORMANCE OF THIS SOURCE CODE.
 *
 * U.S. Government End Users.   This source code is a "commercial item" as
 * that term is defined at  48 C.F.R. 2.101 (OCT 1995), consisting  of
 * "commercial computer  software"  and "commercial computer software
 * documentation" as such terms are  used in 48 C.F.R. 12.212 (SEPT 1995)
 * and is provided to the U.S. Government only as a commercial end item.
 * Consistent with 48 C.F.R.12.212 and 48 C.F.R. 227.7202-1 through
 * 227.7202-4 (JUNE 1995), all U.S. Government End Users acquire the
 * source code with only those rights set forth herein.
 *
 * Any use of this source code in individual and commercial software must
 * include, in the user documentation and internal comments to the code,
 * the above Disclaimer and U.S. Government End Users Notice.
 */

/*
NVML API Reference

The NVIDIA Management Library (NVML) is a C-based programmatic interface for monitoring and
managing various states within NVIDIA Tesla &tm; GPUs. It is intended to be a platform for building
3rd party applications, and is also the underlying library for the NVIDIA-supported nvidia-smi
tool. NVML is thread-safe so it is safe to make simultaneous NVML calls from multiple threads.

API Documentation

Supported platforms:
- Windows:     Windows Server 2008 R2 64bit, Windows Server 2012 R2 64bit, Windows 7 64bit, Windows 8 64bit, Windows 10 64bit
- Linux:       32-bit and 64-bit
- Hypervisors: Windows Server 2008R2/2012 Hyper-V 64bit, Citrix XenServer 6.2 SP1+, VMware ESX 5.1/5.5

Supported products:
- Full Support
    - All Tesla products, starting with the Fermi architecture
    - All Quadro products, starting with the Fermi architecture
    - All vGPU Software products, starting with the Kepler architecture
    - Selected GeForce Titan products
- Limited Support
    - All Geforce products, starting with the Fermi architecture

The NVML library can be found at \%ProgramW6432\%\\"NVIDIA Corporation"\\NVSMI\\ on Windows. It is
not be added to the system path by default. To dynamically link to NVML, add this path to the PATH
environmental variable. To dynamically load NVML, call LoadLibrary with this path.

On Linux the NVML library will be found on the standard library path. For 64 bit Linux, both the 32 bit
and 64 bit NVML libraries will be installed.

Online documentation for this library is available at http://docs.nvidia.com/deploy/nvml-api/index.html
*/

#ifndef __nvml_nvml_h__
#define __nvml_nvml_h__

#ifdef __cplusplus
extern "C"
{
#endif

/*
 * On Windows, set up methods for DLL export
 * define NVML_STATIC_IMPORT when using nvml_loader library
 */
#if defined _WINDOWS
#if !defined NVML_STATIC_IMPORT
#if defined NVML_LIB_EXPORT
#define DECLDIR __declspec(dllexport)
#else
#define DECLDIR __declspec(dllimport)
#endif
#else
#define DECLDIR
#endif
#else
#define DECLDIR
#endif

#define DEVICE_MAX_NUM 32

    /**
     * Return values for NVML API calls.
     */
    typedef enum nvmlReturn_enum
    {
        // cppcheck-suppress *
        NVML_SUCCESS = 0,                          //!< The operation was successful
        NVML_ERROR_UNINITIALIZED = 1,              //!< NVML was not first initialized with nvmlInit()
        NVML_ERROR_INVALID_ARGUMENT = 2,           //!< A supplied argument is invalid
        NVML_ERROR_NOT_SUPPORTED = 3,              //!< The requested operation is not available on target device
        NVML_ERROR_NO_PERMISSION = 4,              //!< The current user does not have permission for operation
        NVML_ERROR_ALREADY_INITIALIZED = 5,        //!< Deprecated: Multiple initializations are now allowed through ref counting
        NVML_ERROR_NOT_FOUND = 6,                  //!< A query to find an object was unsuccessful
        NVML_ERROR_INSUFFICIENT_SIZE = 7,          //!< An input argument is not large enough
        NVML_ERROR_INSUFFICIENT_POWER = 8,         //!< A device's external power cables are not properly attached
        NVML_ERROR_DRIVER_NOT_LOADED = 9,          //!< NVIDIA driver is not loaded
        NVML_ERROR_TIMEOUT = 10,                   //!< User provided timeout passed
        NVML_ERROR_IRQ_ISSUE = 11,                 //!< NVIDIA Kernel detected an interrupt issue with a GPU
        NVML_ERROR_LIBRARY_NOT_FOUND = 12,         //!< NVML Shared Library couldn't be found or loaded
        NVML_ERROR_FUNCTION_NOT_FOUND = 13,        //!< Local version of NVML doesn't implement this function
        NVML_ERROR_CORRUPTED_INFOROM = 14,         //!< infoROM is corrupted
        NVML_ERROR_GPU_IS_LOST = 15,               //!< The GPU has fallen off the bus or has otherwise become inaccessible
        NVML_ERROR_RESET_REQUIRED = 16,            //!< The GPU requires a reset before it can be used again
        NVML_ERROR_OPERATING_SYSTEM = 17,          //!< The GPU control device has been blocked by the operating system/cgroups
        NVML_ERROR_LIB_RM_VERSION_MISMATCH = 18,   //!< RM detects a driver/library version mismatch
        NVML_ERROR_IN_USE = 19,                    //!< An operation cannot be performed because the GPU is currently in use
        NVML_ERROR_MEMORY = 20,                    //!< Insufficient memory
        NVML_ERROR_NO_DATA = 21,                   //!< No data
        NVML_ERROR_VGPU_ECC_NOT_SUPPORTED = 22,    //!< The requested vgpu operation is not available on target device, becasue ECC is enabled
        NVML_ERROR_INSUFFICIENT_RESOURCES = 23,    //!< Ran out of critical resources, other than memory
        NVML_ERROR_FREQ_NOT_SUPPORTED = 24,        //!< Ran out of critical resources, other than memory
        NVML_ERROR_ARGUMENT_VERSION_MISMATCH = 25, //!< The provided version is invalid/unsupported
        NVML_ERROR_DEPRECATED = 26,                //!< The requested functionality has been deprecated
        NVML_ERROR_UNKNOWN = 999                   //!< An internal driver error occurred
    } nvmlReturn_t;

    typedef struct
    {
        struct nvmlDevice_st *handle;
    } nvmlDevice_t;

    /**
     * See \ref nvmlDeviceGetMemoryErrorCounter
     */
    typedef enum nvmlMemoryLocation_enum
    {
        NVML_MEMORY_LOCATION_L1_CACHE = 0,       //!< GPU L1 Cache
        NVML_MEMORY_LOCATION_L2_CACHE = 1,       //!< GPU L2 Cache
        NVML_MEMORY_LOCATION_DRAM = 2,           //!< Turing+ DRAM
        NVML_MEMORY_LOCATION_DEVICE_MEMORY = 2,  //!< GPU Device Memory
        NVML_MEMORY_LOCATION_REGISTER_FILE = 3,  //!< GPU Register File
        NVML_MEMORY_LOCATION_TEXTURE_MEMORY = 4, //!< GPU Texture Memory
        NVML_MEMORY_LOCATION_TEXTURE_SHM = 5,    //!< Shared memory
        NVML_MEMORY_LOCATION_CBU = 6,            //!< CBU
        NVML_MEMORY_LOCATION_SRAM = 7,           //!< Turing+ SRAM
        // Keep this last
        NVML_MEMORY_LOCATION_COUNT //!< This counts the number of memory locations the driver knows about
    } nvmlMemoryLocation_t;

    /**
     * Causes for page retirement
     */
    typedef enum nvmlPageRetirementCause_enum
    {
        NVML_PAGE_RETIREMENT_CAUSE_MULTIPLE_SINGLE_BIT_ECC_ERRORS = 0, //!< Page was retired due to multiple single bit ECC error
        NVML_PAGE_RETIREMENT_CAUSE_DOUBLE_BIT_ECC_ERROR = 1,           //!< Page was retired due to double bit ECC error

        // Keep this last
        NVML_PAGE_RETIREMENT_CAUSE_COUNT
    } nvmlPageRetirementCause_t;

    /**
     * API types that allow changes to default permission restrictions
     */
    typedef enum nvmlRestrictedAPI_enum
    {
        NVML_RESTRICTED_API_SET_APPLICATION_CLOCKS = 0,  //!< APIs that change application clocks, see nvmlDeviceSetApplicationsClocks
                                                         //!< and see nvmlDeviceResetApplicationsClocks
        NVML_RESTRICTED_API_SET_AUTO_BOOSTED_CLOCKS = 1, //!< APIs that enable/disable Auto Boosted clocks
                                                         //!< see nvmlDeviceSetAutoBoostedClocksEnabled
        // Keep this last
        NVML_RESTRICTED_API_COUNT
    } nvmlRestrictedAPI_t;

    /** @} */
    /**
     * Memory allocation information for a device (v1).
     * The total amount is equal to the sum of the amounts of free and used memory.
     */
    typedef struct nvmlMemory_st
    {
        unsigned long long total; //!< Total physical device memory (in bytes)
        unsigned long long free;  //!< Unallocated device memory (in bytes)
        unsigned long long used;  //!< Sum of Reserved and Allocated device memory (in bytes).
                                  //!< Note that the driver/GPU always sets aside a small amount of memory for bookkeeping
    } nvmlMemory_t;

    /**
     * Memory allocation information for a device (v2).
     *
     * Version 2 adds versioning for the struct and the amount of system-reserved memory as an output.
     * @note The \ref nvmlMemory_v2_t.used amount also includes the \ref nvmlMemory_v2_t.reserved amount.
     */
    typedef struct nvmlMemory_v2_st
    {
        unsigned int version;        //!< Structure format version (must be 2)
        unsigned long long total;    //!< Total physical device memory (in bytes)
        unsigned long long reserved; //!< Device memory (in bytes) reserved for system use (driver or firmware)
        unsigned long long free;     //!< Unallocated device memory (in bytes)
        unsigned long long used;     //!< Allocated device memory (in bytes). Note that the driver/GPU always sets aside a small amount of memory for bookkeeping
    } nvmlMemory_v2_t;

    /**
     * Temperature thresholds.
     */
    typedef enum nvmlTemperatureThresholds_enum
    {
        NVML_TEMPERATURE_THRESHOLD_SHUTDOWN = 0,      // Temperature at which the GPU will
                                                      // shut down for HW protection
        NVML_TEMPERATURE_THRESHOLD_SLOWDOWN = 1,      // Temperature at which the GPU will
                                                      // begin HW slowdown
        NVML_TEMPERATURE_THRESHOLD_MEM_MAX = 2,       // Memory Temperature at which the GPU will
                                                      // begin SW slowdown
        NVML_TEMPERATURE_THRESHOLD_GPU_MAX = 3,       // GPU Temperature at which the GPU
                                                      // can be throttled below base clock
        NVML_TEMPERATURE_THRESHOLD_ACOUSTIC_MIN = 4,  // Minimum GPU Temperature that can be
                                                      // set as acoustic threshold
        NVML_TEMPERATURE_THRESHOLD_ACOUSTIC_CURR = 5, // Current temperature that is set as
                                                      // acoustic threshold.
        NVML_TEMPERATURE_THRESHOLD_ACOUSTIC_MAX = 6,  // Maximum GPU temperature that can be
                                                      // set as acoustic threshold.
        // Keep this last
        NVML_TEMPERATURE_THRESHOLD_COUNT
    } nvmlTemperatureThresholds_t;

    /**
     * Temperature sensors.
     */
    typedef enum nvmlTemperatureSensors_enum
    {
        NVML_TEMPERATURE_GPU = 0, //!< Temperature sensor for the GPU die

        // Keep this last
        NVML_TEMPERATURE_COUNT
    } nvmlTemperatureSensors_t;

    /**
     * Generic enable/disable enum.
     */
    typedef enum nvmlEnableState_enum
    {
        NVML_FEATURE_DISABLED = 0, //!< Feature disabled
        NVML_FEATURE_ENABLED = 1   //!< Feature enabled
    } nvmlEnableState_t;

    /**
     * Clock types.
     *
     * All speeds are in Mhz.
     */
    typedef enum nvmlClockType_enum
    {
        NVML_CLOCK_GRAPHICS = 0, //!< Graphics clock domain
        NVML_CLOCK_SM = 1,       //!< SM clock domain
        NVML_CLOCK_MEM = 2,      //!< Memory clock domain
        NVML_CLOCK_VIDEO = 3,    //!< Video encoder/decoder clock domain

        // Keep this last
        NVML_CLOCK_COUNT //!< Count of clock types
    } nvmlClockType_t;

    /**
     * Utilization information for a device.
     * Each sample period may be between 1 second and 1/6 second, depending on the product being queried.
     */
    typedef struct nvmlUtilization_st
    {
        unsigned int gpu;    //!< Percent of time over the past sample period during which one or more kernels was executing on the GPU
        unsigned int memory; //!< Percent of time over the past sample period during which global (device) memory was being read or written
    } nvmlUtilization_t;

    typedef struct nvmlProcessInfo_st
    {
        unsigned int pid;                            //!< Process ID
        unsigned long long usedGpuMemory;            //!< Amount of used GPU memory in bytes.
                                                     //! Under WDDM, \ref NVML_VALUE_NOT_AVAILABLE is always reported
                                                     //! because Windows KMD manages all the memory and not the NVIDIA driver
        unsigned int gpuInstanceId;                  //!< If MIG is enabled, stores a valid GPU instance ID. gpuInstanceId is set to
                                                     //  0xFFFFFFFF otherwise.
        unsigned int computeInstanceId;              //!< If MIG is enabled, stores a valid compute instance ID. computeInstanceId is set to
                                                     //  0xFFFFFFFF otherwise.
        unsigned long long usedGpuCcProtectedMemory; //!< Amount of used GPU conf compute protected memory in bytes.
    } nvmlProcessInfo_t;

    /**
     * Represents level relationships within a system between two GPUs
     * The enums are spaced to allow for future relationships
     */
    typedef enum nvmlGpuLevel_enum
    {
        NVML_TOPOLOGY_INTERNAL = 0,    // e.g. Tesla K80
        NVML_TOPOLOGY_SINGLE = 10,     // all devices that only need traverse a single PCIe switch
        NVML_TOPOLOGY_MULTIPLE = 20,   // all devices that need not traverse a host bridge
        NVML_TOPOLOGY_HOSTBRIDGE = 30, // all devices that are connected to the same host bridge
        NVML_TOPOLOGY_NODE = 40,       // all devices that are connected to the same NUMA node but possibly multiple host bridges
        NVML_TOPOLOGY_SYSTEM = 50      // all devices in the system

        // there is purposefully no COUNT here because of the need for spacing above
    } nvmlGpuTopologyLevel_t;

    /**
     * Information about running compute processes on the GPU, legacy version
     * for older versions of the API.
     */
    typedef struct nvmlProcessInfo_v1_st
    {
        unsigned int pid;                 //!< Process ID
        unsigned long long usedGpuMemory; //!< Amount of used GPU memory in bytes.
                                          //! Under WDDM, \ref NVML_VALUE_NOT_AVAILABLE is always reported
                                          //! because Windows KMD manages all the memory and not the NVIDIA driver
    } nvmlProcessInfo_v1_t;

    /**
     * GPM Metric Identifiers
     */
    typedef enum
    {
        NVML_GPM_METRIC_GRAPHICS_UTIL = 1,            //!< Percentage of time any compute/graphics app was active on the GPU. 0.0 - 100.0
        NVML_GPM_METRIC_SM_UTIL = 2,                  //!< Percentage of SMs that were busy. 0.0 - 100.0
        NVML_GPM_METRIC_SM_OCCUPANCY = 3,             //!< Percentage of warps that were active vs theoretical maximum. 0.0 - 100.0
        NVML_GPM_METRIC_INTEGER_UTIL = 4,             //!< Percentage of time the GPU's SMs were doing integer operations. 0.0 - 100.0
        NVML_GPM_METRIC_ANY_TENSOR_UTIL = 5,          //!< Percentage of time the GPU's SMs were doing ANY tensor operations. 0.0 - 100.0
        NVML_GPM_METRIC_DFMA_TENSOR_UTIL = 6,         //!< Percentage of time the GPU's SMs were doing DFMA tensor operations. 0.0 - 100.0
        NVML_GPM_METRIC_HMMA_TENSOR_UTIL = 7,         //!< Percentage of time the GPU's SMs were doing HMMA tensor operations. 0.0 - 100.0
        NVML_GPM_METRIC_IMMA_TENSOR_UTIL = 9,         //!< Percentage of time the GPU's SMs were doing IMMA tensor operations. 0.0 - 100.0
        NVML_GPM_METRIC_DRAM_BW_UTIL = 10,            //!< Percentage of DRAM bw used vs theoretical maximum. 0.0 - 100.0 */
        NVML_GPM_METRIC_FP64_UTIL = 11,               //!< Percentage of time the GPU's SMs were doing non-tensor FP64 math. 0.0 - 100.0
        NVML_GPM_METRIC_FP32_UTIL = 12,               //!< Percentage of time the GPU's SMs were doing non-tensor FP32 math. 0.0 - 100.0
        NVML_GPM_METRIC_FP16_UTIL = 13,               //!< Percentage of time the GPU's SMs were doing non-tensor FP16 math. 0.0 - 100.0
        NVML_GPM_METRIC_PCIE_TX_PER_SEC = 20,         //!< PCIe traffic from this GPU in MiB/sec
        NVML_GPM_METRIC_PCIE_RX_PER_SEC = 21,         //!< PCIe traffic to this GPU in MiB/sec
        NVML_GPM_METRIC_NVDEC_0_UTIL = 30,            //!< Percent utilization of NVDEC 0. 0.0 - 100.0
        NVML_GPM_METRIC_NVDEC_1_UTIL = 31,            //!< Percent utilization of NVDEC 1. 0.0 - 100.0
        NVML_GPM_METRIC_NVDEC_2_UTIL = 32,            //!< Percent utilization of NVDEC 2. 0.0 - 100.0
        NVML_GPM_METRIC_NVDEC_3_UTIL = 33,            //!< Percent utilization of NVDEC 3. 0.0 - 100.0
        NVML_GPM_METRIC_NVDEC_4_UTIL = 34,            //!< Percent utilization of NVDEC 4. 0.0 - 100.0
        NVML_GPM_METRIC_NVDEC_5_UTIL = 35,            //!< Percent utilization of NVDEC 5. 0.0 - 100.0
        NVML_GPM_METRIC_NVDEC_6_UTIL = 36,            //!< Percent utilization of NVDEC 6. 0.0 - 100.0
        NVML_GPM_METRIC_NVDEC_7_UTIL = 37,            //!< Percent utilization of NVDEC 7. 0.0 - 100.0
        NVML_GPM_METRIC_NVJPG_0_UTIL = 40,            //!< Percent utilization of NVJPG 0. 0.0 - 100.0
        NVML_GPM_METRIC_NVJPG_1_UTIL = 41,            //!< Percent utilization of NVJPG 1. 0.0 - 100.0
        NVML_GPM_METRIC_NVJPG_2_UTIL = 42,            //!< Percent utilization of NVJPG 2. 0.0 - 100.0
        NVML_GPM_METRIC_NVJPG_3_UTIL = 43,            //!< Percent utilization of NVJPG 3. 0.0 - 100.0
        NVML_GPM_METRIC_NVJPG_4_UTIL = 44,            //!< Percent utilization of NVJPG 4. 0.0 - 100.0
        NVML_GPM_METRIC_NVJPG_5_UTIL = 45,            //!< Percent utilization of NVJPG 5. 0.0 - 100.0
        NVML_GPM_METRIC_NVJPG_6_UTIL = 46,            //!< Percent utilization of NVJPG 6. 0.0 - 100.0
        NVML_GPM_METRIC_NVJPG_7_UTIL = 47,            //!< Percent utilization of NVJPG 7. 0.0 - 100.0
        NVML_GPM_METRIC_NVOFA_0_UTIL = 50,            //!< Percent utilization of NVOFA 0. 0.0 - 100.0
        NVML_GPM_METRIC_NVLINK_TOTAL_RX_PER_SEC = 60, //!< NvLink read bandwidth for all links in MiB/sec
        NVML_GPM_METRIC_NVLINK_TOTAL_TX_PER_SEC = 61, //!< NvLink write bandwidth for all links in MiB/sec
        NVML_GPM_METRIC_NVLINK_L0_RX_PER_SEC = 62,    //!< NvLink read bandwidth for link 0 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L0_TX_PER_SEC = 63,    //!< NvLink write bandwidth for link 0 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L1_RX_PER_SEC = 64,    //!< NvLink read bandwidth for link 1 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L1_TX_PER_SEC = 65,    //!< NvLink write bandwidth for link 1 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L2_RX_PER_SEC = 66,    //!< NvLink read bandwidth for link 2 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L2_TX_PER_SEC = 67,    //!< NvLink write bandwidth for link 2 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L3_RX_PER_SEC = 68,    //!< NvLink read bandwidth for link 3 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L3_TX_PER_SEC = 69,    //!< NvLink write bandwidth for link 3 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L4_RX_PER_SEC = 70,    //!< NvLink read bandwidth for link 4 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L4_TX_PER_SEC = 71,    //!< NvLink write bandwidth for link 4 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L5_RX_PER_SEC = 72,    //!< NvLink read bandwidth for link 5 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L5_TX_PER_SEC = 73,    //!< NvLink write bandwidth for link 5 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L6_RX_PER_SEC = 74,    //!< NvLink read bandwidth for link 6 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L6_TX_PER_SEC = 75,    //!< NvLink write bandwidth for link 6 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L7_RX_PER_SEC = 76,    //!< NvLink read bandwidth for link 7 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L7_TX_PER_SEC = 77,    //!< NvLink write bandwidth for link 7 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L8_RX_PER_SEC = 78,    //!< NvLink read bandwidth for link 8 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L8_TX_PER_SEC = 79,    //!< NvLink write bandwidth for link 8 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L9_RX_PER_SEC = 80,    //!< NvLink read bandwidth for link 9 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L9_TX_PER_SEC = 81,    //!< NvLink write bandwidth for link 9 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L10_RX_PER_SEC = 82,   //!< NvLink read bandwidth for link 10 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L10_TX_PER_SEC = 83,   //!< NvLink write bandwidth for link 10 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L11_RX_PER_SEC = 84,   //!< NvLink read bandwidth for link 11 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L11_TX_PER_SEC = 85,   //!< NvLink write bandwidth for link 11 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L12_RX_PER_SEC = 86,   //!< NvLink read bandwidth for link 12 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L12_TX_PER_SEC = 87,   //!< NvLink write bandwidth for link 12 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L13_RX_PER_SEC = 88,   //!< NvLink read bandwidth for link 13 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L13_TX_PER_SEC = 89,   //!< NvLink write bandwidth for link 13 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L14_RX_PER_SEC = 90,   //!< NvLink read bandwidth for link 14 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L14_TX_PER_SEC = 91,   //!< NvLink write bandwidth for link 14 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L15_RX_PER_SEC = 92,   //!< NvLink read bandwidth for link 15 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L15_TX_PER_SEC = 93,   //!< NvLink write bandwidth for link 15 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L16_RX_PER_SEC = 94,   //!< NvLink read bandwidth for link 16 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L16_TX_PER_SEC = 95,   //!< NvLink write bandwidth for link 16 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L17_RX_PER_SEC = 96,   //!< NvLink read bandwidth for link 17 in MiB/sec
        NVML_GPM_METRIC_NVLINK_L17_TX_PER_SEC = 97,   //!< NvLink write bandwidth for link 17 in MiB/sec
        NVML_GPM_METRIC_MAX = 98,                     //!< Maximum value above +1. Note that changing this should also change NVML_GPM_METRICS_GET_VERSION due to struct size change
    } nvmlGpmMetricId_t;

    /** @} */ // @defgroup nvmlGpmEnums

    /***************************************************************************************************/
    /** @defgroup nvmlGpmStructs GPM Structs
     *  @{
     */
    /***************************************************************************************************/

    /**
     * Handle to an allocated GPM sample allocated with nvmlGpmSampleAlloc(). Free this with nvmlGpmSampleFree().
     */
    typedef struct
    {
        struct nvmlGpmSample_st *handle;
    } nvmlGpmSample_t;

    /**
     * GPM metric information.
     */
    typedef struct
    {
        unsigned int metricId;   //!<  IN: NVML_GPM_METRIC_? #define of which metric to retrieve
        nvmlReturn_t nvmlReturn; //!<  OUT: Status of this metric. If this is nonzero, then value is not valid
        double value;            //!<  OUT: Value of this metric. Is only valid if nvmlReturn is 0 (NVML_SUCCESS)
        struct
        {
            char *shortName;
            char *longName;
            char *unit;
        } metricInfo; //!< OUT: Metric name and unit. Those can be NULL if not defined
    } nvmlGpmMetric_t;

    /**
     * GPM buffer information.
     */
    typedef struct
    {
        unsigned int version;                         //!< IN: Set to NVML_GPM_METRICS_GET_VERSION
        unsigned int numMetrics;                      //!< IN: How many metrics to retrieve in metrics[]
        nvmlGpmSample_t sample1;                      //!< IN: Sample buffer
        nvmlGpmSample_t sample2;                      //!< IN: Sample buffer
        nvmlGpmMetric_t metrics[NVML_GPM_METRIC_MAX]; //!< IN/OUT: Array of metrics. Set metricId on call. See nvmlReturn and value on return
    } nvmlGpmMetricsGet_t;

#define NVML_GPM_METRICS_GET_VERSION 1

    /**
     * GPM device information.
     */
    typedef struct
    {
        unsigned int version;           //!< IN: Set to NVML_GPM_SUPPORT_VERSION
        unsigned int isSupportedDevice; //!< OUT: Indicates device support
    } nvmlGpmSupport_t;

/**
 * Maximum limit on Physical Bridges per Board
 */
#define NVML_MAX_PHYSICAL_BRIDGE (128)

    /**
     * Enum to represent type of bridge chip
     */
    typedef enum nvmlBridgeChipType_enum
    {
        NVML_BRIDGE_CHIP_PLX = 0,
        NVML_BRIDGE_CHIP_BRO4 = 1
    } nvmlBridgeChipType_t;

    /**
     * Information about the Bridge Chip Firmware
     */
    typedef struct nvmlBridgeChipInfo_st
    {
        nvmlBridgeChipType_t type; //!< Type of Bridge Chip
        unsigned int fwVersion;    //!< Firmware Version. 0=Version is unavailable
    } nvmlBridgeChipInfo_t;

    /**
     * This structure stores the complete Hierarchy of the Bridge Chip within the board. The immediate
     * bridge is stored at index 0 of bridgeInfoList, parent to immediate bridge is at index 1 and so forth.
     */
    typedef struct nvmlBridgeChipHierarchy_st
    {
        unsigned char bridgeCount;                                     //!< Number of Bridge Chips on the Board
        nvmlBridgeChipInfo_t bridgeChipInfo[NVML_MAX_PHYSICAL_BRIDGE]; //!< Hierarchy of Bridge Chips on the board
    } nvmlBridgeChipHierarchy_t;

    /**
     *  Represents Type of Sampling Event
     */
    typedef enum nvmlSamplingType_enum
    {
        NVML_TOTAL_POWER_SAMPLES = 0, //!< To represent total power drawn by GPU
        NVML_GPU_UTILIZATION_SAMPLES =
            1, //!< To represent percent of time during which one or more kernels was executing on the GPU
        NVML_MEMORY_UTILIZATION_SAMPLES =
            2,                            //!< To represent percent of time during which global (device) memory was being read or written
        NVML_ENC_UTILIZATION_SAMPLES = 3, //!< To represent percent of time during which NVENC remains busy
        NVML_DEC_UTILIZATION_SAMPLES = 4, //!< To represent percent of time during which NVDEC remains busy
        NVML_PROCESSOR_CLK_SAMPLES = 5,   //!< To represent processor clock samples
        NVML_MEMORY_CLK_SAMPLES = 6,      //!< To represent memory clock samples

        // Keep this last
        NVML_SAMPLINGTYPE_COUNT
    } nvmlSamplingType_t;

    /**
     * Represents the queryable PCIe utilization counters
     */
    typedef enum nvmlPcieUtilCounter_enum
    {
        NVML_PCIE_UTIL_TX_BYTES = 0, // 1KB granularity
        NVML_PCIE_UTIL_RX_BYTES = 1, // 1KB granularity

        // Keep this last
        NVML_PCIE_UTIL_COUNT
    } nvmlPcieUtilCounter_t;

    /**
     * Represents the type for sample value returned
     */
    typedef enum nvmlValueType_enum
    {
        NVML_VALUE_TYPE_DOUBLE = 0,
        NVML_VALUE_TYPE_UNSIGNED_INT = 1,
        NVML_VALUE_TYPE_UNSIGNED_LONG = 2,
        NVML_VALUE_TYPE_UNSIGNED_LONG_LONG = 3,
        NVML_VALUE_TYPE_SIGNED_LONG_LONG = 4,

        // Keep this last
        NVML_VALUE_TYPE_COUNT
    } nvmlValueType_t;

    /**
     * Union to represent different types of Value
     */
    typedef union nvmlValue_st
    {
        double dVal;               //!< If the value is double
        unsigned int uiVal;        //!< If the value is unsigned int
        unsigned long ulVal;       //!< If the value is unsigned long
        unsigned long long ullVal; //!< If the value is unsigned long long
        signed long long sllVal;   //!< If the value is signed long long
    } nvmlValue_t;

    /**
     * Information for Sample
     */
    typedef struct nvmlSample_st
    {
        unsigned long long timeStamp; //!< CPU Timestamp in microseconds
        nvmlValue_t sampleValue;      //!< Sample Value
    } nvmlSample_t;

    /**
     * Represents type of perf policy for which violation times can be queried
     */
    typedef enum nvmlPerfPolicyType_enum
    {
        NVML_PERF_POLICY_POWER = 0,      //!< How long did power violations cause the GPU to be below application clocks
        NVML_PERF_POLICY_THERMAL = 1,    //!< How long did thermal violations cause the GPU to be below application clocks
        NVML_PERF_POLICY_SYNC_BOOST = 2, //!< How long did sync boost cause the GPU to be below application clocks
        NVML_PERF_POLICY_BOARD_LIMIT =
            3, //!< How long did the board limit cause the GPU to be below application clocks
        NVML_PERF_POLICY_LOW_UTILIZATION =
            4, //!< How long did low utilization cause the GPU to be below application clocks
        NVML_PERF_POLICY_RELIABILITY =
            5, //!< How long did the board reliability limit cause the GPU to be below application clocks

        NVML_PERF_POLICY_TOTAL_APP_CLOCKS =
            10,                                  //!< Total time the GPU was held below application clocks by any limiter (0 - 5 above)
        NVML_PERF_POLICY_TOTAL_BASE_CLOCKS = 11, //!< Total time the GPU was held below base clocks

        // Keep this last
        NVML_PERF_POLICY_COUNT
    } nvmlPerfPolicyType_t;

    /**
     * Compute mode.
     *
     * NVML_COMPUTEMODE_EXCLUSIVE_PROCESS was added in CUDA 4.0.
     * Earlier CUDA versions supported a single exclusive mode,
     * which is equivalent to NVML_COMPUTEMODE_EXCLUSIVE_THREAD in CUDA 4.0 and beyond.
     */
    typedef enum nvmlComputeMode_enum
    {
        NVML_COMPUTEMODE_DEFAULT = 0,           //!< Default compute mode -- multiple contexts per device
        NVML_COMPUTEMODE_EXCLUSIVE_THREAD = 1,  //!< Support Removed
        NVML_COMPUTEMODE_PROHIBITED = 2,        //!< Compute-prohibited mode -- no contexts per device
        NVML_COMPUTEMODE_EXCLUSIVE_PROCESS = 3, //!< Compute-exclusive-process mode -- only one context per device,
                                                //!< usable from multiple threads at a time

        // Keep this last
        NVML_COMPUTEMODE_COUNT
    } nvmlComputeMode_t;

    /**
     * Struct to hold perf policy violation status data
     */
    typedef struct nvmlViolationTime_st
    {
        unsigned long long referenceTime; //!< referenceTime represents CPU timestamp in microseconds
        unsigned long long violationTime; //!< violationTime in Nanoseconds
    } nvmlViolationTime_t;

#define NVML_GPM_SUPPORT_VERSION 1
/**
 * Buffer size guaranteed to be large enough for storing GPU identifiers.
 */
#define NVML_DEVICE_UUID_BUFFER_SIZE 80

/**
 * Buffer size guaranteed to be large enough for \ref nvmlSystemGetDriverVersion
 */
#define NVML_SYSTEM_DRIVER_VERSION_BUFFER_SIZE 80

/**
 * Buffer size guaranteed to be large enough for \ref nvmlSystemGetNVMLVersion
 */
#define NVML_SYSTEM_NVML_VERSION_BUFFER_SIZE 80

/**
 * Buffer size guaranteed to be large enough for storing GPU device names.
 */
#define NVML_DEVICE_NAME_BUFFER_SIZE 64

/**
 * Buffer size guaranteed to be large enough for \ref nvmlDeviceGetName
 */
#define NVML_DEVICE_NAME_V2_BUFFER_SIZE 96

/**
 * Buffer size guaranteed to be large enough for \ref nvmlDeviceGetBoardPartNumber
 */
#define NVML_DEVICE_PART_NUMBER_BUFFER_SIZE 80

/**
 * Buffer size guaranteed to be large enough for \ref nvmlDeviceGetSerial
 */
#define NVML_DEVICE_SERIAL_BUFFER_SIZE 30

/**
 * Buffer size guaranteed to be large enough for \ref nvmlDeviceGetVbiosVersion
 */
#define NVML_DEVICE_VBIOS_VERSION_BUFFER_SIZE 32

/**
 * Buffer size guaranteed to be large enough for pci bus id
 */
#define NVML_DEVICE_PCI_BUS_ID_BUFFER_SIZE 32

/**
 * Buffer size guaranteed to be large enough for pci bus id for ::busIdLegacy
 */
#define NVML_DEVICE_PCI_BUS_ID_BUFFER_V2_SIZE 16

#define ixmlHealthSYSHUBError 0x0000000000000001LL
#define ixmlHealthMCError 0x0000000000000002LL
#define ixmlHealthOverTempError 0x0000000000000004LL
#define ixmlHealthOverVoltageError 0x0000000000000008LL
#define ixmlHealthECCError 0x0000000000000010LL
#define ixmlHealthMemoryError 0x0000000000000020LL
#define ixmlHealthPCIEError 0x0000000000000040LL
#define ixmlHealthOK 0x0000000000000000LL

    /**
     * PCI information about a GPU device.
     */
    typedef struct nvmlPciInfo_st
    {
        char busIdLegacy[NVML_DEVICE_PCI_BUS_ID_BUFFER_V2_SIZE]; //!< The legacy tuple domain:bus:device.function PCI identifier (&amp; NULL terminator)
        unsigned int domain;                                     //!< The PCI domain on which the device's bus resides, 0 to 0xffffffff
        unsigned int bus;                                        //!< The bus on which the device resides, 0 to 0xff
        unsigned int device;                                     //!< The device's id on the bus, 0 to 31
        unsigned int pciDeviceId;                                //!< The combined 16-bit device id and 16-bit vendor id

        // Added in NVML 2.285 API
        unsigned int pciSubSystemId; //!< The 32-bit Sub System Device ID

        char busId[NVML_DEVICE_PCI_BUS_ID_BUFFER_SIZE]; //!< The tuple domain:bus:device.function PCI identifier (&amp; NULL terminator)
    } nvmlPciInfo_t;

    /**
     * Initialize NVML, but don't initialize any GPUs yet.
     *
     * \note nvmlInit_v3 introduces a "flags" argument, that allows passing boolean values
     *       modifying the behaviour of nvmlInit().
     * \note In NVML 5.319 new nvmlInit_v2 has replaced nvmlInit"_v1" (default in NVML 4.304 and older) that
     *       did initialize all GPU devices in the system.
     *
     * This allows NVML to communicate with a GPU
     * when other GPUs in the system are unstable or in a bad state.  When using this API, GPUs are
     * discovered and initialized in nvmlDeviceGetHandleBy* functions instead.
     *
     * \note To contrast nvmlInit_v2 with nvmlInit"_v1", NVML 4.304 nvmlInit"_v1" will fail when any detected GPU is in
     *       a bad or unstable state.
     *
     * For all products.
     *
     * This method, should be called once before invoking any other methods in the library.
     * A reference count of the number of initializations is maintained.  Shutdown only occurs
     * when the reference count reaches zero.
     *
     * @return
     *         - \ref NVML_SUCCESS                   if NVML has been properly initialized
     *         - \ref NVML_ERROR_DRIVER_NOT_LOADED   if NVIDIA driver is not running
     *         - \ref NVML_ERROR_NO_PERMISSION       if NVML does not have permission to talk to the driver
     *         - \ref NVML_ERROR_UNKNOWN             on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlInit_v2(void);

    /**
     * Shut down NVML by releasing all GPU resources previously allocated with \ref nvmlInit_v2().
     *
     * For all products.
     *
     * This method should be called after NVML work is done, once for each call to \ref nvmlInit_v2()
     * A reference count of the number of initializations is maintained.  Shutdown only occurs
     * when the reference count reaches zero.  For backwards compatibility, no error is reported if
     * nvmlShutdown() is called more times than nvmlInit().
     *
     * @return
     *         - \ref NVML_SUCCESS                 if NVML has been properly shut down
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlShutdown(void);

    /**
     * Retrieves the number of compute devices in the system. A compute device is a single GPU.
     *
     * For all products.
     *
     * Note: New nvmlDeviceGetCount_v2 (default in NVML 5.319) returns count of all devices in the system
     *       even if nvmlDeviceGetHandleByIndex_v2 returns NVML_ERROR_NO_PERMISSION for such device.
     *       Update your code to handle this error, or use NVML 4.304 or older nvml header file.
     *       For backward binary compatibility reasons _v1 version of the API is still present in the shared
     *       library.
     *       Old _v1 version of nvmlDeviceGetCount doesn't count devices that NVML has no permission to talk to.
     *
     * @param deviceCount                          Reference in which to return the number of accessible devices
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a deviceCount has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a deviceCount is NULL
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetCount_v2(unsigned int *deviceCount);

    /**
     * Acquire the handle for a particular device, based on its index.
     *
     * For all products.
     *
     * Valid indices are derived from the \a accessibleDevices count returned by
     *   \ref nvmlDeviceGetCount_v2(). For example, if \a accessibleDevices is 2 the valid indices
     *   are 0 and 1, corresponding to GPU 0 and GPU 1.
     *
     * The order in which NVML enumerates devices has no guarantees of consistency between reboots. For that reason it
     *   is recommended that devices be looked up by their PCI ids or UUID. See
     *   \ref nvmlDeviceGetHandleByUUID() and \ref nvmlDeviceGetHandleByPciBusId_v2().
     *
     * Note: The NVML index may not correlate with other APIs, such as the CUDA device index.
     *
     * Starting from NVML 5, this API causes NVML to initialize the target GPU
     * NVML may initialize additional GPUs if:
     *  - The target GPU is an SLI slave
     *
     * Note: New nvmlDeviceGetCount_v2 (default in NVML 5.319) returns count of all devices in the system
     *       even if nvmlDeviceGetHandleByIndex_v2 returns NVML_ERROR_NO_PERMISSION for such device.
     *       Update your code to handle this error, or use NVML 4.304 or older nvml header file.
     *       For backward binary compatibility reasons _v1 version of the API is still present in the shared
     *       library.
     *       Old _v1 version of nvmlDeviceGetCount doesn't count devices that NVML has no permission to talk to.
     *
     *       This means that nvmlDeviceGetHandleByIndex_v2 and _v1 can return different devices for the same index.
     *       If you don't touch macros that map old (_v1) versions to _v2 versions at the top of the file you don't
     *       need to worry about that.
     *
     * @param index                                The index of the target GPU, >= 0 and < \a accessibleDevices
     * @param device                               Reference in which to return the device handle
     *
     * @return
     *         - \ref NVML_SUCCESS                  if \a device has been set
     *         - \ref NVML_ERROR_UNINITIALIZED      if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT   if \a index is invalid or \a device is NULL
     *         - \ref NVML_ERROR_INSUFFICIENT_POWER if any attached devices have improperly attached external power cables
     *         - \ref NVML_ERROR_NO_PERMISSION      if the user doesn't have permission to talk to this device
     *         - \ref NVML_ERROR_IRQ_ISSUE          if NVIDIA kernel detected an interrupt issue with the attached GPUs
     *         - \ref NVML_ERROR_GPU_IS_LOST        if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN            on any unexpected error
     *
     * @see nvmlDeviceGetIndex
     * @see nvmlDeviceGetCount
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetHandleByIndex_v2(unsigned int index, nvmlDevice_t *device);

    /**
     * Acquire the handle for a particular device, based on its globally unique immutable UUID associated with each device.
     *
     * For all products.
     *
     * @param uuid                                 The UUID of the target GPU or MIG instance
     * @param device                               Reference in which to return the device handle or MIG device handle
     *
     * Starting from NVML 5, this API causes NVML to initialize the target GPU
     * NVML may initialize additional GPUs as it searches for the target GPU
     *
     * @return
     *         - \ref NVML_SUCCESS                  if \a device has been set
     *         - \ref NVML_ERROR_UNINITIALIZED      if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT   if \a uuid is invalid or \a device is null
     *         - \ref NVML_ERROR_NOT_FOUND          if \a uuid does not match a valid device on the system
     *         - \ref NVML_ERROR_INSUFFICIENT_POWER if any attached devices have improperly attached external power cables
     *         - \ref NVML_ERROR_IRQ_ISSUE          if NVIDIA kernel detected an interrupt issue with the attached GPUs
     *         - \ref NVML_ERROR_GPU_IS_LOST        if any GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN            on any unexpected error
     *
     * @see nvmlDeviceGetUUID
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetHandleByUUID(const char *uuid, nvmlDevice_t *device);

    /**
     * Acquire the handle for a particular device, based on its PCI bus id.
     *
     * For all products.
     *
     * This value corresponds to the nvmlPciInfo_t::busId returned by \ref nvmlDeviceGetPciInfo_v3().
     *
     * Starting from NVML 5, this API causes NVML to initialize the target GPU
     * NVML may initialize additional GPUs if:
     *  - The target GPU is an SLI slave
     *
     * \note NVML 4.304 and older version of nvmlDeviceGetHandleByPciBusId"_v1" returns NVML_ERROR_NOT_FOUND
     *       instead of NVML_ERROR_NO_PERMISSION.
     *
     * @param pciBusId                             The PCI bus id of the target GPU
     * @param device                               Reference in which to return the device handle
     *
     * @return
     *         - \ref NVML_SUCCESS                  if \a device has been set
     *         - \ref NVML_ERROR_UNINITIALIZED      if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT   if \a pciBusId is invalid or \a device is NULL
     *         - \ref NVML_ERROR_NOT_FOUND          if \a pciBusId does not match a valid device on the system
     *         - \ref NVML_ERROR_INSUFFICIENT_POWER if the attached device has improperly attached external power cables
     *         - \ref NVML_ERROR_NO_PERMISSION      if the user doesn't have permission to talk to this device
     *         - \ref NVML_ERROR_IRQ_ISSUE          if NVIDIA kernel detected an interrupt issue with the attached GPUs
     *         - \ref NVML_ERROR_GPU_IS_LOST        if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN            on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetHandleByPciBusId_v2(const char *pciBusId, nvmlDevice_t *device);

    /**
     * Retrieves minor number for the device. The minor number for the device is such that the Nvidia device node file for
     * each GPU will have the form /dev/nvidia[minor number].
     *
     * For all products.
     * Supported only for Linux
     *
     * @param device                                The identifier of the target device
     * @param minorNumber                           Reference in which to return the minor number for the device
     * @return
     *         - \ref NVML_SUCCESS                 if the minor number is successfully retrieved
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a minorNumber is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if this query is not supported by the device
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetMinorNumber(nvmlDevice_t device, unsigned int *minorNumber);

    /**
     * Retrieves the globally unique immutable UUID associated with this device, as a 5 part hexadecimal string,
     * that augments the immutable, board serial identifier.
     *
     * For all products.
     *
     * The UUID is a globally unique identifier. It is the only available identifier for pre-Fermi-architecture products.
     * It does NOT correspond to any identifier printed on the board.  It will not exceed 96 characters in length
     * (including the NULL terminator).  See \ref nvmlConstants::NVML_DEVICE_UUID_V2_BUFFER_SIZE.
     *
     * When used with MIG device handles the API returns globally unique UUIDs which can be used to identify MIG
     * devices across both GPU and MIG devices. UUIDs are immutable for the lifetime of a MIG device.
     *
     * @param device                               The identifier of the target device
     * @param uuid                                 Reference in which to return the GPU UUID
     * @param length                               The maximum allowed length of the string returned in \a uuid
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a uuid has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid, or \a uuid is NULL
     *         - \ref NVML_ERROR_INSUFFICIENT_SIZE if \a length is too small
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support this feature
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetUUID(nvmlDevice_t device, char *uuid, unsigned int length);

    /**
     * Retrieves the name of this device.
     *
     * For all products.
     *
     * The name is an alphanumeric string that denotes a particular product, e.g. Tesla &tm; C2070. It will not
     * exceed 96 characters in length (including the NULL terminator).  See \ref
     * nvmlConstants::NVML_DEVICE_NAME_V2_BUFFER_SIZE.
     *
     * When used with MIG device handles the API returns MIG device names which can be used to identify devices
     * based on their attributes.
     *
     * @param device                               The identifier of the target device
     * @param name                                 Reference in which to return the product name
     * @param length                               The maximum allowed length of the string returned in \a name
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a name has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid, or \a name is NULL
     *         - \ref NVML_ERROR_INSUFFICIENT_SIZE if \a length is too small
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetName(nvmlDevice_t device, char *name, unsigned int length);

    /**
     * Retrieves the version of the system's graphics driver.
     *
     * For all products.
     *
     * The version identifier is an alphanumeric string.  It will not exceed 80 characters in length
     * (including the NULL terminator).  See \ref nvmlConstants::NVML_SYSTEM_DRIVER_VERSION_BUFFER_SIZE.
     *
     * @param version                              Reference in which to return the version identifier
     * @param length                               The maximum allowed length of the string returned in \a version
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a version has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a version is NULL
     *         - \ref NVML_ERROR_INSUFFICIENT_SIZE if \a length is too small
     */
    nvmlReturn_t DECLDIR nvmlSystemGetDriverVersion(char *version, unsigned int length);

    /**
     * Retrieves the version of the NVML library.
     *
     * For all products.
     *
     * The version identifier is an alphanumeric string.  It will not exceed 80 characters in length
     * (including the NULL terminator).  See \ref nvmlConstants::NVML_SYSTEM_NVML_VERSION_BUFFER_SIZE.
     *
     * @param version                              Reference in which to return the version identifier
     * @param length                               The maximum allowed length of the string returned in \a version
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a version has been set
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a version is NULL
     *         - \ref NVML_ERROR_INSUFFICIENT_SIZE if \a length is too small
     */
    nvmlReturn_t DECLDIR nvmlSystemGetNVMLVersion(char *version, unsigned int length);

    /**
     * Retrieves the version of the CUDA driver.
     *
     * For all products.
     *
     * The CUDA driver version returned will be retreived from the currently installed version of CUDA.
     * If the cuda library is not found, this function will return a known supported version number.
     *
     * @param cudaDriverVersion                    Reference in which to return the version identifier
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a cudaDriverVersion has been set
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a cudaDriverVersion is NULL
     */
    nvmlReturn_t DECLDIR nvmlSystemGetCudaDriverVersion(int *cudaDriverVersion);

    /**
     * Retrieves the version of the CUDA driver from the shared library.
     *
     * For all products.
     *
     * The returned CUDA driver version by calling cuDriverGetVersion()
     *
     * @param cudaDriverVersion                    Reference in which to return the version identifier
     *
     * @return
     *         - \ref NVML_SUCCESS                  if \a cudaDriverVersion has been set
     *         - \ref NVML_ERROR_INVALID_ARGUMENT   if \a cudaDriverVersion is NULL
     *         - \ref NVML_ERROR_LIBRARY_NOT_FOUND  if \a libcuda.so.1 or libcuda.dll is not found
     *         - \ref NVML_ERROR_FUNCTION_NOT_FOUND if \a cuDriverGetVersion() is not found in the shared library
     */
    nvmlReturn_t DECLDIR nvmlSystemGetCudaDriverVersion_v2(int *cudaDriverVersion);

    /**
     * Retrieves the current temperature readings for the device, in degrees C.
     *
     * For all products.
     *
     * See \ref nvmlTemperatureSensors_t for details on available temperature sensors.
     *
     * @param device                               The identifier of the target device
     * @param sensorType                           Flag that indicates which sensor reading to retrieve
     * @param temp                                 Reference in which to return the temperature reading
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a temp has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid, \a sensorType is invalid or \a temp is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not have the specified sensor
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetTemperature(nvmlDevice_t device, nvmlTemperatureSensors_t sensorType, unsigned int *temp);

    /**
     * Retrieves the temperature threshold for the GPU with the specified threshold type in degrees C.
     *
     * For Kepler &tm; or newer fully supported devices.
     *
     * See \ref nvmlTemperatureThresholds_t for details on available temperature thresholds.
     *
     * @param device                               The identifier of the target device
     * @param thresholdType                        The type of threshold value queried
     * @param temp                                 Reference in which to return the temperature reading
     * @return
     *         - \ref NVML_SUCCESS                 if \a temp has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid, \a thresholdType is invalid or \a temp is
     * NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not have a temperature sensor or is unsupported
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetTemperatureThreshold(nvmlDevice_t device,
                                                           nvmlTemperatureThresholds_t thresholdType,
                                                           unsigned int *temp);

    /**
     * Retrieves the intended operating speed of the device's fan.
     *
     * Note: The reported speed is the intended fan speed.  If the fan is physically blocked and unable to spin, the
     * output will not match the actual fan speed.
     *
     * For all discrete products with dedicated fans.
     *
     * The fan speed is expressed as a percentage of the product's maximum noise tolerance fan speed.
     * This value may exceed 100% in certain cases.
     *
     * @param device                               The identifier of the target device
     * @param speed                                Reference in which to return the fan speed percentage
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a speed has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a speed is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not have a fan
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetFanSpeed(nvmlDevice_t device, unsigned int *speed);

    /**
     * Retrieves the current clock speeds for the device.
     *
     * For Fermi &tm; or newer fully supported devices.
     *
     * See \ref nvmlClockType_t for details on available clock information.
     *
     * @param device                               The identifier of the target device
     * @param type                                 Identify which clock domain to query
     * @param clock                                Reference in which to return the clock speed in MHz
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a clock has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a clock is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device cannot report the specified clock
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetClockInfo(nvmlDevice_t device, nvmlClockType_t type, unsigned int *clock);

    /**
     * Retrieves the amount of used, free, reserved and total memory available on the device, in bytes.
     * The reserved amount is supported on version 2 only.
     *
     * For all products.
     *
     * Enabling ECC reduces the amount of total available memory, due to the extra required parity bits.
     * Under WDDM most device memory is allocated and managed on startup by Windows.
     *
     * Under Linux and Windows TCC, the reported amount of used memory is equal to the sum of memory allocated
     * by all active channels on the device.
     *
     * See \ref nvmlMemory_v2_t for details on available memory info.
     *
     * @note In MIG mode, if device handle is provided, the API returns aggregate
     *       information, only if the caller has appropriate privileges. Per-instance
     *       information can be queried by using specific MIG device handles.
     *
     * @note nvmlDeviceGetMemoryInfo_v2 adds additional memory information.
     *
     * @param device                               The identifier of the target device
     * @param memory                               Reference in which to return the memory information
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a memory has been populated
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_NO_PERMISSION     if the user doesn't have permission to perform this operation
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a memory is NULL
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetMemoryInfo(nvmlDevice_t device, nvmlMemory_t *memory);
    nvmlReturn_t DECLDIR nvmlDeviceGetMemoryInfo_v2(nvmlDevice_t device, nvmlMemory_v2_t *memory);

    /**
     * Retrieves the intended operating speed of the device's specified fan.
     *
     * Note: The reported speed is the intended fan speed. If the fan is physically blocked and unable to spin, the
     * output will not match the actual fan speed.
     *
     * For all discrete products with dedicated fans.
     *
     * The fan speed is expressed as a percentage of the product's maximum noise tolerance fan speed.
     * This value may exceed 100% in certain cases.
     *
     * @param device                                The identifier of the target device
     * @param fan                                   The index of the target fan, zero indexed.
     * @param speed                                 Reference in which to return the fan speed percentage
     *
     * @return
     *        - \ref NVML_SUCCESS                   if \a speed has been set
     *        - \ref NVML_ERROR_UNINITIALIZED       if the library has not been successfully initialized
     *        - \ref NVML_ERROR_INVALID_ARGUMENT    if \a device is invalid, \a fan is not an acceptable index, or \a speed is NULL
     *        - \ref NVML_ERROR_NOT_SUPPORTED       if the device does not have a fan or is newer than Maxwell
     *        - \ref NVML_ERROR_GPU_IS_LOST         if the target GPU has fallen off the bus or is otherwise inaccessible
     *        - \ref NVML_ERROR_UNKNOWN             on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetFanSpeed_v2(nvmlDevice_t device, unsigned int fan, unsigned int *speed);

    /**
     * Retrieves the current utilization rates for the device's major subsystems.
     *
     * For Fermi &tm; or newer fully supported devices.
     *
     * See \ref nvmlUtilization_t for details on available utilization rates.
     *
     * \note During driver initialization when ECC is enabled one can see high GPU and Memory Utilization readings.
     *       This is caused by ECC Memory Scrubbing mechanism that is performed during driver initialization.
     *
     * @note On MIG-enabled GPUs, querying device utilization rates is not currently supported.
     *
     * @param device                               The identifier of the target device
     * @param utilization                          Reference in which to return the utilization information
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a utilization has been populated
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a utilization is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support this feature
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetUtilizationRates(nvmlDevice_t device, nvmlUtilization_t *utilization);

    /**
     * Retrieves the current compute mode for the device.
     *
     * For all products.
     *
     * See \ref nvmlComputeMode_t for details on allowed compute modes.
     *
     * @param device                               The identifier of the target device
     * @param mode                                 Reference in which to return the current compute mode
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a mode has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a mode is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support this feature
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     *
     * @see nvmlDeviceSetComputeMode()
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetComputeMode(nvmlDevice_t device, nvmlComputeMode_t *mode);

    /**
     * Retrieves the CUDA compute capability of the device.
     *
     * For all products.
     *
     * Returns the major and minor compute capability version numbers of the
     * device.  The major and minor versions are equivalent to the
     * CU_DEVICE_ATTRIBUTE_COMPUTE_CAPABILITY_MINOR and
     * CU_DEVICE_ATTRIBUTE_COMPUTE_CAPABILITY_MAJOR attributes that would be
     * returned by CUDA's cuDeviceGetAttribute().
     *
     * @param device                               The identifier of the target device
     * @param major                                Reference in which to return the major CUDA compute capability
     * @param minor                                Reference in which to return the minor CUDA compute capability
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a major and \a minor have been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a major or \a minor are NULL
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetCudaComputeCapability(nvmlDevice_t device, int *major, int *minor);

    /**
     * Retrieves the current and pending ECC modes for the device.
     *
     * For Fermi &tm; or newer fully supported devices.
     * Only applicable to devices with ECC.
     * Requires \a NVML_INFOROM_ECC version 1.0 or higher.
     *
     * Changing ECC modes requires a reboot. The "pending" ECC mode refers to the target mode following
     * the next reboot.
     *
     * See \ref nvmlEnableState_t for details on allowed modes.
     *
     * @param device                               The identifier of the target device
     * @param current                              Reference in which to return the current ECC mode
     * @param pending                              Reference in which to return the pending ECC mode
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a current and \a pending have been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or either \a current or \a pending is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support this feature
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     *
     * @see nvmlDeviceSetEccMode()
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetEccMode(nvmlDevice_t device,
                                              nvmlEnableState_t *current,
                                              nvmlEnableState_t *pending);

    /**
     * Retrieves the default ECC modes for the device.
     *
     * For Fermi &tm; or newer fully supported devices.
     * Only applicable to devices with ECC.
     * Requires \a NVML_INFOROM_ECC version 1.0 or higher.
     *
     * See \ref nvmlEnableState_t for details on allowed modes.
     *
     * @param device                               The identifier of the target device
     * @param defaultMode                          Reference in which to return the default ECC mode
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a current and \a pending have been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a default is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support this feature
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     *
     * @see nvmlDeviceSetEccMode()
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetDefaultEccMode(nvmlDevice_t device, nvmlEnableState_t *defaultMode);

    /**
     * Retrieves the device boardId from 0-N.
     * Devices with the same boardId indicate GPUs connected to the same PLX.  Use in conjunction with
     *  \ref nvmlDeviceGetMultiGpuBoard() to decide if they are on the same board as well.
     *  The boardId returned is a unique ID for the current configuration.  Uniqueness and ordering across
     *  reboots and system configurations is not guaranteed (i.e. if a Tesla K40c returns 0x100 and
     *  the two GPUs on a Tesla K10 in the same system returns 0x200 it is not guaranteed they will
     *  always return those values but they will always be different from each other).
     *
     *
     * For Fermi &tm; or newer fully supported devices.
     *
     * @param device                               The identifier of the target device
     * @param boardId                              Reference in which to return the device's board ID
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a boardId has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a boardId is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support this feature
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetBoardId(nvmlDevice_t device, unsigned int *boardId);

    /**
     * Retrieves whether the device is on a Multi-GPU Board
     * Devices that are on multi-GPU boards will set \a multiGpuBool to a non-zero value.
     *
     * For Fermi &tm; or newer fully supported devices.
     *
     * @param device                               The identifier of the target device
     * @param multiGpuBool                         Reference in which to return a zero or non-zero value
     *                                                 to indicate whether the device is on a multi GPU board
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a multiGpuBool has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a multiGpuBool is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support this feature
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetMultiGpuBoard(nvmlDevice_t device, unsigned int *multiGpuBool);

    /**
     * Retrieves the PCI attributes of this device.
     *
     * For all products.
     *
     * See \ref nvmlPciInfo_t for details on the available PCI info.
     *
     * @param device                               The identifier of the target device
     * @param pci                                  Reference in which to return the PCI info
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a pci has been populated
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a pci is NULL
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetPciInfo_v3(nvmlDevice_t device, nvmlPciInfo_t *pci);

    /**
     * Retrieves the maximum PCIe link generation possible with this device and system
     *
     * I.E. for a generation 2 PCIe device attached to a generation 1 PCIe bus the max link generation this function
     * will report is generation 1.
     *
     * For Fermi &tm; or newer fully supported devices.
     *
     * @param device                               The identifier of the target device
     * @param maxLinkGen                           Reference in which to return the max PCIe link generation
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a maxLinkGen has been populated
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a maxLinkGen is null
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if PCIe link information is not available
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetMaxPcieLinkGeneration(nvmlDevice_t device, unsigned int *maxLinkGen);

    /**
     * Retrieves the maximum PCIe link width possible with this device and system
     *
     * I.E. for a device with a 16x PCIe bus width attached to a 8x PCIe system bus this function will report
     * a max link width of 8.
     *
     * For Fermi &tm; or newer fully supported devices.
     *
     * @param device                               The identifier of the target device
     * @param maxLinkWidth                         Reference in which to return the max PCIe link generation
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a maxLinkWidth has been populated
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a maxLinkWidth is null
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if PCIe link information is not available
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetMaxPcieLinkWidth(nvmlDevice_t device, unsigned int *maxLinkWidth);

    /**
     * Retrieves the current PCIe link generation
     *
     * For Fermi &tm; or newer fully supported devices.
     *
     * @param device                               The identifier of the target device
     * @param currLinkGen                          Reference in which to return the current PCIe link generation
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a currLinkGen has been populated
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a currLinkGen is null
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if PCIe link information is not available
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetCurrPcieLinkGeneration(nvmlDevice_t device, unsigned int *currLinkGen);

    /**
     * Retrieves the current PCIe link width
     *
     * For Fermi &tm; or newer fully supported devices.
     *
     * @param device                               The identifier of the target device
     * @param currLinkWidth                        Reference in which to return the current PCIe link generation
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a currLinkWidth has been populated
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a currLinkWidth is null
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if PCIe link information is not available
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetCurrPcieLinkWidth(nvmlDevice_t device, unsigned int *currLinkWidth);

    /**
     * Retrieve PCIe utilization information.
     * This function is querying a byte counter over a 20ms interval and thus is the
     *   PCIe throughput over that interval.
     *
     * For Maxwell &tm; or newer fully supported devices.
     *
     * This method is not supported in virtual machines running virtual GPU (vGPU).
     *
     * @param device                               The identifier of the target device
     * @param counter                              The specific counter that should be queried \ref
     * nvmlPcieUtilCounter_t
     * @param value                                Reference in which to return throughput in KB/s
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a value has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device or \a counter is invalid, or \a value is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support this feature
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetPcieThroughput(nvmlDevice_t device,
                                                     nvmlPcieUtilCounter_t counter,
                                                     unsigned int *value);

    /**
     * @deprecated Use \ref nvmlDeviceGetCurrentClocksEventReasons instead
     */

    /**
     * Retrieve the PCIe replay counter.
     *
     * For Kepler &tm; or newer fully supported devices.
     *
     * @param device                               The identifier of the target device
     * @param value                                Reference in which to return the counter's value
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a value has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid, or \a value is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support this feature
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetPcieReplayCounter(nvmlDevice_t device, unsigned int *value);

    /**
     * Retrieves the NVML index of this device.
     *
     * For all products.
     *
     * Valid indices are derived from the \a accessibleDevices count returned by
     *   \ref nvmlDeviceGetCount_v2(). For example, if \a accessibleDevices is 2 the valid indices
     *   are 0 and 1, corresponding to GPU 0 and GPU 1.
     *
     * The order in which NVML enumerates devices has no guarantees of consistency between reboots. For that reason it
     *   is recommended that devices be looked up by their PCI ids or GPU UUID. See
     *   \ref nvmlDeviceGetHandleByPciBusId_v2() and \ref nvmlDeviceGetHandleByUUID().
     *
     * When used with MIG device handles this API returns indices that can be
     * passed to \ref nvmlDeviceGetMigDeviceHandleByIndex to retrieve an identical handle.
     * MIG device indices are unique within a device.
     *
     * Note: The NVML index may not correlate with other APIs, such as the CUDA device index.
     *
     * @param device                               The identifier of the target device
     * @param index                                Reference in which to return the NVML index of the device
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a index has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid, or \a index is NULL
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     *
     * @see nvmlDeviceGetHandleByIndex()
     * @see nvmlDeviceGetCount()
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetIndex(nvmlDevice_t device, unsigned int *index);

    /**
     * Retrieves the globally unique board serial number associated with this device's board.
     *
     * For all products with an inforom.
     *
     * The serial number is an alphanumeric string that will not exceed 30 characters (including the NULL terminator).
     * This number matches the serial number tag that is physically attached to the board.  See \ref
     * nvmlConstants::NVML_DEVICE_SERIAL_BUFFER_SIZE.
     *
     * @param device                               The identifier of the target device
     * @param serial                               Reference in which to return the board/module serial number
     * @param length                               The maximum allowed length of the string returned in \a serial
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a serial has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid, or \a serial is NULL
     *         - \ref NVML_ERROR_INSUFFICIENT_SIZE if \a length is too small
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support this feature
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetSerial(nvmlDevice_t device, char *serial, unsigned int length);

    /**
     * Acquire the handle for a particular device, based on its board serial number.
     *
     * For Fermi &tm; or newer fully supported devices.
     *
     * This number corresponds to the value printed directly on the board, and to the value returned by
     *   \ref nvmlDeviceGetSerial().
     *
     * @deprecated Since more than one GPU can exist on a single board this function is deprecated in favor
     *             of \ref nvmlDeviceGetHandleByUUID.
     *             For dual GPU boards this function will return NVML_ERROR_INVALID_ARGUMENT.
     *
     * Starting from NVML 5, this API causes NVML to initialize the target GPU
     * NVML may initialize additional GPUs as it searches for the target GPU
     *
     * @param serial                               The board serial number of the target GPU
     * @param device                               Reference in which to return the device handle
     *
     * @return
     *         - \ref NVML_SUCCESS                  if \a device has been set
     *         - \ref NVML_ERROR_UNINITIALIZED      if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT   if \a serial is invalid, \a device is NULL or more than one
     *                                              device has the same serial (dual GPU boards)
     *         - \ref NVML_ERROR_NOT_FOUND          if \a serial does not match a valid device on the system
     *         - \ref NVML_ERROR_INSUFFICIENT_POWER if any attached devices have improperly attached external power
     * cables
     *         - \ref NVML_ERROR_IRQ_ISSUE          if NVIDIA kernel detected an interrupt issue with the attached GPUs
     *         - \ref NVML_ERROR_GPU_IS_LOST        if any GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN            on any unexpected error
     *
     * @see nvmlDeviceGetSerial
     * @see nvmlDeviceGetHandleByUUID
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetHandleBySerial(const char *serial, nvmlDevice_t *device);

    /**
     * Get VBIOS version of the device.
     *
     * For all products.
     *
     * The VBIOS version may change from time to time. It will not exceed 32 characters in length
     * (including the NULL terminator).  See \ref nvmlConstants::NVML_DEVICE_VBIOS_VERSION_BUFFER_SIZE.
     *
     * @param device                               The identifier of the target device
     * @param version                              Reference to which to return the VBIOS version
     * @param length                               The maximum allowed length of the string returned in \a version
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a version has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid, or \a version is NULL
     *         - \ref NVML_ERROR_INSUFFICIENT_SIZE if \a length is too small
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetVbiosVersion(nvmlDevice_t device, char *version, unsigned int length);

    /**
     * Retrieves the the device board part number which is programmed into the board's InfoROM
     *
     * For all products.
     *
     * @param device                                Identifier of the target device
     * @param partNumber                            Reference to the buffer to return
     * @param length                                Length of the buffer reference
     *
     * @return
     *         - \ref NVML_SUCCESS                  if \a partNumber has been set
     *         - \ref NVML_ERROR_UNINITIALIZED      if the library has not been successfully initialized
     *         - \ref NVML_ERROR_NOT_SUPPORTED      if the needed VBIOS fields have not been filled
     *         - \ref NVML_ERROR_INVALID_ARGUMENT   if \a device is invalid or \a serial is NULL
     *         - \ref NVML_ERROR_GPU_IS_LOST        if the target GPU has fallen off the bus or is otherwise
     * inaccessible
     *         - \ref NVML_ERROR_UNKNOWN            on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetBoardPartNumber(nvmlDevice_t device, char *partNumber, unsigned int length);

    /**
     * Set the ECC mode for the device.
     *
     * For Kepler &tm; or newer fully supported devices.
     * Only applicable to devices with ECC.
     * Requires \a NVML_INFOROM_ECC version 1.0 or higher.
     * Requires root/admin permissions.
     *
     * The ECC mode determines whether the GPU enables its ECC support.
     *
     * This operation takes effect after the next reboot.
     *
     * See \ref nvmlEnableState_t for details on available modes.
     *
     * @param device                               The identifier of the target device
     * @param ecc                                  The target ECC mode
     *
     * @return
     *         - \ref NVML_SUCCESS                 if the ECC mode was set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a ecc is invalid
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support this feature
     *         - \ref NVML_ERROR_NO_PERMISSION     if the user doesn't have permission to perform this operation
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     *
     * @see nvmlDeviceGetEccMode()
     */
    nvmlReturn_t DECLDIR nvmlDeviceSetEccMode(nvmlDevice_t device, nvmlEnableState_t ecc);

    /**
     * Retrieves the current and pending ECC modes for the device.
     *
     * For Fermi &tm; or newer fully supported devices.
     * Only applicable to devices with ECC.
     * Requires \a NVML_INFOROM_ECC version 1.0 or higher.
     *
     * Changing ECC modes requires a reboot. The "pending" ECC mode refers to the target mode following
     * the next reboot.
     *
     * See \ref nvmlEnableState_t for details on allowed modes.
     *
     * @param device                               The identifier of the target device
     * @param current                              Reference in which to return the current ECC mode
     * @param pending                              Reference in which to return the pending ECC mode
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a current and \a pending have been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or either \a current or \a pending is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support this feature
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     *
     * @see nvmlDeviceSetEccMode()
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetEccMode(nvmlDevice_t device,
                                              nvmlEnableState_t *current,
                                              nvmlEnableState_t *pending);

    /**
     * Retrieves the default ECC modes for the device.
     *
     * For Fermi &tm; or newer fully supported devices.
     * Only applicable to devices with ECC.
     * Requires \a NVML_INFOROM_ECC version 1.0 or higher.
     *
     * See \ref nvmlEnableState_t for details on allowed modes.
     *
     * @param device                               The identifier of the target device
     * @param defaultMode                          Reference in which to return the default ECC mode
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a current and \a pending have been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a default is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support this feature
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     *
     * @see nvmlDeviceSetEccMode()
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetDefaultEccMode(nvmlDevice_t device, nvmlEnableState_t *defaultMode);

    /**
     * Retrieves power usage for this GPU in milliwatts and its associated circuitry (e.g. memory)
     *
     * For Fermi &tm; or newer fully supported devices.
     *
     * On Fermi and Kepler GPUs the reading is accurate to within +/- 5% of current power draw.
     *
     * It is only available if power management mode is supported. See \ref nvmlDeviceGetPowerManagementMode.
     *
     * @param device                               The identifier of the target device
     * @param power                                Reference in which to return the power usage information
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a power has been populated
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a power is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support power readings
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetPowerUsage(nvmlDevice_t device, unsigned int *power);

    /**
     * Check if the GPU devices are on the same physical board.
     *
     * For all fully supported products.
     *
     * @param device1                               The first GPU device
     * @param device2                               The second GPU device
     * @param onSameBoard                           Reference in which to return the status.
     *                                              Non-zero indicates that the GPUs are on the same board.
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a onSameBoard has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a dev1 or \a dev2 are invalid or \a onSameBoard is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if this check is not supported by the device
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the either GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceOnSameBoard(nvmlDevice_t device1, nvmlDevice_t device2, int *onSameBoard);

    /**
     * Get information about processes with a compute context on a device
     *
     * For Fermi &tm; or newer fully supported devices.
     *
     * This function returns information only about compute running processes (e.g. CUDA application which have
     * active context). Any graphics applications (e.g. using OpenGL, DirectX) won't be listed by this function.
     *
     * To query the current number of running compute processes, call this function with *infoCount = 0. The
     * return code will be NVML_ERROR_INSUFFICIENT_SIZE, or NVML_SUCCESS if none are running. For this call
     * \a infos is allowed to be NULL.
     *
     * The usedGpuMemory field returned is all of the memory used by the application.
     *
     * Keep in mind that information returned by this call is dynamic and the number of elements might change in
     * time. Allocate more space for \a infos table in case new compute processes are spawned.
     *
     * @note In MIG mode, if device handle is provided, the API returns aggregate information, only if
     *       the caller has appropriate privileges. Per-instance information can be queried by using
     *       specific MIG device handles.
     *       Querying per-instance information using MIG device handles is not supported if the device is in vGPU Host virtualization mode.
     *
     * @param device                               The device handle or MIG device handle
     * @param infoCount                            Reference in which to provide the \a infos array size, and
     *                                             to return the number of returned elements
     * @param infos                                Reference in which to return the process information
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a infoCount and \a infos have been populated
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INSUFFICIENT_SIZE if \a infoCount indicates that the \a infos array is too small
     *                                             \a infoCount will contain minimal amount of space necessary for
     *                                             the call to complete
     *         - \ref NVML_ERROR_NO_PERMISSION     if the user doesn't have permission to perform this operation
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid, either of \a infoCount or \a infos is NULL
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if this query is not supported by \a device
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     *
     * @see \ref nvmlSystemGetProcessName
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetComputeRunningProcesses(nvmlDevice_t device, unsigned int *infoCount, nvmlProcessInfo_v1_t *infos);

    /** @} */ // @defgroup nvmlGPMStructs

    /***************************************************************************************************/
    /** @defgroup nvmlGpmFunctions GPM Functions
     *  @{
     */
    /***************************************************************************************************/

    /**
     * Calculate GPM metrics from two samples.
     *
     * For Hopper &tm; or newer fully supported devices.
     *
     * @param metricsGet             IN/OUT: populated \a nvmlGpmMetricsGet_t struct
     *
     * @return
     *         - \ref NVML_SUCCESS on success
     *         - Nonzero NVML_ERROR_? enum on error
     */
    nvmlReturn_t DECLDIR nvmlGpmMetricsGet(nvmlGpmMetricsGet_t *metricsGet);

    /**
     * Indicate whether the supplied device supports GPM
     *
     * @param device                NVML device to query for
     * @param gpmSupport            Structure to indicate GPM support \a nvmlGpmSupport_t. Indicates
     *                              GPM support per system for the supplied device
     *
     * @return
     *         - NVML_SUCCESS on success
     *         - Nonzero NVML_ERROR_? enum if there is an error in processing the query
     */
    nvmlReturn_t DECLDIR nvmlGpmQueryDeviceSupport(nvmlDevice_t device, nvmlGpmSupport_t *gpmSupport);

    /**
     * Free an allocated sample buffer that was allocated with \ref nvmlGpmSampleAlloc()
     *
     * For Hopper &tm; or newer fully supported devices.
     *
     * @param gpmSample              Sample to free
     *
     * @return
     *         - \ref NVML_SUCCESS                on success
     *         - \ref NVML_ERROR_INVALID_ARGUMENT if an invalid pointer is provided
     */
    nvmlReturn_t DECLDIR nvmlGpmSampleFree(nvmlGpmSample_t gpmSample);

    /**
     * Allocate a sample buffer to be used with NVML GPM . You will need to allocate
     * at least two of these buffers to use with the NVML GPM feature
     *
     * For Hopper &tm; or newer fully supported devices.
     *
     * @param gpmSample             Where  the allocated sample will be stored
     *
     * @return
     *         - \ref NVML_SUCCESS                on success
     *         - \ref NVML_ERROR_INVALID_ARGUMENT if an invalid pointer is provided
     *         - \ref NVML_ERROR_MEMORY           if system memory is insufficient
     */
    nvmlReturn_t DECLDIR nvmlGpmSampleAlloc(nvmlGpmSample_t *gpmSample);

    /**
     * Read a sample of GPM metrics into the provided \a gpmSample buffer. After
     * two samples are gathered, you can call nvmlGpmMetricGet on those samples to
     * retrive metrics
     *
     * For Hopper &tm; or newer fully supported devices.
     *
     * @param device                Device to get samples for
     * @param gpmSample             Buffer to read samples into
     *
     * @return
     *         - \ref NVML_SUCCESS on success
     *         - Nonzero NVML_ERROR_? enum on error
     */
    nvmlReturn_t DECLDIR nvmlGpmSampleGet(nvmlDevice_t device, nvmlGpmSample_t gpmSample);

    /**
     * Retrieves the power management limit associated with this device.
     *
     * For Fermi &tm; or newer fully supported devices.
     *
     * The power limit defines the upper boundary for the card's power draw. If
     * the card's total power draw reaches this limit the power management algorithm kicks in.
     *
     * This reading is only available if power management mode is supported.
     * See \ref nvmlDeviceGetPowerManagementMode.
     *
     * @param device                               The identifier of the target device
     * @param limit                                Reference in which to return the power management limit in milliwatts
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a limit has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a limit is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support this feature
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetPowerManagementLimit(nvmlDevice_t device, unsigned int *limit);

    /**
     * Retrieves information about possible values of power management limits on this device.
     *
     * For Kepler &tm; or newer fully supported devices.
     *
     * @param device                               The identifier of the target device
     * @param minLimit                             Reference in which to return the minimum power management limit in
     * milliwatts
     * @param maxLimit                             Reference in which to return the maximum power management limit in
     * milliwatts
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a minLimit and \a maxLimit have been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a minLimit or \a maxLimit is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support this feature
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     *
     * @see nvmlDeviceSetPowerManagementLimit
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetPowerManagementLimitConstraints(nvmlDevice_t device,
                                                                      unsigned int *minLimit,
                                                                      unsigned int *maxLimit);

    /**
     * Retrieves default power management limit on this device, in milliwatts.
     * Default power management limit is a power management limit that the device boots with.
     *
     * For Kepler &tm; or newer fully supported devices.
     *
     * @param device                               The identifier of the target device
     * @param defaultLimit                         Reference in which to return the default power management limit in
     * milliwatts
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a defaultLimit has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a defaultLimit is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support this feature
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetPowerManagementDefaultLimit(nvmlDevice_t device, unsigned int *defaultLimit);

    /**
     * Retrieves current clocks throttling reasons.
     *
     * For all fully supported products.
     *
     * \note More than one bit can be enabled at the same time. Multiple reasons can be affecting clocks at once.
     *
     * @param device                                The identifier of the target device
     * @param clocksThrottleReasons                 Reference in which to return bitmask of active clocks throttle
     *                                                  reasons
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a clocksThrottleReasons has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a clocksThrottleReasons is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device does not support this feature
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     *
     * @see nvmlClocksThrottleReasons
     * @see nvmlDeviceGetSupportedClocksThrottleReasons
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetCurrentClocksThrottleReasons(nvmlDevice_t device,
                                                                   unsigned long long *clocksThrottleReasons);

    /**
     * Retrieves bitmask of supported clocks throttle reasons that can be returned by
     * \ref nvmlDeviceGetCurrentClocksThrottleReasons
     *
     * For all fully supported products.
     *
     * This method is not supported in virtual machines running virtual GPU (vGPU).
     *
     * @param device                               The identifier of the target device
     * @param supportedClocksThrottleReasons       Reference in which to return bitmask of supported
     *                                              clocks throttle reasons
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a supportedClocksThrottleReasons has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device is invalid or \a supportedClocksThrottleReasons is NULL
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     *
     * @see nvmlClocksThrottleReasons
     * @see nvmlDeviceGetCurrentClocksThrottleReasons
     */
    nvmlReturn_t DECLDIR
    nvmlDeviceGetSupportedClocksThrottleReasons(nvmlDevice_t device,
                                                unsigned long long *supportedClocksThrottleReasons);

    /**
     * Retrieve the common ancestor for two devices
     * For all products.
     * Supported on Linux only.
     *
     * @param device1                              The identifier of the first device
     * @param device2                              The identifier of the second device
     * @param pathInfo                             A \ref nvmlGpuTopologyLevel_t that gives the path type
     *
     * @return
     *         - \ref NVML_SUCCESS                 if \a pathInfo has been set
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a device1, or \a device2 is invalid, or \a pathInfo is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the device or OS does not support this feature
     *         - \ref NVML_ERROR_UNKNOWN           an error has occurred in underlying topology discovery
     */

    /** @} */
    nvmlReturn_t DECLDIR nvmlDeviceGetTopologyCommonAncestor(nvmlDevice_t device1, nvmlDevice_t device2, nvmlGpuTopologyLevel_t *pathInfo);

    /***************************************************************************************************/
    /** @addtogroup nvmlEvents
     *  @{
     */
    /***************************************************************************************************/

    /**
     * Handle to an event set
     */
    typedef struct
    {
        struct nvmlEventSet_st *handle;
    } nvmlEventSet_t;

/** @defgroup nvmlEventType Event Types
 * @{
 * Event Types which user can be notified about.
 * See description of particular functions for details.
 *
 * See \ref nvmlDeviceRegisterEvents and \ref nvmlDeviceGetSupportedEventTypes to check which devices
 * support each event.
 *
 * Types can be combined with bitwise or operator '|' when passed to \ref nvmlDeviceRegisterEvents
 */
//! Event about single bit ECC errors
/**
 * \note A corrected texture memory error is not an ECC error, so it does not generate a single bit event
 */
#define nvmlEventTypeSingleBitEccError 0x0000000000000001LL

//! Event about double bit ECC errors
/**
 * \note An uncorrected texture memory error is not an ECC error, so it does not generate a double bit event
 */
#define nvmlEventTypeDoubleBitEccError 0x0000000000000002LL

//! Event about PState changes
/**
 *  \note On Fermi architecture PState changes are also an indicator that GPU is throttling down due to
 *  no work being executed on the GPU, power capping or thermal capping. In a typical situation,
 *  Fermi-based GPU should stay in P0 for the duration of the execution of the compute process.
 */
#define nvmlEventTypePState 0x0000000000000004LL

//! Event that Xid critical error occurred
#define nvmlEventTypeXidCriticalError 0x0000000000000008LL

//! Event about clock changes
/**
 * Kepler only
 */
#define nvmlEventTypeClock 0x0000000000000010LL

//! Event about AC/Battery power source changes
#define nvmlEventTypePowerSourceChange 0x0000000000000080LL

//! Event about MIG configuration changes
#define nvmlEventMigConfigChange 0x0000000000000100LL

//! Mask with no events
#define nvmlEventTypeNone 0x0000000000000000LL

//! Mask of all events
#define nvmlEventTypeAll                                                                                         \
    (nvmlEventTypeNone | nvmlEventTypeSingleBitEccError | nvmlEventTypeDoubleBitEccError | nvmlEventTypePState | \
     nvmlEventTypeClock | nvmlEventTypeXidCriticalError | nvmlEventTypePowerSourceChange | nvmlEventMigConfigChange)
    /** @} */

    /**
     * Information about occurred event
     */
    typedef struct nvmlEventData_st
    {
        nvmlDevice_t device;          //!< Specific device where the event occurred
        unsigned long long eventType; //!< Information about what specific event occurred
        unsigned long long
            eventData; //!< Stores XID error for the device in the event of nvmlEventTypeXidCriticalError,
                       //   eventData is 0 for any other event. eventData is set as 999 for unknown xid error.
        unsigned int
            gpuInstanceId; //!< If MIG is enabled and nvmlEventTypeXidCriticalError event is attributable to a GPU
                           //   instance, stores a valid GPU instance ID. gpuInstanceId is set to 0xFFFFFFFF
                           //   otherwise.
        unsigned int
            computeInstanceId; //!< If MIG is enabled and nvmlEventTypeXidCriticalError event is attributable to a
                               //   compute instance, stores a valid compute instance ID. computeInstanceId is set to
                               //   0xFFFFFFFF otherwise.
    } nvmlEventData_t;

    /** @} */

    /***************************************************************************************************/
    /** @defgroup nvmlEvents Event Handling Methods
     * This chapter describes methods that NVML can perform against each device to register and wait for
     * some event to occur.
     *  @{
     */
    /***************************************************************************************************/

    /**
     * Create an empty set of events.
     * Event set should be freed by \ref nvmlEventSetFree
     *
     * For Fermi &tm; or newer fully supported devices.
     * @param set                                  Reference in which to return the event handle
     *
     * @return
     *         - \ref NVML_SUCCESS                 if the event has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a set is NULL
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     *
     * @see nvmlEventSetFree
     */
    nvmlReturn_t DECLDIR nvmlEventSetCreate(nvmlEventSet_t *set);

    /**
     * Starts recording of events on a specified devices and add the events to specified \ref nvmlEventSet_t
     *
     * For Fermi &tm; or newer fully supported devices.
     * Ecc events are available only on ECC enabled devices (see \ref nvmlDeviceGetTotalEccErrors)
     * Power capping events are available only on Power Management enabled devices (see \ref
     * nvmlDeviceGetPowerManagementMode)
     *
     * For Linux only.
     *
     * \b IMPORTANT: Operations on \a set are not thread safe
     *
     * This call starts recording of events on specific device.
     * All events that occurred before this call are not recorded.
     * Checking if some event occurred can be done with \ref nvmlEventSetWait_v2
     *
     * If function reports NVML_ERROR_UNKNOWN, event set is in undefined state and should be freed.
     * If function reports NVML_ERROR_NOT_SUPPORTED, event set can still be used. None of the requested eventTypes
     *     are registered in that case.
     *
     * @param device                               The identifier of the target device
     * @param eventTypes                           Bitmask of \ref nvmlEventType to record
     * @param set                                  Set to which add new event types
     *
     * @return
     *         - \ref NVML_SUCCESS                 if the event has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a eventTypes is invalid or \a set is NULL
     *         - \ref NVML_ERROR_NOT_SUPPORTED     if the platform does not support this feature or some of requested
     * event types
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     *
     * @see nvmlEventType
     * @see nvmlDeviceGetSupportedEventTypes
     * @see nvmlEventSetWait
     * @see nvmlEventSetFree
     */
    nvmlReturn_t DECLDIR nvmlDeviceRegisterEvents(nvmlDevice_t device,
                                                  unsigned long long eventTypes,
                                                  nvmlEventSet_t set);

    /**
     * Returns information about events supported on device
     *
     * For Fermi &tm; or newer fully supported devices.
     *
     * Events are not supported on Windows. So this function returns an empty mask in \a eventTypes on Windows.
     *
     * @param device                               The identifier of the target device
     * @param eventTypes                           Reference in which to return bitmask of supported events
     *
     * @return
     *         - \ref NVML_SUCCESS                 if the eventTypes has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a eventType is NULL
     *         - \ref NVML_ERROR_GPU_IS_LOST       if the target GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     *
     * @see nvmlEventType
     * @see nvmlDeviceRegisterEvents
     */
    nvmlReturn_t DECLDIR nvmlDeviceGetSupportedEventTypes(nvmlDevice_t device, unsigned long long *eventTypes);

    /**
     * Waits on events and delivers events
     *
     * For Fermi &tm; or newer fully supported devices.
     *
     * If some events are ready to be delivered at the time of the call, function returns immediately.
     * If there are no events ready to be delivered, function sleeps till event arrives
     * but not longer than specified timeout. This function in certain conditions can return before
     * specified timeout passes (e.g. when interrupt arrives)
     *
     * On Windows, in case of xid error, the function returns the most recent xid error type seen by the system.
     * If there are multiple xid errors generated before nvmlEventSetWait is invoked then the last seen xid error
     * type is returned for all xid error events.
     *
     * On Linux, every xid error event would return the associated event data and other information if applicable.
     *
     * In MIG mode, if device handle is provided, the API reports all the events for the available instances,
     * only if the caller has appropriate privileges. In absence of required privileges, only the events which
     * affect all the instances (i.e. whole device) are reported.
     *
     * This API does not currently support per-instance event reporting using MIG device handles.
     *
     * @param set                                  Reference to set of events to wait on
     * @param data                                 Reference in which to return event data
     * @param timeoutms                            Maximum amount of wait time in milliseconds for registered event
     *
     * @return
     *         - \ref NVML_SUCCESS                 if the data has been set
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_INVALID_ARGUMENT  if \a data is NULL
     *         - \ref NVML_ERROR_TIMEOUT           if no event arrived in specified timeout or interrupt arrived
     *         - \ref NVML_ERROR_GPU_IS_LOST       if a GPU has fallen off the bus or is otherwise inaccessible
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     *
     * @see nvmlEventType
     * @see nvmlDeviceRegisterEvents
     */
    nvmlReturn_t DECLDIR nvmlEventSetWait_v2(nvmlEventSet_t set, nvmlEventData_t *data, unsigned int timeoutms);

    /**
     * Releases events in the set
     *
     * For Fermi &tm; or newer fully supported devices.
     *
     * @param set                                  Reference to events to be released
     *
     * @return
     *         - \ref NVML_SUCCESS                 if the event has been successfully released
     *         - \ref NVML_ERROR_UNINITIALIZED     if the library has not been successfully initialized
     *         - \ref NVML_ERROR_UNKNOWN           on any unexpected error
     *
     * @see nvmlDeviceRegisterEvents
     */
    nvmlReturn_t DECLDIR nvmlEventSetFree(nvmlEventSet_t set);

    /** @} */

    nvmlReturn_t DECLDIR ixmlDeviceGetBoardPosition(nvmlDevice_t device, unsigned int *position);

    nvmlReturn_t DECLDIR ixmlDeviceGetGPUVoltage(nvmlDevice_t device, unsigned int *integer, unsigned int *decimal);

    nvmlReturn_t DECLDIR ixmlDeviceGetEccErros(nvmlDevice_t device, unsigned int *single_error, unsigned int *double_error);

    nvmlReturn_t DECLDIR ixmlDeviceGetHealth(nvmlDevice_t device, unsigned long long *health);

#ifdef __cplusplus
}
#endif

#endif
