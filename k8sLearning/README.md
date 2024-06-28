# k8s learning

## 特点

自我修复：发现一个容器挂了马上找个替代品  

弹性伸缩：我设置数字 他就给我创建多少个容器  

自动部署和回滚：让用户感受不到变化  

负载均衡： 可以不需要nginx了

配置管理



![](img/截屏2024-06-20%2011.00.27.png)
  

**集群架构：一个集群搭建完会至少包含下面这些组件**

![](img/截屏2024-06-20%2011.18.18.png)


客户端主要有两种
- 命令行 `kubectl`
- 一些市面上的ui `可视化界面`

![](img/截屏2024-06-20%2014.59.51.png)

 


# K8s in action

## chapter 3

Therefore, you need to run each process in its own container. That’s how Docker
and Kubernetes are meant to be used.

One thing to stress here is that because containers in a pod run in the same Network namespace, they share the same IP address and port space. 

Containers of different pods can never run into port conflicts, because each pod has a separate port space. 

pods are logical hosts and behave much like physical hosts or VMs in the non-container world. Processes running in the same pod are like processes running on the same physical or virtual machine, except that each process is encapsulated in a container.

Kubernetes can’t horizontally scale individual contain- ers; instead, it scales whole pods.

### Using kubectl create to create the pod

```sh
kubectl create -f .yaml
kubectl get po zilin-gohttp-base-5ddc8b54d4-nts4q -o yaml
kubectl get po zilin-gohttp-base-5ddc8b54d4-nts4q -o json
kubectl get pod zilin-gohttp-base-5ddc8b54d4-nts4q -o json | jq -r '.status.containerStatuses[].containerID'
kubectl get pod zilin-gohttp-base-5ddc8b54d4-nts4q -o json | jq -r '.spec.containers[].name'
kubectl logs zilin-gohttp-base-5ddc8b54d4-nts4q #print the output log of container
```

### Sending requests to the pod

```sh
➜  ~ kubectl port-forward zilin-gohttp-base-5ddc8b54d4-nts4q 8888:8080
Forwarding from 127.0.0.1:8888 -> 8080
Forwarding from [::1]:8888 -> 8080


curl localhost:8888

```

### labels

```sh
kubectl get pod --show-labels
NAME              READY   STATUS            RESTARTS   AGE   LABELS
kubia-manual-v2   0/2     PodInitializing   0          34s   creation_method=manual,env=prod,security.istio.io/tlsMode=istio,service.istio.io/canonical-name=kubia-manual-v2,service.istio.io/canonical-revision=latest


➜  yaml kubectl get pod -L creation_method,env
NAME              READY   STATUS             RESTARTS   AGE     CREATION_METHOD   ENV
kubia-manual-v2   1/2     ImagePullBackOff   0          2m29s   manual            prod



# add label 
➜  yaml kubectl label pod kubia-manual-v1 creation_method=manual
pod/kubia-manual-v1 labeled
➜  yaml kubectl get pod -L creation_method,env
NAME              READY   STATUS             RESTARTS   AGE   CREATION_METHOD   ENV
kubia-manual-v1   1/2     ImagePullBackOff   0          10m   manual
kubia-manual-v2   1/2     ImagePullBackOff   0          15m   manual            prod

# change label
➜  yaml kubectl label pod kubia-manual-v2 env=debug --overwrite
pod/kubia-manual-v2 labeled
➜  yaml kubectl get pod -L creation_method,env
NAME              READY   STATUS             RESTARTS   AGE   CREATION_METHOD   ENV
kubia-manual-v1   1/2     ImagePullBackOff   0          12m   manual
kubia-manual-v2   1/2     ImagePullBackOff   0          16m   manual            debug
```

### 3.4 Listing subsets of pods through label selectors

```sh
➜  yaml kubectl label pod kubia-manual-v1 creation_method= --overwrite
pod/kubia-manual-v1 labeled
➜  yaml kubectl get pod -L creation_method,env
NAME              READY   STATUS             RESTARTS   AGE   CREATION_METHOD   ENV
kubia-manual-v1   1/2     ImagePullBackOff   0          17m
kubia-manual-v2   1/2     ImagePullBackOff   0          22m   manual            debug

➜  yaml kubectl get pod -l creation_method=manual
NAME              READY   STATUS             RESTARTS   AGE
kubia-manual-v2   1/2     ImagePullBackOff   0          23m

➜  yaml kubectl get pod -l env
NAME              READY   STATUS             RESTARTS   AGE
kubia-manual-v2   1/2     ImagePullBackOff   0          23m

➜  yaml kubectl get pod -l '!env'
NAME              READY   STATUS             RESTARTS   AGE
kubia-manual-v1   1/2     ImagePullBackOff   0          20m
```

### Using labels and selectors to constrain pod scheduling

```sh
➜  yaml kubectl label node cn-zhangjiakou.192.168.104.100 gpu=true
node/cn-zhangjiakou.192.168.104.100 labeled
➜  yaml kubectl get nodes -l gpu=true
NAME                             STATUS   ROLES    AGE   VERSION
cn-zhangjiakou.192.168.104.100   Ready    <none>   59d   v1.24.6-aliyun.1

# know the pod's belong node
kubectl get pod kubia-manual-v1 -o wide
```

### Annotating pods

Annotations are also key-value pairs, so in essence, they’re similar to labels, but they aren’t meant to hold identifying information. 


A great use of annotations is adding descriptions for each pod or other API object, so that everyone using the cluster can quickly look up information about each individ- ual object

**get the annotations of a pod**

```sh
kubectl get pods kubia-manual-v1 -o json | jq '.metadata.annotations'
```

