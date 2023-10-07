# kubectr


`kubectr` is a utility that list all containers in a pod, supporting both `initContainers` and `ephemeralContainers`.

# Demo
```
$ kubectl ctr csi-do-controller-0 -n kube-system

NAME           	  READY	  STATUS 	  RESTARTS     	  AGE  	  PORTS	  IMAGE                                             	  PULLPOLICY  	  TYPE
csi-provisioner	  1    	  Running	  4            	  6h10m	  -    	  registry.k8s.io/sig-storage/csi-provisioner:v3.5.0	  IfNotPresent	  container
csi-attacher   	  1    	  Running	  4 (6h ago)   	  6h10m	  -    	  registry.k8s.io/sig-storage/csi-attacher:v4.3.0   	  IfNotPresent	  container
csi-snapshotter	  1    	  Running	  4 (6h ago)   	  6h10m	  -    	  registry.k8s.io/sig-storage/csi-snapshotter:v6.2.2	  IfNotPresent	  container
csi-resizer    	  1    	  Running	  4 (6h ago)   	  6h9m 	  -    	  registry.k8s.io/sig-storage/csi-resizer:v1.8.0    	  IfNotPresent	  container
csi-do-plugin  	  0    	  Waiting	  1069 (5m ago)	  -    	  -    	  digitalocean/do-csi-plugin:v4.7.1                 	  Always      	  container
```

# Installation
There are several ways to install kubectr. The recommended installation method is via brew.

## Via krew

Krew is a `kubectl` plugin manager. If you have not yet installed `krew`, get it at [https://github.com/kubernetes-sigs/krew](https://github.com/kubernetes-sigs/krew). Then installation is as simple as

```
kubectl krew install ctr
```

## Binaries
download binaries file from https://github.com/cfanbo/kubectr/releases

## From Sourcee

```shell
git clone https://github.com/cfanbo/kubectr.git
make 
bin/kubectr
```

# Usage

## As kubectl plugin
Most users will have installed `kubectr` via [krew](https://github.com/kubernetes-sigs/krew), so the plugin is already correctly installed. Otherwise, rename `kubectr` to `kubectl-ctr` and put it in some directory from your $PATH variable. Then you can invoke the plugin via `kubectl ctr`

## Standalone
Put the `kubectr` binary in some directory from your `$PATH` variable. For example
```
sudo mv -i kubectr /usr/bin/kubectr
```
Then you can invoke the plugin via `kubectr`