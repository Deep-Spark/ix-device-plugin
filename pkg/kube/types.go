package kube

import (
	"gitee.com/deep-spark/ix-device-plugin/pkg/gpuallocator"
	"google.golang.org/grpc"
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubelet/pkg/apis/podresources/v1alpha1"
)

// EventType string type of event
type EventType string

const (
	// EventTypeAdd is used when a new resource is created
	EventTypeAdd EventType = "add"
	// EventTypeUpdate is used when an existing resource is modified
	EventTypeUpdate EventType = "update"
	// EventTypeDelete is used when an existing resource is deleted
	EventTypeDelete EventType = "delete"
)

type metaData struct {
	Annotation map[string]string `json:"annotations"`
}

type podMetaData map[string]metaData

type PodResource struct {
	conn   *grpc.ClientConn
	client v1alpha1.PodResourcesListerClient
}

type PodDeviceInfo struct {
	Pod        v1.Pod
	KltDevice  []string
	RealDevice []string
}

type PodDevice struct {
	ResourceName string
	DeviceIds    []string
}

type P2PLink struct {
	TypeName  string
	TypeIndex gpuallocator.P2PLinkType
}

type DeviceInfo struct {
	Name  string
	UUID  string
	Links map[string][]P2PLink
}

type NodeDeviceInfo struct {
	DeviceInfo map[string]DeviceInfo
	UpdateTime int64
}

type NodeDeviceList struct {
	DeviceList []string
	UpdateTime int64
}

type NodeDeviceInfoCache struct {
	DeviceInfo NodeDeviceInfo
}