```sh
➜  yaml kubectl get pods kubia-manual-v1 -o json | jq '.metadata.annotations'
{
  "k8s.aliyun.com/pod-ips": "192.168.162.160",
  "kubectl.kubernetes.io/default-container": "kubia",
  "kubectl.kubernetes.io/default-logs-container": "kubia",
  "kubernetes.io/psp": "ack.privileged",
  "prometheus.io/path": "/stats/prometheus",
  "prometheus.io/port": "15020",
  "prometheus.io/scrape": "true",
  "sidecar.istio.io/status": "{\"initContainers\":[\"istio-init\"],\"containers\":[\"istio-proxy\"],\"volumes\":[\"istio-envoy\",\"istio-data\",\"istio-podinfo\",\"istio-token\",\"istiod-ca-cert\"],\"imagePullSecrets\":null,\"revision\":\"default\"}"
}
➜  yaml kubectl annotate pod kubia-manual-v1 description="this is zilin learning"
pod/kubia-manual-v1 annotated
➜  yaml kubectl get pods kubia-manual-v1 -o json | jq '.metadata.annotations'
{
  "description": "this is zilin learning",
  "k8s.aliyun.com/pod-ips": "192.168.162.160",
  "kubectl.kubernetes.io/default-container": "kubia",
  "kubectl.kubernetes.io/default-logs-container": "kubia",
  "kubernetes.io/psp": "ack.privileged",
  "prometheus.io/path": "/stats/prometheus",
  "prometheus.io/port": "15020",
  "prometheus.io/scrape": "true",
  "sidecar.istio.io/status": "{\"initContainers\":[\"istio-init\"],\"containers\":[\"istio-proxy\"],\"volumes\":[\"istio-envoy\",\"istio-data\",\"istio-podinfo\",\"istio-token\",\"istiod-ca-cert\"],\"imagePullSecrets\":null,\"revision\":\"default\"}"
}
```
### namespace

Using multiple namespaces allows you to split complex systems with numerous components into smaller distinct groups.

```sh
➜  yaml kubectl get ns
NAME                          STATUS   AGE
arms-prom                     Active   698d
canary-operator               Active   728d
csdr                          Active   2y1d
default                       Active   2y95d
do                            Active   657d
elastic-env-operator-system   Active   460d
hw-mcp                        Active   643d
istio-system                  Active   638d
karmada-cluster               Active   644d
keda                          Active   144d
knative-eventing              Active   2y3d
knative-serving               Active   2y3d
knative-sources               Active   2y3d
kube-node-lease               Active   2y95d
kube-public                   Active   2y95d
kube-system                   Active   2y95d
kubevela                      Active   634d
prod                          Active   670d
spire                         Active   19d
spire-system                  Active   19d
sqb                           Active   728d
sreworks-client               Active   2y95d
vela-system                   Active   634d
➜  yaml kubectl get pod --namespace sqb
NAME              READY   STATUS             RESTARTS   AGE
kubia-gpu         1/2     ImagePullBackOff   0          32m
kubia-manual-v1   1/2     ImagePullBackOff   0          122m
kubia-manual-v2   1/2     ImagePullBackOff   0          127m
```


### stop and remove pods 

```sh
➜  yaml kubectl get pods
NAME              READY   STATUS             RESTARTS   AGE
kubia-gpu         1/2     ImagePullBackOff   0          65m
kubia-manual-v1   1/2     ImagePullBackOff   0          156m
kubia-manual-v2   1/2     ImagePullBackOff   0          160m
➜  yaml kubectl delete pods kubia-gpu
pod "kubia-gpu" deleted
➜  yaml kubectl get pods
NAME              READY   STATUS             RESTARTS   AGE
kubia-manual-v1   1/2     ImagePullBackOff   0          159m
kubia-manual-v2   1/2     ImagePullBackOff   0          164m
```

**delete with labels**

```sh
➜  yaml kubectl get pods -L creation_method,env
NAME              READY   STATUS             RESTARTS   AGE    CREATION_METHOD   ENV
kubia-manual-v1   1/2     ImagePullBackOff   0          161m
kubia-manual-v2   1/2     ImagePullBackOff   0          165m   manual            debug
➜  yaml kubectl delete pods -l creation_method=manual
pod "kubia-manual-v2" deleted
➜  yaml kubectl get pods
NAME              READY   STATUS             RESTARTS   AGE
kubia-manual-v1   1/2     ImagePullBackOff   0          162m
```

## Chapter 5

The service address doesn’t change even if the pod’s IP address changes. Additionally, by creating the service, you also enable the frontend pods to easily find the backend service by its name through either environment variables or DNS.

```sh
yaml kubectl create -f kubia-svc.yaml
service/kubia created
➜  yaml kubectl get svc
NAME                           TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
app-edge                       ClusterIP   172.16.205.35    <none>        80/TCP           728d
boss-circle-op                 ClusterIP   172.16.181.235   <none>        80/TCP,443/TCP   725d
details                        ClusterIP   172.16.167.204   <none>        9080/TCP         711d
emenu-h5                       ClusterIP   172.16.163.6     <none>        8080/TCP         725d
heart-magic                    ClusterIP   172.16.97.70     <none>        80/TCP           725d
jfrog-flow-test                ClusterIP   172.16.10.200    <none>        80/TCP           460d
kubia                          ClusterIP   172.16.232.114   <none>        80/TCP           4s
```


You can send requests to your service from within the cluster in a few ways:
 The obvious way is to create a pod that will send the request to the service’s cluster IP and log the response. You can then examine the pod’s log to see what the service’s response was.
 You can ssh into one of the Kubernetes nodes and use the curl command.
 You can execute the curl command inside one of your existing pods through
the kubectl exec command.




