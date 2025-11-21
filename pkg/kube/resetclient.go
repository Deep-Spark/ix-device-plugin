package kube

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"gitee.com/deep-spark/ix-device-plugin/pkg/ixml"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
)

type GpuResetInfo struct {
	NodeName string          `json:"nodename"             yaml:"nodename"`
	Reset    bool            `json:"reset"                yaml:"reset"`
	Occupy   map[string]bool `json:"occupy"                yaml:"occupy"`
}

type ResetClient struct {
	client      *KubeClient
	resetInfo   *GpuResetInfo
	cmLabel     string
	cmName      string
	cmNamespace string
	// only one allocated pod can start to reset gpu
	resetLock sync.Mutex
	// rwlock for ResetClient.resetInfo.Occupy
	resetInfoLock sync.RWMutex
}

func NewResetClient() *ResetClient {
	ki, err := NewKubeClient()
	if err != nil {
		klog.Errorf("Failed to create kube client: %v", err)
		os.Exit(1)
	}

	nodeName := os.Getenv("NODE_NAME")
	return &ResetClient{
		client:      ki,
		cmName:      ResetConfigName + nodeName,
		cmLabel:     "gpuReset",
		cmNamespace: ki.Namespace,
		resetInfo: &GpuResetInfo{
			NodeName: nodeName,
			Reset:    false,
			Occupy: map[string]bool{
				DevicePluginName: true,
			},
		},
	}
}

func (rc *ResetClient) syncWriteOccupy(k string, v bool) {
	rc.resetInfoLock.Lock()
	defer rc.resetInfoLock.Unlock()
	rc.resetInfo.Occupy[k] = v
}

func (rc *ResetClient) syncReadResetInfo() GpuResetInfo {
	rc.resetInfoLock.RLock()
	defer rc.resetInfoLock.RUnlock()

	snap := GpuResetInfo{
		NodeName: rc.resetInfo.NodeName,
		Reset:    rc.resetInfo.Reset,
		Occupy:   make(map[string]bool, len(rc.resetInfo.Occupy)),
	}
	for k, v := range rc.resetInfo.Occupy {
		snap.Occupy[k] = v
	}
	return snap
}

func (rc *ResetClient) InitCmInformer() {
	informerFactory := informers.NewSharedInformerFactory(rc.client.Client, 0)
	cmInformer := informerFactory.Core().V1().ConfigMaps().Informer()
	cmInformer.AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: rc.informerConfigmapFilter,
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				klog.Infof("add reset configmap: %v", obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				rc.updateResetInfo(newObj)
			},
			DeleteFunc: func(obj interface{}) {
				klog.Infof("delete reset configmap: %v", obj)
			},
		},
	})
	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)

	rc.CreateResetInfo()
}

func (rc *ResetClient) CreateResetInfo() {
	var resetInfoData []byte
	var err error
	if resetInfoData, err = yaml.Marshal(rc.syncReadResetInfo()); err != nil {
		klog.Errorf("Failed to marshal resetInfo: %v", err)
		return
	}

	resetInfoCM := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rc.cmName,
			Namespace: rc.cmNamespace,
		},
	}

	resetInfoCM.Data = map[string]string{
		rc.cmLabel: string(resetInfoData),
	}

	klog.Infof("write gpu reset info %v into cm: %s/%s.", rc.resetInfo, rc.cmName, rc.cmNamespace)

	err = wait.PollImmediate(UpdateInternalTime*time.Second, UpdateTimeoutTime*time.Second, func() (bool, error) {
		if err := rc.client.createOrUpdateDeviceCM(resetInfoCM); err != nil {
			return false, err
		}
		return true, nil
	})
	if err != nil {
		klog.Errorf("Failed to create resetInfo configmap: %v", err)
	}
}

func (rc *ResetClient) updateResetCm(resetInfoData []byte) error {
	ctx := context.Background()
	key := rc.cmLabel
	val := string(resetInfoData)

	backoff := PatchWaitTime * time.Millisecond
	for i := 0; i < RetryUpdateCount; i++ {
		cur, err := rc.client.Client.CoreV1().ConfigMaps(rc.cmNamespace).
			Get(ctx, rc.cmName, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				klog.Errorf("ConfigMap %s not found", rc.cmName)
				return err
			}
			klog.Warningf("get configmap failed: %v", err)
			time.Sleep(backoff)
			continue
		}

		patchObj := map[string]any{
			"metadata": map[string]any{"resourceVersion": cur.ResourceVersion},
			"data":     map[string]string{key: val},
		}
		b, _ := json.Marshal(patchObj)

		_, err = rc.client.Client.CoreV1().ConfigMaps(rc.cmNamespace).
			Patch(ctx, rc.cmName, types.StrategicMergePatchType, b, metav1.PatchOptions{})
		if err == nil {
			klog.Infof("patch configmap %s success", rc.cmName)
			return nil
		}
		if errors.IsConflict(err) {
			// retry with exponential backoff
			klog.Infof("patch conflict on %s, retrying...", rc.cmName)
			time.Sleep(backoff)
			if backoff < time.Second {
				backoff *= 2
			}
			continue
		}
		if errors.IsNotFound(err) {
			klog.Errorf("ConfigMap %s not found", rc.cmName)
			return err
		}
		klog.Warningf("patch configmap failed: %v, try again", err)
		time.Sleep(backoff)
	}
	return fmt.Errorf("patch configmap %s exceeded retries", rc.cmName)
}

