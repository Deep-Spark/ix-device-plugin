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

import (
	"sort"
)

type replicaPolicy struct{}

type ReplicaPolicyArgs struct {
	PolicyArgs
	Device              ReplicaDeviceMap
	Available, Required []string
	Size                int
}

func NewReplicaPolicy() Policy {
	return &replicaPolicy{}
}

func (p *replicaPolicy) Allocate(arg PolicyArgs) []string {
	replicaArg, ok := (arg).(ReplicaPolicyArgs)
	if !ok {
		return []string{}
	}
	available := replicaArg.Available
	required := replicaArg.Required
	size := replicaArg.Size
	ReplicaMap := replicaArg.Device
	var ret []string

	replicaSet := make(map[string]*struct {
		Avail, Total int
		ReplicaStr   []string
	})
	candidatesMap := ReplicaMap.Subset(available).Difference(ReplicaMap.Subset(required))
	candidates := candidatesMap.GetIDs()
	candidatesPrefix, _ := PrefixUUID(candidates)

	for uuid, replicaDev := range candidatesMap {
		_, ok := replicaSet[replicaDev.Parent.UUID]
		if !ok {
			replicaSet[replicaDev.Parent.UUID] = new(struct {
				Avail, Total int
				ReplicaStr   []string
			})
		}
		tmp := replicaSet[replicaDev.Parent.UUID]
		tmp.Total = replicaDev.Parent.Replicas
		tmp.ReplicaStr = append(tmp.ReplicaStr, uuid)
		tmp.Avail++
	}
	needed := size - len(required)

	for i := 0; i < needed; i++ {
		sort.SliceStable(candidatesPrefix, func(i, j int) bool {
			idiff := replicaSet[candidatesPrefix[i]].Total - replicaSet[candidatesPrefix[i]].Avail
			jdiff := replicaSet[candidatesPrefix[j]].Total - replicaSet[candidatesPrefix[j]].Avail
			return idiff < jdiff
		})
		fetchSet := replicaSet[candidatesPrefix[0]]
		fetchSet.Avail--
		fetchId := fetchSet.ReplicaStr[0]
		fetchSet.ReplicaStr = fetchSet.ReplicaStr[1:]
		ret = append(ret, fetchId)
	}

	ret = append(ret, required...)
	return ret
}
