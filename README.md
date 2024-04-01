# IX device plugin for Kubernetes

## Overview

The IX device plugin is a Daemonset for Kubernetes, which can help to expose the Iluvatar GPU in the Kubernetes cluster.

## Build

Before building the IX device plugin, it's mandatory to prepare `Corex SDK`, the default `COREX SDK` path is `/usr/local/corex/`.

Make sure `golang >= 1.11` and build the IX device plugin as follows:

```shell
make plugin
```

or

```shell
make
```


## Deployment

Once the Kubernetes cluster is ready, you can enable GPU support by deploying the following Daemonset:

```shell
kubectl create -f ix-device-plugin.yaml
```

## Example

GPU can be exposed to a pod by adding `iluvatar.ai/gpu` to the pod definition, and you can restrict the GPU resource by adding `resources.limits` to the pod definition. Example following:

```shell
kubectl create -f corex-example.yaml
```

## License

Copyright (c) 2024 Iluvatar CoreX. All rights reserved. This project has an Apache-2.0 license, as found in the [LICENSE](LICENSE) file.
