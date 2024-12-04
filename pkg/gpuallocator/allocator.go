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

package gpuallocator

type Allocator struct {
	GPUs []*Device

	policy    Policy
	remaining DeviceSet
	allocated DeviceSet
}

type PolicyArgs interface {
}

// Policy defines an interface for pluggable allocation policies to be added
// to an Allocator.
type Policy interface {
	// Allocate is meant to do the heavy-lifting of implementing the actual
	// allocation strategy of the policy. It takes a slice of devices to
	// allocate GPUs from, and an amount 'size' to allocate from that slice. It
	// then returns a subset of devices of length 'size'. If the policy is
	// unable to allocate 'size' GPUs from the slice of input devices, it
	// returns an empty slice.
	Allocate(arg PolicyArgs) []string
}
