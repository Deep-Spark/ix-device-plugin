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

package config

import (
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
	"sigs.k8s.io/yaml"
)

type Flags struct {
	SplitBoard bool `json:"splitboard"                yaml:"board"`
}

type ReplicatedResources struct {
	Replicas int `json:"replicas"         yaml:"replicas"`
}

// Sharing encapsulates the set of sharing strategies that are supported.
type Sharing struct {
	// TimeSlicing defines the set of replicas to be made for timeSlicing available resources.
	TimeSlicing ReplicatedResources `json:"timeSlicing,omitempty" yaml:"timeSlicing,omitempty"`
	// MPS defines the set of replicas to be shared using MPS
	MPS *ReplicatedResources `json:"mps,omitempty"         yaml:"mps,omitempty"`
}

// Config is a versioned struct used to hold configuration information.
type Config struct {
	Version      string  `json:"version"             yaml:"version"`
	ResourceName string  `json:"resourceName"         yaml:"resourceName"`
	Flags        Flags   `json:"flags,omitempty"     yaml:"flags,omitempty"`
	Sharing      Sharing `json:"sharing,omitempty"   yaml:"sharing,omitempty"`
}

func parseConfigFrom(reader io.Reader) (*Config, error) {
	var err error
	var configYaml []byte

	configYaml, err = io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read error: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(configYaml, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	if cfg.Version == "" {
		cfg.Version = VERSION
	}

	if cfg.Version != VERSION {
		return nil, fmt.Errorf("unknown version: %v", cfg.Version)
	}

	return &cfg, nil
}

func (f *Flags) UpdateFromCLIFlags(c *cli.Context, flags []cli.Flag) {
	for _, flag := range flags {
		for _, n := range flag.Names() {
			if !c.IsSet(n) {
				continue
			}
			// Common flags
			switch n {
			case "splitboard":
				f.SplitBoard = c.Bool(n)
			default:
				panic(fmt.Errorf("unsupported flag type for %v", n))
			}
		}
	}
}

func LoadConfig(c *cli.Context, flags []cli.Flag) (*Config, error) {
	reader, err := os.Open(ConfigDirectory)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}

	defer reader.Close()

	config, err := parseConfigFrom(reader)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	config.Flags.UpdateFromCLIFlags(c, flags)

	return config, nil
}