func (rc *ResetClient) updateResetInfo(obj interface{}) {

	cm, ok := obj.(*v1.ConfigMap)
	if !ok {
		klog.Errorf("Failed to get configmap from informer")
		return
	}

	resetinfo := &GpuResetInfo{}

	if err := yaml.Unmarshal([]byte(cm.Data[rc.cmLabel]), resetinfo); err != nil {
		klog.Errorf("Failed to unmarshal resetInfo: %v", err)
		return
	}

	for k, v := range resetinfo.Occupy {
		if k == DevicePluginName {
			continue
		}
		rc.syncWriteOccupy(k, v)
	}

	klog.Info("Updating resetInfo from configmap: ", resetinfo)
	klog.Info("Updated resetInfo: ", rc.resetInfo)
}

func (rc *ResetClient) informerConfigmapFilter(obj interface{}) bool {
	cm, ok := obj.(*v1.ConfigMap)
	if !ok {
		klog.Errorf("Failed to get configmap from informer")
		return false
	}
	return rc.checkConfigMapIsResetConfig(cm)
}

func (rc *ResetClient) checkConfigMapIsResetConfig(cm *v1.ConfigMap) bool {
	return cm.Namespace == rc.cmNamespace && cm.Name == rc.cmName
}

func (rc *ResetClient) prepareForReset() error {
	klog.Info("prepareForReset")
	var resetInfoData []byte
	var err error

	rc.resetInfo.Reset = true
	rc.syncWriteOccupy(DevicePluginName, false)

	if resetInfoData, err = yaml.Marshal(rc.syncReadResetInfo()); err != nil {
		klog.Errorf("Failed to marshal resetInfo: %v", err)
		return err
	}

	return rc.updateResetCm(resetInfoData)
}

func (rc *ResetClient) doneForReset() error {
	var resetInfoData []byte
	var err error

	rc.resetInfo.Reset = false
	rc.syncWriteOccupy(DevicePluginName, true)

	if resetInfoData, err = yaml.Marshal(rc.syncReadResetInfo()); err != nil {
		klog.Errorf("Failed to marshal resetInfo: %v", err)
		return err
	}

	rc.updateResetCm(resetInfoData)

	for {
		isRecover := true
		stillDown := []string{}
		for k, v := range rc.syncReadResetInfo().Occupy {
			if !v {
				isRecover = false
				stillDown = append(stillDown, k)
			}
		}

		if !isRecover {
			klog.Info("waiting for other to recover...")
			klog.Infof("Still down: %v", stillDown)
		} else {
			break
		}
		time.Sleep(1 * time.Second)
	}

	klog.Info("doneForReset")
	return nil
}

func (rc *ResetClient) ResetGpus(indexs []int) error {
	rc.resetLock.Lock()
	klog.Info("reset gpu locking")
	defer func() {
		rc.resetLock.Unlock()
		klog.Info("reset gpu unlocking")
	}()

	klog.Info("Start Reset Process")
	klog.Info("Shutdown of IXML for reset gpu returned:", ixml.Shutdown())

	err := rc.prepareForReset()
	if err != nil {
		return err
	}

	for {
		isCompleted := true
		stillOccupy := []string{}

		for k, v := range rc.syncReadResetInfo().Occupy {
			if v {
				isCompleted = false
				stillOccupy = append(stillOccupy, k)
			}
		}

		if isCompleted {
			break
		}
		klog.Infof("Waiting for other plugins to complete reset gpu...")
		klog.Infof("Still occupy reset gpu by %v", stillOccupy)
		time.Sleep(5 * time.Second)
	}

	sindex := ""
	for i, index := range indexs {
		sindex += strconv.Itoa(index)
		if i != len(indexs)-1 {
			sindex += ","
		}
	}
	cmd := exec.Command("/usr/local/corex/bin/ixsmi", "-r", "-i", sindex)

	output, cmdErr := cmd.Output()
	if cmdErr != nil {
		klog.Errorf("Failed to reset gpu: %v", cmdErr)
	}
	klog.Infof("reset gpu output: %s", output)

	err = rc.doneForReset()
	if err != nil {
		return err
	}

	klog.Info("Loading IXML")
	err = ixml.Init()
	if err != nil {
		klog.Errorf("Failed to initialize IXML: %v", err)
		return fmt.Errorf("%v", err)
	} else {
		klog.Info("IXML load success")
	}

	return cmdErr
}
