/*
Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
Copyright (c) NVIDIA CORPORATION. All Rights Reserved.

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

package gpuallocator

import (
	goixml "gitee.com/deep-spark/go-ixml/pkg/ixml"
	"gitee.com/deep-spark/ix-device-plugin/pkg/ixml"
)

type P2PLinkType uint32

type P2PLink struct {
	GPU  *Device
	Type P2PLinkType
}

const (
	P2PLinkUnknown P2PLinkType = iota
	P2PLinkCrossCPU
	P2PLinkSameCPU
	P2PLinkHostBridge
	P2PLinkMultiSwitch
	P2PLinkSingleSwitch
	P2PLinkSameBoard
	SingleNVLINKLink
	TwoNVLINKLinks
	ThreeNVLINKLinks
	FourNVLINKLinks
	FiveNVLINKLinks
	SixNVLINKLinks
	SevenNVLINKLinks
	EightNVLINKLinks
	NineNVLINKLinks
	TenNVLINKLinks
	ElevenNVLINKLinks
	TwelveNVLINKLinks
	ThirteenNVLINKLinks
	FourteenNVLINKLinks
	FifteenNVLINKLinks
	SixteenNVLINKLinks
	SeventeenNVLINKLinks
	EighteenNVLINKLinks
)

func P2PLinkTypeToString(linkType P2PLinkType) string {
	switch linkType {
	case P2PLinkUnknown:
		return "P2PLinkUnknown"
	case P2PLinkCrossCPU:
		return "P2PLinkCrossCPU"
	case P2PLinkSameCPU:
		return "P2PLinkSameCPU"
	case P2PLinkHostBridge:
		return "P2PLinkHostBridge"
	case P2PLinkMultiSwitch:
		return "P2PLinkMultiSwitch"
	case P2PLinkSingleSwitch:
		return "P2PLinkSingleSwitch"
	case P2PLinkSameBoard:
		return "P2PLinkSameBoard"
	case SingleNVLINKLink:
		return "SingleNVLINKLink"
	case TwoNVLINKLinks:
		return "TwoNVLINKLinks"
	case ThreeNVLINKLinks:
		return "ThreeNVLINKLinks"
	case FourNVLINKLinks:
		return "FourNVLINKLinks"
	case FiveNVLINKLinks:
		return "FiveNVLINKLinks"
	case SixNVLINKLinks:
		return "SixNVLINKLinks"
	case SevenNVLINKLinks:
		return "SevenNVLINKLinks"
	case EightNVLINKLinks:
		return "EightNVLINKLinks"
	case NineNVLINKLinks:
		return "NineNVLINKLinks"
	case TenNVLINKLinks:
		return "TenNVLINKLinks"
	case ElevenNVLINKLinks:
		return "ElevenNVLINKLinks"
	case TwelveNVLINKLinks:
		return "TwelveNVLINKLinks"
	case ThirteenNVLINKLinks:
		return "ThirteenNVLINKLinks"
	case FourteenNVLINKLinks:
		return "FourteenNVLINKLinks"
	case FifteenNVLINKLinks:
		return "FifteenNVLINKLinks"
	case SixteenNVLINKLinks:
		return "SixteenNVLINKLinks"
	case SeventeenNVLINKLinks:
		return "SeventeenNVLINKLinks"
	case EighteenNVLINKLinks:
		return "EighteenNVLINKLinks"
	default:
		return "Unknown P2PLinkType"
	}
}

func GetP2PLink(dev1 ixml.Device, dev2 ixml.Device) P2PLinkType {
	level, err := dev1.DeviceGetTopology(&dev2)
	if err != nil {
		return P2PLinkUnknown
	}
	switch level {
	case goixml.TOPOLOGY_INTERNAL:
		return P2PLinkSameBoard
	case goixml.TOPOLOGY_SINGLE:
		return P2PLinkSingleSwitch
	case goixml.TOPOLOGY_MULTIPLE:
		return P2PLinkMultiSwitch
	case goixml.TOPOLOGY_HOSTBRIDGE:
		return P2PLinkHostBridge
	case goixml.TOPOLOGY_NODE: // NVML_TOPOLOGY_CPU was renamed NVML_TOPOLOGY_NODE
		return P2PLinkSameCPU
	case goixml.TOPOLOGY_SYSTEM:
		return P2PLinkCrossCPU
	}

	return P2PLinkUnknown
}
