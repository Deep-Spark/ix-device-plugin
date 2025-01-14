package kube

import (
	"context"
	"fmt"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"k8s.io/kubelet/pkg/apis/podresources/v1alpha1"
	"k8s.io/kubernetes/pkg/kubelet/apis/podresources"
)

const (
	socketPath                 = "/var/lib/kubelet/pod-resources/kubelet.sock"
	defaultPodResourcesMaxSize = 1024 * 1024 * 16
	callTimeout                = 2 * time.Second
)

// start starts the gRPC server, registers the pod resource with the Kubelet
func (pr *PodResource) start() error {
	pr.stop()
	var err error
	if pr.client, pr.conn, err = podresources.GetV1alpha1Client("unix://"+socketPath, callTimeout,
		defaultPodResourcesMaxSize); err != nil {
		klog.Errorf("get pod resource client failed, %v", err)
		return err
	}
	// klog.Infof("pod resource client init success.")
	return nil
}

// stop the connection
func (pr *PodResource) stop() {
	if pr == nil {
		klog.Errorf("invalid interface receiver")
		return
	}
	if pr.conn != nil {
		if err := pr.conn.Close(); err != nil {
			klog.Errorf("stop connect failed, err: %v", err)
		}
		pr.conn = nil
		pr.client = nil
	}
}

func (pr *PodResource) GetPodResource() (map[string]PodDevice, error) {
	if err := pr.start(); err != nil {
		return nil, err
	}
	defer pr.stop()
	return pr.assemblePodResource()
}

func (pr *PodResource) assemblePodResource() (map[string]PodDevice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), callTimeout)
	defer cancel()

	resp, err := pr.client.List(ctx, &v1alpha1.ListPodResourcesRequest{})
	if err != nil {
		return nil, fmt.Errorf("list pod resource failed, err: %v", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("invalid list response")
	}

	device := make(map[string]PodDevice, 1)
	for _, pod := range resp.PodResources {
		if pod == nil {
			klog.Warningf("invalid pod")
			continue
		}
		resourceName, podDevice, err := pr.getDeviceFromPod(pod)
		if err != nil || resourceName == "" || len(podDevice) == 0 {
			if err != nil {
				klog.Warningf("get device from pod maybe err: %v", err)
			}

			continue
		}
		device[pod.Namespace+"_"+pod.Name] = PodDevice{
			ResourceName: resourceName,
			DeviceIds:    podDevice,
		}
	}
	return device, nil
}

func (pr *PodResource) getDeviceFromPod(podResources *v1alpha1.PodResources) (string, []string, error) {
	var podDevice = make([]string, 0)
	var resourceName = ""
	for _, containerResource := range podResources.Containers {
		containerResourceName, containerDevices, err := pr.getContainerResource(containerResource)
		if err != nil {
			return "", nil, err
		}
		if containerResourceName == "" {
			continue
		}
		if resourceName == "" {
			resourceName = containerResourceName
		}
		podDevice = append(podDevice, containerDevices...)
	}
	return resourceName, podDevice, nil
}

func (pr *PodResource) getContainerResource(containerResource *v1alpha1.ContainerResources) (string, []string, error) {
	if containerResource == nil {
		return "", nil, fmt.Errorf("invalid container resource")
	}
	var deviceIds = make([]string, 0)
	ixResourceName := ResourceNamePrefix + GPU
	for _, containerDevice := range containerResource.Devices {
		if containerDevice == nil {
			klog.Warningf("invalid container device")
			continue
		}
		if containerDevice.ResourceName != ixResourceName {
			continue
		}

		deviceIds = append(deviceIds, containerDevice.DeviceIds...)
	}
	return ixResourceName, deviceIds, nil
}

func NewPodResource() *PodResource {
	return &PodResource{}
}

func GetKltAndRealAllocateDev(podList []v1.Pod) ([]*PodDeviceInfo, error) {
	prClient := NewPodResource()
	podDevice, err := prClient.GetPodResource()
	if err != nil {
		return nil, fmt.Errorf("get pod resource failed, %v", err)
	}

	var podDeviceInfo = make([]*PodDeviceInfo, 0)
	for _, pod := range podList {
		podKey := pod.Namespace + "_" + pod.Name
		podResource, exist := podDevice[podKey]
		if !exist {
			continue
		}

		realDeviceList := podResource.DeviceIds

		devStr, exist := pod.Annotations[ResourceNamePrefix+PodDevVolcano]
		if exist {
			realDeviceList = strings.Split(devStr, ",")
		}

		podDeviceInfo = append(podDeviceInfo, &PodDeviceInfo{Pod: pod, KltDevice: podResource.DeviceIds,
			RealDevice: realDeviceList})
	}
	return podDeviceInfo, nil
}
