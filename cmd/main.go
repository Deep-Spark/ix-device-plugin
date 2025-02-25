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

package main

import (
	"os"

	"gitee.com/deep-spark/ix-device-plugin/pkg/config"
	dpm "gitee.com/deep-spark/ix-device-plugin/pkg/dpm"
	"github.com/urfave/cli/v2"
	"k8s.io/klog/v2"
)

func main() {
	manager := dpm.NewManager()
	c := cli.NewApp()
	c.Action = func(ctx *cli.Context) error {
		return manager.Run(ctx, c.Flags)
	}
	c.Version = config.VERSION
	c.Name = "Iluvatar Device Plugin"
	c.Usage = "Iluvatar device plugin for Kubernetes"
	c.Flags = []cli.Flag{&cli.BoolFlag{
		Name:    "splitboard",
		Usage:   "chip is not exposed and managed by device plugin, versa board is managed by device plugin:\n\t\t[false, true]",
		EnvVars: []string{"SPLIT_BOARD"},
	}}

	defer klog.Flush()
	err := c.Run(os.Args)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
}
