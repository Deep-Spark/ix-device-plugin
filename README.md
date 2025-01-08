# IX device plugin for Kubernetes

## Table of Contents
- [About](#about)
- [Prerequisites](#prerequisites)
- [Building the IX device plugin](#building-the-ix-device-plugin)
- [Configuring the IX device plugin](#configuring-the-ix-device-plugin)
- [Enabling GPU Support in Kubernetes](#enabling-gpu-support-in-kubernetes)
- [Running GPU Jobs](#running-gpu-jobs)
- [Split GPU Board to Multiple GPU Devices](#split-gpu-board-to-multiple-gpu-devices)
- [Shared Access to GPUs](#shared-access-to-gpus)

## About

The IX device plugin for Kubernetes is a Daemonset that allows you to automatically:
- Expose the number of GPUs on each nodes of your cluster
- Keep track of the health of your GPUs
- Run GPU enabled containers in your Kubernetes cluster.

## Prerequisites

The list of prerequisites for running the IX device plugin is described below:
* Iluvatar driver and software stack >= v1.1.0
* Kubernetes version >= 1.10

## Building the IX device plugin

```shell
make all
```
This will build the ix-device-plugin binary and ix-device-plugin image, see logging for more details.

## Configuring the IX device plugin

The IX device plugin has a number of options that can be configured for it.
These options can be configured via a config file when launching the device plugin. Here we explain what
each of these options are and how to configure them in configmap.
```yaml
# ix-config.yaml
apiVersion: v1
kind: ConfigMap
data:
ix-config: |-
    version: "4.2.0"
    flags:
      splitboard: false
    sharing:
      timeSlicing:
          replicas: 4 

metadata:
name: ix-config
namespace: kube-system
```
```shell
kubectl create -f ix-config.yaml
```
| `Field`|        `Type `               |   `Description` |
|--------|------------------------------|------------------|
| `flags.splitboard`       | boolean  | Split GPU devices in every board(eg.BI-V150) if `splitboard` is `true`|
| `sharing.timeSlicing.replicas`       | integer  | Specifies the number of GPU time-slicing ​​replicas for shared access|

## Enabling GPU Support in Kubernetes

Once you have configured the options above on all the GPU nodes in your
cluster, you can enable GPU support by deploying the following Daemonset:
```yaml
# ix-device-plugin.yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: iluvatar-device-plugin
  namespace: kube-system
  labels:
    app.kubernetes.io/name: iluvatar-device-plugin
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: iluvatar-device-plugin
  template:
    metadata:
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ""
      labels:
        app.kubernetes.io/name: iluvatar-device-plugin
    spec:
      priorityClassName: "system-node-critical"
      securityContext:
        null
      containers:
        - name: iluvatar-device-plugin
          securityContext:
            capabilities:
              drop:
              - ALL
            privileged: true
          image: "ix-device-plugin:4.2.0"
          imagePullPolicy: IfNotPresent
          livenessProbe:
            exec:
              command:
              - ls
              - /var/lib/kubelet/device-plugins/iluvatar-gpu.sock
            periodSeconds: 5
          startupProbe:
            exec:
              command:
              - ls
              - /var/lib/kubelet/device-plugins/iluvatar-gpu.sock
            periodSeconds: 5
          resources:
            {}
          volumeMounts:
            - mountPath: /var/lib/kubelet/device-plugins
              name: device-plugin
            - mountPath: /run/udev
              name: udev-ctl
              readOnly: true
            - mountPath: /sys
              name: sys
              readOnly: true
            - mountPath: /dev
              name: dev
            - name: ixc
              mountPath: /ixconfig
      volumes:
        - hostPath:
            path: /var/lib/kubelet/device-plugins
          name: device-plugin
        - hostPath:
            path: /run/udev
          name: udev-ctl
        - hostPath:
            path: /sys
          name: sys
        - hostPath:
            path: /etc/udev/
          name: udev-etc
        - hostPath:
            path: /dev
          name: dev
        - name: ixc
          configMap:
              name: ix-config
```
```shell
kubectl create -f ix-device-plugin.yaml
```

## Running GPU Jobs

GPU can be exposed to a pod by adding `iluvatar.com/gpu` to the pod definition, and you can restrict the GPU resource by adding `resources.limits` to the pod definition.

```yaml
$ cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: corex-example
spec:
  containers:
  - name: corex-example
    image: corex:4.0.0
    command: ["/usr/local/corex/bin/ixsmi"]
    args: ["-l"]
    resources:
      limits:
        iluvatar.com/gpu: 1 # requesting 1 GPUs
EOF
```

```shell
kubectl logs corex-example
+-----------------------------------------------------------------------------+
|  IX-ML: 4.0.0       Driver Version: 4.1.0       CUDA Version: N/A           |
|-------------------------------+----------------------+----------------------|
| GPU  Name                     | Bus-Id               | Clock-SM  Clock-Mem  |
| Fan  Temp  Perf  Pwr:Usage/Cap|      Memory-Usage    | GPU-Util  Compute M. |
|===============================+======================+======================|
| 0    Iluvatar BI-V150S        | 00000000:8A:00.0     | 500MHz    1600MHz    |
| 0%   33C   P0    N/A / N/A    | 114MiB / 32768MiB    | 0%        Default    |
+-------------------------------+----------------------+----------------------+

+-----------------------------------------------------------------------------+
| Processes:                                                       GPU Memory |
|  GPU        PID      Process name                                Usage(MiB) |
|=============================================================================|
|  No running processes found                                                 |
+-----------------------------------------------------------------------------+
```

## Split GPU Board to Multiple GPU Devices

The IX device plugin allows splitting one GPU board into multiple GPU Devices through a set of
extended options in its configuration file. 

### With SplitBoard

The extended options for splitting board can be seen below:

```yaml
version: "4.2.0"
flags:
    splitboard: false
```

That is, `flags.splitboard`, a boolean flag can now be specified. If this flag is set to true, the plugin will split the GPU board into multiple GPUs and
kubelet will advertise multiple `iluvatar.com/gpu` resources to Kubernetes instead of 1 for one GPU board. Otherwise, the plugin will advertise only 1 `iluvatar.com/gpu` resource for one GPU board.

For example:

```yaml
version: "4.2.0"
flags:
    splitboard: true
```

If this configuration were applied to a node with 1 GPUs(eg. Bi-V150, which has 2 GPU chips on it) on it, the plugin
would now advertise 2 `iluvatar.com/gpu` resources to Kubernetes instead of 1.

```
$ kubectl describe node
...
Capacity:
  iluvatar.com/gpu: 2
...
```

## Shared Access to GPUs

The IX device plugin allows oversubscription of GPUs through a set of
extended options in its configuration file. 

### With Time-Slicing

The extended options for sharing using time-slicing can be seen below:

```yaml
version: "4.2.0"
sharing:
    timeSlicing:
        replicas: <num-replicas>
    ...
```

That is, `sharing.timeSlicing.replicas`, a number of replicas can now be specified. These replicas represent the number of shared accesses that will be granted for a GPU.

For example:

```yaml
version: "4.2.0"
flags:
    splitboard: false
sharing:
    timeSlicing:
        replicas: 4
```

If this configuration were applied to a node with 2 GPUs on it, the plugin
would now advertise 8 `iluvatar.com/gpu` resources to Kubernetes instead of 2.

```
$ kubectl describe node
...
Capacity:
  iluvatar.com/gpu: 8
...
```
