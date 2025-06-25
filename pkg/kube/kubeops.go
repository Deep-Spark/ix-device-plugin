package kube

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
)

func (ki *KubeClient) WriteDevicelisttoCM(devices []string) error {

	var nodeDeviceListData = NodeDeviceList{
		DeviceList: devices,
		UpdateTime: time.Now().Unix(),
	}

	var devicelistdata []byte
	var err error

	if devicelistdata, err = json.Marshal(nodeDeviceListData); err != nil {
		return fmt.Errorf("marshal nodeDeviceListData failed: %v", err)
	}

	newConfigmapData := map[string]interface{}{
		"data": map[string]interface{}{
			"DeviceListCfg": string(devicelistdata),
		},
	}
	configmapUpdateData, err := json.Marshal(newConfigmapData)
	if err != nil {
		klog.Errorf("Failed to marshal configmap data: %v", err)
		return err
	}

	for i := 0; i < RetryUpdateCount; i++ {
		if _, err = ki.Client.CoreV1().ConfigMaps(DeviceInfoCMNameSpace).Patch(context.Background(),
			ki.DeviceInfoName, types.StrategicMergePatchType, configmapUpdateData, metav1.PatchOptions{}); err == nil {
			return nil
		}

		if errors.IsNotFound(err) {
			return err
		}

		klog.Warningf("patch configmap failed: %v, try again", err)
		time.Sleep(PatchWaitTime * time.Millisecond)
	}

	return err
}

func (ki *KubeClient) WriteDeviceInfoDataIntoCM(devices []string, deviceinfo map[string]DeviceInfo, verbose bool) error {

	updatetime := time.Now().Unix()
	var nodeDeviceInfoData = NodeDeviceInfo{
		DeviceInfo: deviceinfo,
		UpdateTime: updatetime,
	}

	var nodeDeviceListData = NodeDeviceList{
		DeviceList: devices,
		UpdateTime: updatetime,
	}

	var deviceinfodata, devicelistdata []byte
	var err error
	if deviceinfodata, err = json.Marshal(nodeDeviceInfoData); err != nil {
		return fmt.Errorf("marshal nodeDeviceInfoData failed: %v", err)
	}

	if devicelistdata, err = json.Marshal(nodeDeviceListData); err != nil {
		return fmt.Errorf("marshal nodeDeviceListData failed: %v", err)
	}

	deviceInfoCM := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ki.DeviceInfoName,
			Namespace: DeviceInfoCMNameSpace,
		},
	}

	deviceInfoCM.Data = map[string]string{
		DeviceInfoCMDataKey: string(deviceinfodata),
		DeviceListCMDataKey: string(devicelistdata),
	}
	if verbose {
		klog.Infof("write device info cache into cm: %s/%s.", deviceInfoCM.Namespace, deviceInfoCM.Name)
	}
	err = wait.PollImmediate(UpdateInternalTime*time.Second, UpdateTimeoutTime*time.Second, func() (bool, error) {
		if err := ki.createOrUpdateDeviceCM(deviceInfoCM); err != nil {
			return false, err
		}
		return true, nil
	})
	return err
}

func (ki *KubeClient) TryUpdatePodCacheAnnotation(pod *v1.Pod, annotation map[string]string) error {
	if err := ki.TryUpdatePodAnnotation(pod, annotation); err != nil {
		klog.Errorf("update pod annotation in api server failed, err: %v", err)
		return err
	}
	// update cache
	lock.Lock()
	defer lock.Unlock()
	for i, podInCache := range podCache {
		if podInCache.Namespace == pod.Namespace && podInCache.Name == pod.Name {
			for k, v := range annotation {
				podCache[i].Annotations[k] = v
			}
			klog.Infof("update annotation in pod cache success, name: %s, namespace: %s", pod.Name, pod.Namespace)
			return nil
		}
	}
	klog.Warningf("no pod found in cache when update annotation, name: %s, namespace: %s", pod.Name, pod.Namespace)
	return nil
}

func (ki *KubeClient) TryUpdatePodAnnotation(pod *v1.Pod, annotation map[string]string) error {

	newPodMetaData := podMetaData{"metadata": metaData{Annotation: annotation}}
	podUpdateMetaData, err := json.Marshal(newPodMetaData)
	if err != nil {
		klog.Errorf("Failed to marshal pod metadata: %v", err)
		return err
	}

	for i := 0; i < RetryUpdateCount; i++ {
		if _, err = ki.PatchPod(pod, podUpdateMetaData); err == nil {
			return nil
		}

		// There is no need to retry if the pod does not exist
		if errors.IsNotFound(err) {
			return err
		}

		klog.Warningf("patch pod annotation failed: %v, try again", err)
		time.Sleep(PatchWaitTime * time.Millisecond)
	}

	return fmt.Errorf("patch pod annotation failed, exceeded max number of retries")
}
