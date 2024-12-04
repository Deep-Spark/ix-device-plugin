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

package dpm

import (
	"fmt"
	"os"
	"syscall"

	"gitee.com/deep-spark/ix-device-plugin/pkg/config"
	"gitee.com/deep-spark/ix-device-plugin/pkg/ixml"
	"github.com/fsnotify/fsnotify"
	udev "github.com/jochenvg/go-udev"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/context"
	"k8s.io/klog/v2"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

// Manager contains the main machinery of iluvatar device plugin framwork.
type Manager struct {
	fsWatcher   *fsnotify.Watcher
	udevWatcher <-chan *udev.Device

	sigs chan os.Signal
}

// NewManager initialize Manger structure.
func NewManager() *Manager {
	return &Manager{
		fsWatcher: nil,
	}
}

// Run starts the Manager
func (m *Manager) Run(c *cli.Context, flags []cli.Flag) error {
	klog.Info("Loading configuration.")
	cfg, err := config.LoadConfig(c, flags)
	if err != nil {
		return fmt.Errorf("unable to load config: %v", err)
	}

	klog.Info("Loading IXML")
	err = ixml.Init()
	if err != nil {
		klog.Errorf("Failed to initialize IXML: %v", err)
		return fmt.Errorf("%v", err)
	}
	defer func() {
		klog.Info("Shutdown of IXML returned:", ixml.Shutdown())
	}()

	klog.Info("Starting FS watcher.")
	m.fsWatcher, err = newFSWatcher(pluginapi.DevicePluginPath)
	if err != nil {
		return fmt.Errorf("Failed to create FS watcher: %v", err)
	}
	defer m.fsWatcher.Close()

	klog.Info("Starting OS watcher.")
	m.sigs = newOSWatcher(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer close(m.sigs)

	//ix device udev watchr
	klog.Info("Starting IX device watcher.")
	ctx := context.Background()
	u := udev.Udev{}
	udevMonitor := u.NewMonitorFromNetlink("kernel")
	if udevMonitor == nil {
		return fmt.Errorf("Failed to create udev context: %v", err)
	}

	err = udevMonitor.FilterAddMatchSubsystem(config.UdevWatcherSubsystem)
	if err != nil {
		return fmt.Errorf("Failed to add udev subsystem: %v", err)
	}

	m.udevWatcher, err = udevMonitor.DeviceChan(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create udev watcher: %v", err)
	}

	server := newServer(cfg)
Restart:
	err = server.start()
	if err != nil {
		klog.Info("Failed to start plugin.")

		return fmt.Errorf("Failed to start plugin: %v", err)
	}

	/*
	 * 1. Stop plugin if kubelet exit.
	 * 2. Start plugin if kubelet is running.
	 * 3. Stop plugin if interrupted.
	 */
HandleEvents:
	for {
		select {
		case event := <-m.fsWatcher.Events:
			if event.Name == pluginapi.KubeletSocket {
				if event.Op&fsnotify.Create == fsnotify.Create {
					klog.Infof("Notify '%s' created, restarting plugin.", pluginapi.KubeletSocket)
					server.stop()
					goto Restart
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					klog.Infof("Detect '%s' removed, stopping plugin.", pluginapi.KubeletSocket)
					server.stop()
				}
			}
		case ixdev := <-m.udevWatcher:
			klog.Infof("udev:%v\n", ixdev.Sysname())
			server.updateUdev(ixdev)
		case s := <-m.sigs:
			switch s {
			case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
				klog.Infof("Received signal %v, shutting down.", s)
				server.stop()
				break HandleEvents
			}
		}
	}

	return nil
}
