<!--
SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
SPDX-License-Identifier: Apache-2.0
-->

# Deploying onos-topo

This guide deploys `onos-topo` through its [Helm] chart assumes you have a [Kubernetes] cluster running 
with an atomix controller deployed in a namespace.
`onos-topo` Helm chart is based on Helm 3.0 version, with no need for the Tiller pod to be present. 
If you don't have a cluster running and want to try on your local machine please follow first 
the [Kubernetes] setup steps outlined in [deploy with Helm](https://docs.onosproject.org/developers/deploy_with_helm/).
The following steps assume you have the setup outlined in that page, including the `micro-onos` namespace configured. 

## Installing
To install the `onos-topo` published chart in the `micro-onos` namespace run:
```bash
$ helm -n micro-onos install onos-topo onosproject/onos-topo
```
Note, that this assumes that `https://charts.onosproject.org` was added as `onosproject` Helm repo.

Alternately, you can deploy your own version of the chart from the root directory of the `onos-helm-charts`
repo via the following:
```bash
$ helm install -n micro-onos onos-uenib onos-uenib
```
In either case, the output of the `helm install` command should look similar to the following:
```bash
NAME: onos-topo
LAST DEPLOYED: Wed Jul 14 14:30:27 2021
NAMESPACE: micro-onos
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

`helm install` assigns a unique name to the chart and displays all the k8s resources that were
created by it. To list the charts that are installed and view their statuses, run `helm ls`:

```bash
$ helm -n micro-onos ls
NAME     	NAMESPACE 	REVISION	UPDATED                            	STATUS  	CHART           	APP VERSION
onos-topo	micro-onos	1       	2021-07-14 14:30:27.81784 -0700 PDT	deployed	onos-topo-1.0.15	v0.7.9
```

### Partition Set
The `onos-topo` chart also deploys a custom Atomix `PartitionSet` resource to store all the 
topology information in a replicated and fail-safe manner. 
In the following example there is only one partition set deployed `onos-topo-consensus-store-1-0`.

```bash
$ kubectl -n micro-onos get pods
NAME                            READY   STATUS    RESTARTS   AGE
onos-topo-854bf799f-lrnkd       3/3     Running   0          2m42s
onos-topo-consensus-store-1-0   1/1     Running   0          2m41s
```

One can customize the number of partitions and replicas by modifying, in `values.yaml`, under `store/raft` 
the values of 
```bash 
partitions: 1
partitionSize: 1
```

## Uninstalling

To uninstall the `onos-topo` chart, run the following:
```bash
 helm delete -n micro-onos onos-topo
```

## Troubleshooting

Helm offers two flags to help you debug your chart. This can be useful if your chart does not install,
the pod is not running for some reason, or you want to trouble-shoot custom configuration values,

* `--dry-run` check the chart without actually installing the pod.
* `--debug` prints out more information about your chart

```bash
helm install -n micro-onos onos-topo --debug --dry-run onos-topo/
```

[Helm]: https://helm.sh/
[Kubernetes]: https://kubernetes.io/
[kind]: https://kind.sigs.k8s.io

