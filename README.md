# podchaosmonkey

Application to test resilience in services deployed in Kubernetes.

## Description

The program must runs inside the cluster, interacts with the kube-apiserver, and deletes on pod at random in a particular namespace on a schedule.

It is assumed that a schedule is a repetitive event over the time with a specific frequency.

### App parameters

| Parameter | Type | Default value | Description |
| --------- | ---- | ------------- | ----------- |
| `targetNamespace` | String | `workloads` | Namespace used to remove pods. |
| `excludeNamespaces` | String | `kube-system` | Namespaces were pods can't be removed. |
| `scheduler` | String | `5s` | Scheduler to delete a random pod. e.g `10s`, `2m`, `4h`. |
| `labelSelector` | String | `""` | Label selector to filter the list of pods. |
| `gracePeriod` | Int64 | `0` | Grace period to remove the pod. |

### Requirements to test the app

* Have container management tool installed.
* Have a kubernetes cluster created.
* Have Kubectl installed.
* Have Helm installed.

## Create Kubernetes cluster

[Kind](https://kind.sigs.k8s.io/) can be used to create a Kubernetes cluster if you don't have your own cluster.

```
kind version
```

<details>
  <summary>Output</summary>

  ```
  kind v0.15.0 go1.19 linux/amd64
  ```
</details>

Create a new cluster with `kind`.

```
kind create cluster
```

<details>
  <summary>Output</summary>

  ```
  Creating cluster "kind" ...
    ✓ Ensuring node image (kindest/node:v1.25.0) 🖼 
    ✓ Preparing nodes 📦  
    ✓ Writing configuration 📜 
    ✓ Starting control-plane 🕹️ 
    ✓ Installing CNI 🔌 
    ✓ Installing StorageClass 💾 
    Set kubectl context to "kind-kind"
    You can now use your cluster with:

    kubectl cluster-info --context kind-kind

    Have a nice day! 👋
  ```
</details>

Check cluster status.

```
kubectl get nodes
```

<details>
  <summary>Output</summary>

  ```
  NAME                 STATUS   ROLES           AGE   VERSION
  kind-control-plane   Ready    control-plane   62s   v1.25.0
  ```
</details>

## Deploy pod examples into workloads namespace

Create `workloads` namespace

```
kubectl create namespace workloads.
```

<details>
  <summary>Output</summary>

  ```
  namespace/workloads created
  ```
</details>

Deploy `example-1`.

```
kubectl --namespace workloads apply \
  --filename https://raw.githubusercontent.com/mmorejon/erase-una-vez-k8s/main/deployments/deploy-01.yaml
```

<details>
  <summary>Output</summary>

  ```
  deployment.apps/deploy-example-1 created
  ```
</details>

Deploy `example-2`.

```
kubectl --namespace workloads apply \
  --filename https://raw.githubusercontent.com/mmorejon/erase-una-vez-k8s/main/deployments/deploy-02.yaml
```

<details>
  <summary>Output</summary>

  ```
  deployment.apps/deploy-example-2 created
  ```
</details>

List all pods created in the workloads namespace.

```
kubectl --namespace workloads get pods
```

<details>
  <summary>Output</summary>

  ```
  NAME                                READY   STATUS    RESTARTS   AGE
  deploy-example-1-7bd69c4c97-7sts9   1/1     Running   0          4m15s
  deploy-example-1-7bd69c4c97-88k6v   1/1     Running   0          4m15s
  deploy-example-1-7bd69c4c97-br4xf   1/1     Running   0          4m15s
  deploy-example-1-7bd69c4c97-jht7s   1/1     Running   0          4m15s
  deploy-example-1-7bd69c4c97-jj88k   1/1     Running   0          4m15s
  deploy-example-1-7bd69c4c97-kb6g6   1/1     Running   0          4m15s
  deploy-example-1-7bd69c4c97-l25vf   1/1     Running   0          4m15s
  deploy-example-2-5d6ffd8d74-5zpj2   1/1     Running   0          96s
  deploy-example-2-5d6ffd8d74-82mwz   1/1     Running   0          96s
  deploy-example-2-5d6ffd8d74-8z58n   1/1     Running   0          96s
  deploy-example-2-5d6ffd8d74-9fr7p   1/1     Running   0          96s
  deploy-example-2-5d6ffd8d74-9ws6h   1/1     Running   0          96s
  deploy-example-2-5d6ffd8d74-d9x7x   1/1     Running   0          96s
  deploy-example-2-5d6ffd8d74-nbtrr   1/1     Running   0          96s
  deploy-example-2-5d6ffd8d74-pzzgw   1/1     Running   0          96s
  deploy-example-2-5d6ffd8d74-s4mv6   1/1     Running   0          96s
  deploy-example-2-5d6ffd8d74-vdbhs   1/1     Running   0          96s
  ```
</details>

## Deploy podchaosmonkey app

```
helm upgrade --install podchaosmonkey \
  --namespace podchaosmonkey --create-namespace \
  chart/podchaosmonkey
```

<details>
  <summary>Output</summary>

  ```
  Release "podchaosmonkey" does not exist. Installing it now.
  NAME: podchaosmonkey
  LAST DEPLOYED: Wed Sep 14 11:55:46 2022
  NAMESPACE: podchaosmonkey
  STATUS: deployed
  REVISION: 1
  TEST SUITE: None
  ```
</details>

Check podchaosmonkey pod deployment.

```
kubectl --namespace podchaosmonkey get pod
```

<details>
  <summary>Output</summary>

  ```
  NAME                             READY   STATUS    RESTARTS   AGE
  podchaosmonkey-9c9bc4586-l68rc   1/1     Running   0          8s
  ```
</details>

See podchaosmonkey logs.

```
kubectl --namespace podchaosmonkey logs \
  --selector app.kubernetes.io/name=podchaosmonkey --follow
```

<details>
  <summary>Output</summary>

  ```
  Starting chaos process ...
  Pods in the namespace workloads will be removed every 5s.

  Waiting for the next schedule.
  It is time to remove a new pod ...
  Number of pods available 17
  The pod deploy-example-1-7bd69c4c97-88k6v was removed.

  Waiting for the next schedule.
  It is time to remove a new pod ...
  Number of pods available 17
  The pod deploy-example-2-5d6ffd8d74-nbtrr was removed.

  Waiting for the next schedule.
  It is time to remove a new pod ...
  Number of pods available 17
  The pod deploy-example-1-7bd69c4c97-jj88k was removed.
  ```
</details>