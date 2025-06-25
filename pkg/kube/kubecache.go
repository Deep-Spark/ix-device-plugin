package kube

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

var (
	nodeDeviceInfoCache *NodeDeviceInfoCache
	podCache            = map[types.UID]*podInfo{}
	lock                = sync.Mutex{}
)

type podInfo struct {
	*v1.Pod
	updateTime time.Time
}

func (ki *KubeClient) InitPodInformer() {
	factory := informers.NewSharedInformerFactoryWithOptions(ki.Client, 0,
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.FieldSelector = "spec.nodeName=" + ki.NodeName
		}))
	podInformer := factory.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			UpdatePodList(nil, obj, EventTypeAdd)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			if !reflect.DeepEqual(oldObj, newObj) {
				UpdatePodList(oldObj, newObj, EventTypeUpdate)
			}
		},
		DeleteFunc: func(obj interface{}) {
			UpdatePodList(nil, obj, EventTypeDelete)
		},
	})
	podInformer.SetWatchErrorHandler(func(r *cache.Reflector, err error) {
		klog.Errorf("pod informer watch error: %v", err)
	})
	factory.Start(make(chan struct{}))

	cache.WaitForCacheSync(wait.NeverStop, podInformer.HasSynced)

	ki.PodInformer = podInformer
}

func UpdatePodList(oldObj, newObj interface{}, operator EventType) {
	newPod, ok := newObj.(*v1.Pod)
	if !ok {
		return
	}
	lock.Lock()
	defer lock.Unlock()
	switch operator {
	case EventTypeAdd, EventTypeUpdate:
		// klog.Infof("pod(%s/%s) is %s to cache", newPod.Namespace, newPod.Name, operator)
		podCache[newPod.UID] = &podInfo{
			Pod:        newPod,
			updateTime: time.Now(),
		}
	case EventTypeDelete:
		// klog.Infof("pod(%s/%s) is deleted from cache", newPod.Namespace, newPod.Name)
		delete(podCache, newPod.UID)
	default:
		klog.Errorf("operator is undefined, find operater: %s", operator)
	}
}

func (ki *KubeClient) GetActivePodListCache() []v1.Pod {
	newPodList := make([]v1.Pod, 0)
	lock.Lock()
	defer lock.Unlock()
	for _, pi := range podCache {
		if pi.Status.Phase == v1.PodFailed || pi.Status.Phase == v1.PodSucceeded {
			continue
		}
		newPodList = append(newPodList, *pi.Pod)
	}

	return newPodList
}

func (ki *KubeClient) GetActivePodList() ([]v1.Pod, error) {
	fieldSelector, err := fields.ParseSelector("spec.nodeName=" + ki.NodeName + "," +
		"status.phase!=" + string(v1.PodSucceeded) + ",status.phase!=" + string(v1.PodFailed))
	if err != nil {
		return nil, err
	}
	podList, err := ki.getPodListByCondition(fieldSelector)
	if err != nil {
		return nil, err
	}
	return checkPodList(podList)
}

func checkPodList(podList *v1.PodList) ([]v1.Pod, error) {
	var pods = make([]v1.Pod, 0)
	for _, pod := range podList.Items {
		pods = append(pods, pod)
	}
	return pods, nil
}

func (ki *KubeClient) getPodListByCondition(selector fields.Selector) (*v1.PodList, error) {
	newPodList, err := ki.Client.CoreV1().Pods(v1.NamespaceAll).List(context.Background(), metav1.ListOptions{
		FieldSelector:   selector.String(),
		ResourceVersion: "0",
	})

	return newPodList, err
}

func checkAnnotationAllocateValid(requestDevices []string, pod *v1.Pod) bool {
	if predicateTime, ok := pod.Annotations[PodPredicateTime]; ok {
		if predicateTime == strconv.FormatUint(math.MaxUint64, 10) {
			klog.Warningf("The pod has been mounted to a device, pod name: %s", pod.Name)
			return false
		}
	}
	devStr, exist := pod.Annotations[ResourceNamePrefix+PodDevVolcano]
	if !exist {
		return false
	}

	allocateDevice := strings.Split(devStr, CommaSepDev)
	return len(allocateDevice) == len(requestDevices)
}

func isShouldDeletePod(pod *v1.Pod) bool {
	if pod.DeletionTimestamp != nil {
		return true
	}
	for _, status := range pod.Status.ContainerStatuses {
		if status.State.Waiting != nil &&
			strings.Contains(status.State.Waiting.Message, "PreStartContainer check failed") {
			return true
		}
	}
	return pod.Status.Reason == "UnexpectedAdmissionError"
}

func FilterPods(pods []v1.Pod, conditionFunc func(pod *v1.Pod) bool) []v1.Pod {
	var res = make([]v1.Pod, 0)
	for _, pod := range pods {
		if isShouldDeletePod(&pod) {
			continue
		}
		if conditionFunc != nil && !conditionFunc(&pod) {
			continue
		}
		res = append(res, pod)
	}
	return res
}

func (ki *KubeClient) GetMatchedPod(requestDevices []string) (*v1.Pod, error) {
	conditionFunc := func(pod *v1.Pod) bool {
		return checkAnnotationAllocateValid(requestDevices, pod)
	}
	var filteredPods = make([]v1.Pod, 0)
	var allPods = make([]v1.Pod, 0)
	for i := 0; i < GetPodFromInformerTime; i++ {
		if i == GetPodFromInformerTime-1 {
			// in the last time of retry, get the pod from api server instead of cache
			noneCachedPod, err := ki.GetActivePodList()
			if err != nil {
				klog.Errorf("get active pod from api server failed")
				return nil, err
			}
			allPods = noneCachedPod
		} else {
			allPods = ki.GetActivePodListCache()
		}
		filteredPods = FilterPods(allPods, conditionFunc)
		if len(filteredPods) != 0 {
			break
		}
		klog.Warningf("no pod passed the filter, request device: %v, retry: %d", requestDevices, i)
		time.Sleep(time.Second)
	}
	oldestPod := getOldestPod(filteredPods)
	if oldestPod == nil {
		return nil, fmt.Errorf("not get valid pod")
	}

	return oldestPod, nil
}

func getOldestPod(pods []v1.Pod) *v1.Pod {
	if len(pods) == 0 {
		return nil
	}
	oldest := pods[0]
	for _, pod := range pods {
		if getPredicateTimeFromPodAnnotation(&oldest) > getPredicateTimeFromPodAnnotation(&pod) {
			oldest = pod
		}
	}
	klog.Infof("oldest pod %#v, predicate time: %#v", oldest.Name, oldest.Annotations[PodPredicateTime])
	return &oldest
}

func getPredicateTimeFromPodAnnotation(pod *v1.Pod) uint64 {
	assumeTimeStr, ok := pod.Annotations[PodPredicateTime]
	if !ok {
		klog.Warningf("volcano not write timestamp, pod Name: %s", pod.Name)
		return math.MaxUint64
	}
	predicateTime, err := strconv.ParseUint(assumeTimeStr, 10, 64)
	if err != nil {
		klog.Errorf("parse timestamp failed, %v", err)
		return math.MaxUint64
	}
	return predicateTime
}
