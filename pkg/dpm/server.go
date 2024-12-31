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
	"net"
	"os"
	"path"
	"time"

	"gitee.com/deep-spark/ix-device-plugin/pkg/config"
	"gitee.com/deep-spark/ix-device-plugin/pkg/gpuallocator"
	"github.com/jochenvg/go-udev"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"k8s.io/klog/v2"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const iluvatarDevicePluginSocket string = "iluvatar-gpu.sock"

// server is a grpc implementation between kubelet and iluvatar device plugin.
type server struct {
	// iluvatar device plugin implementation
	iluvatarDevicePlugin

	// socket for iluvatar device plugin
	socket string

	// socket for kubelet
	kubeletSocket string

	// iluvatar device plugin grpc server
	grpcServer *grpc.Server
}

func newServer(cfg *config.Config) *server {
	ret := &server{
		socket:        pluginapi.DevicePluginPath + iluvatarDevicePluginSocket,
		kubeletSocket: pluginapi.KubeletSocket,
		grpcServer:    nil,
		iluvatarDevicePlugin: iluvatarDevicePlugin{
			iluvatarDevice: iluvatarDevice{
				devSet:        gpuallocator.BuildDeviceSet(cfg),
				stopCheckHeal: make(chan struct{}),
				deviceCh:      make(chan *gpuallocator.Device),
			},
			name:     resourceName,
			stopList: make(chan struct{}),
		},
	}

	ret.devSet.ShowLayout()

	return ret
}

func (s *server) start() error {
	s.grpcServer = grpc.NewServer([]grpc.ServerOption{}...)

	err := s.createServer()
	if err != nil {
		klog.Errorf("Failed to create gprc server for '%s': %s", s.name, err)
		s.cleanup()
		return err
	}
	klog.Infof("Create grpc server '%s' on '%s'", s.name, s.socket)

	err = s.register()
	if err != nil {
		klog.Errorf("Failed to register device plugin: '%s'", s.name)
		s.stop()
		return err
	}
	klog.Infof("Register device plugin for '%s' with Kubelet", s.name)

	go s.checkHealth()

	return nil
}

func (s *server) stop() error {
	if s == nil || s.grpcServer == nil {
		return nil
	}
	klog.Infof("Stopping serve '%s' on %s", s.name, s.socket)
	s.grpcServer.Stop()
	if err := os.Remove(s.socket); err != nil && !os.IsNotExist(err) {
		return err
	}

	s.stopList <- struct{}{}
	s.stopCheckHeal <- struct{}{}

	s.cleanup()
	return nil
}

func (s *server) cleanup() {
	s.grpcServer = nil
}

func (s *server) createServer() error {
	// Remove the socket if exist.
	os.Remove(s.socket)

	// Create and Listen announces on the socket.
	sock, err := net.Listen("unix", s.socket)
	if err != nil {
		return err
	}

	pluginapi.RegisterDevicePluginServer(s.grpcServer, s)

	go func() {
		lastCrashTime := time.Now()
		restartCount := 0
		for {
			klog.Infof("Starting GRPC server for '%s'", s.name)
			err := s.grpcServer.Serve(sock)
			if err == nil {
				break
			}

			klog.Infof("GRPC server for '%s' crashed with error: %v", s.name, err)

			// restart if it has not been too often
			// i.e. if server has crashed more than 5 times and it didn't last more than one hour each time
			if restartCount > 5 {
				// quit
				klog.Fatalf("GRPC server for '%s' has repeatedly crashed recently. Quitting", s.name)
			}
			timeSinceLastCrash := time.Since(lastCrashTime).Seconds()
			lastCrashTime = time.Now()
			if timeSinceLastCrash > 3600 {
				// it has been one hour since the last crash.. reset the count
				// to reflect on the frequency
				restartCount = 1
			} else {
				restartCount++
			}
		}
	}()

	// Wait for server to start by launching a blocking connexion
	conn, err := s.dial(s.socket, 5*time.Second)
	if err != nil {
		return err
	}
	conn.Close()

	return nil
}

func (s *server) register() error {
	conn, err := s.dial(pluginapi.KubeletSocket, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pluginapi.NewRegistrationClient(conn)
	reqt := &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     path.Base(s.socket),
		ResourceName: s.name,
		Options: &pluginapi.DevicePluginOptions{
			GetPreferredAllocationAvailable: true,
		},
	}

	_, err = client.Register(context.Background(), reqt)
	if err != nil {
		return err
	}
	return nil
}

// dial establishes the gRPC communication with the registered device plugin.
func (s *server) dial(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
	c, err := grpc.Dial(unixSocketPath, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(timeout),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *server) updateUdev(dev *udev.Device) {
	s.devSet.UpdateUdev(dev)
}
