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
package kube

import (
	"context"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

type KubeClient struct {
	Client         kubernetes.Interface
	NodeName       string
	DeviceInfoName string
	NeedRefresh    bool
	PodInformer    cache.SharedIndexInformer
	Queue          workqueue.RateLimitingInterface
}

func NewKubeClient() (*KubeClient, error) {
	clientCfg, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		klog.Errorf("Failed to get client config: %v", err.Error())
		return nil, err
	}

	client, err := kubernetes.NewForConfig(clientCfg)
	if err != nil {
		klog.Errorf("Failed to get client: %v", err.Error())
		return nil, err
	}

	nodeName, err := GetNodeNameFromEnv()
	if err != nil {
		return nil, err
	}

	return &KubeClient{
		Client:         client,
		NodeName:       nodeName,
		DeviceInfoName: DeviceInfoCMNamePrefix + nodeName,
		Queue:          workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		NeedRefresh:    false,
	}, nil
}

func GetNodeNameFromEnv() (string, error) {
	nodeName := os.Getenv("NODE_NAME")
	if err := checkNodeName(nodeName); err != nil {
		return "", fmt.Errorf("check node name failed: %v", err.Error())
	}
	return nodeName, nil
}

func checkNodeName(nodeName string) error {
	if len(nodeName) == 0 {
		return fmt.Errorf("the env variable whose key is NODE_NAME must be set")
	}
	if len(nodeName) > KubeEnvMaxLength {
		return fmt.Errorf("node name length %d is bigger than %d", len(nodeName), KubeEnvMaxLength)
	}
	pattern := NamePatterns["nodeName"]
	if match := pattern.MatchString(nodeName); !match {
		return fmt.Errorf("node name %s is illegal", nodeName)
	}
	return nil
}

func (ki *KubeClient) CreateConfigMap(cm *v1.ConfigMap) (*v1.ConfigMap, error) {
	newCM, err := ki.Client.CoreV1().ConfigMaps(cm.ObjectMeta.Namespace).
		Create(context.TODO(), cm, metav1.CreateOptions{})
	if err != nil {
		klog.Errorf("create configmap failed: %v", err)
	}

	return newCM, err
}

func (ki *KubeClient) GetConfigMap(cmName, cmNameSpace string) (*v1.ConfigMap, error) {
	newCM, err := ki.Client.CoreV1().ConfigMaps(cmNameSpace).Get(context.TODO(), cmName, metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		klog.Errorf("get configmap failed: %v", err)
	}

	return newCM, err
}

func (ki *KubeClient) UpdateConfigMap(cm *v1.ConfigMap) (*v1.ConfigMap, error) {
	if cm == nil {
		return nil, fmt.Errorf("param cm is nil")
	}
	newCM, err := ki.Client.CoreV1().ConfigMaps(cm.ObjectMeta.Namespace).
		Update(context.TODO(), cm, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorf("update configmap failed: %v", err)
	}

	return newCM, err
}

func (ki *KubeClient) createOrUpdateDeviceCM(cm *v1.ConfigMap) error {
	// use update first
	if _, err := ki.UpdateConfigMap(cm); errors.IsNotFound(err) {
		if _, err := ki.CreateConfigMap(cm); err != nil {
			return fmt.Errorf("unable to create configmap, %v", err)
		}
		return nil
	} else {
		return err
	}
}

func (ki *KubeClient) PatchPod(pod *v1.Pod, data []byte) (*v1.Pod, error) {
	return ki.Client.CoreV1().Pods(pod.Namespace).Patch(context.Background(),
		pod.Name, types.StrategicMergePatchType, data, metav1.PatchOptions{})
}
