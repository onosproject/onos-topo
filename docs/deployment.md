# Deploying onos-topo

This guide deploys `onos-topo` through it's [Helm] chart assumes you have a [Kubernets] cluster running 
with an atomix controller deployed in a namespace.  
If you dont' specify the `--namespace` in the commands below atomix controller must be deployed in the `default`
`onos-topo` Helm chart is based on Helm 3.0 version, with no need for the Tiller pod to be present. 
If you don't have a cluster running and want to try on your local machine please follow first 
the [Kubernetes] setup steps outlined in [deploy with Helm](https://docs.onosproject.org/developers/deploy_with_helm/).

## Installing the Chart

To install the chart, simply run `helm install deployments/helm/onos-topo` from
the root directory of this project:

```bash
helm install onos-topo deployments/helm/onos-topo
NAME: onos-topo
LAST DEPLOYED: Tue Nov 26 13:31:42 2019
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

`helm install` assigns a unique name to the chart and displays all the k8s resources that were
created by it. To list the charts that are installed and view their statuses, run `helm ls`:

```bash
helm ls
NAME          	REVISION	UPDATED                 	STATUS  	CHART                    	APP VERSION	NAMESPACE
...
onos-topo	1       	Tue May 14 18:56:39 2019	DEPLOYED	onos-topo-0.0.1	        0.0.1      	default
```

### Onos Topo Partition Set

The `onos-topo` chart also deployes a `PartitionSet` custom Atomix resource to store all the 
configuration in a replicated and fail safe manner. 
In the following example there is only one partition set deployed
`onos-topo-1-0`.

```bash
NAMESPACE     NAME                                         READY   STATUS    RESTARTS   AGE
default       atomix-controller-b579b9f48-lgvxf            1/1     Running   0          63m
default       onos-topo-1-0                              1/1     Running   0          61m
default       onos-topo-77765c9dc4-vsjjn                 1/1     Running   0          61m
```

One can customize the number of partitions and replicas by modifying, in `values.yaml`, under `store/raft` 
the values of 
```bash 
partitions: 1
partitionSize: 1
```

### Installing the chart in a different namespace.

To install the chart in a different namespace please modify the `default` occurances in the `values.yaml` file.   
Please be aware to change also `atomix-controller.default.svc.cluster.local:5679` 
to `atomix-controller.<your_name_space_here>.svc.cluster.local:5679`.  
Then issue the `helm install` command
```bash
helm install onos-topo --namespace <your_name_space> deployments/helm/onos-topo
```
### Installing the chart with debug. 
`onos-topo` offers the capability to open a debug port (4000) to the image.
To enable the debug capabilities please set the debug flag to true in `values.yaml` or pass it to `helm install`
```bash
helm install onos-topo deployments/helm/onos-topo --set debug=true
```

### Troubleshoot

If your chart does not install or the pod is not running for some reason and/or you modified values Helm offers two flags to help you
debug your chart:  

* `--dry-run` check the chart without actually installing the pod. 
* `--debug` prints out more information about your chart

```bash
helm install onos-topo --debug --dry-run ./deployments/helm/onos-topo/
```
## Uninstalling the chart.

To remove the `onos-topo` pod issue
```bash
 helm delete onos-topo
```
## Pod Information

To view the pods that are deployed, run `kubectl get pods`.

[Helm]: https://helm.sh/
[Kubernetes]: https://kubernetes.io/
[kind]: https://kind.sigs.k8s.io

