# sample-deployment-scaler
カスタムリソースとカスタムコントローラのサンプル。学習用

## Description
DeploymentScaler（カスタムリソース）に定義したreplicasで、対象のDeploymentのreplicasを保ちます。

## 動かし方

### kindでクラスター作成

```
kind create cluster --name deploy-scaler
```

### カスタムコントローラのビルド

```
export IMG=deployment-scaler:dev-$(date +%s)
make docker-build IMG=$IMG
```

### kindにDocker imageをpush

```
kind load docker-image $IMG --name deploy-scaler
```

### カスタムコントローラーのデプロイ

```
make deploy IMG=$IMG
```

### deploymentとカスタムリソースの適用

次のコマンドを実行し、下記のようなログが出ていれば成功

```
❯ kubectl apply -k config/samples/
deployment.apps/server created
deploymentscaler.scaling.example.com/deploymentscaler-sample created
```

## 動作確認

1. Deploymentのreplicasを確認

```
> kubectl describe deployment server
...
Replicas:               4 desired | 4 updated | 4 total | 4 available | 0 unavailable
```

2. Deploymentのreplicasを手動で変更

```
❯ kubectl scale deployment server --replicas=2
deployment.apps/server scaled
```

3. Deploymentのreplicasが1のreplicasから変更ないことを確認（DeploymentScalerにより2で変更したが１の状態に調整された）

```
❯ kubectl describe deployment server
...
Replicas:               4 desired | 4 updated | 4 total | 4 available | 0 unavailable
```

### ログでReconciliation Loopを確認する

1. DeploymentScalerのPodを確認

```
❯ kubectl get pod -n sample-deployment-scaler-system
NAME                                                           READY   STATUS    RESTARTS   AGE
sample-deployment-scaler-controller-manager-67c5d4d885-7htnh   1/1     Running   0          66s
```

2. 1で取得したPodのログを確認. 以下のようなログが出ていればDeploymentScalerによりreplicasが調整されている

```
❯ kubectl -n sample-deployment-scaler-system logs sample-deployment-scaler-controller-manager-67c5d4d885-7htnh -c manager -f

INFO    Reconciling DeploymentScaler    {"controller": "deploymentscaler", "controllerGroup": "scaling.example.com", "controllerKind": "DeploymentScaler", "DeploymentScaler": {"name":"deploymentscaler-sample","namespace":"default"}, "namespace": "default", "name": "deploymentscaler-sample", "reconcileID": "46841168-ab0b-488d-9a4e-319e01b99b4f", "request": {"name":"deploymentscaler-sample","namespace":"default"}, "namespace": "default", "name": "deploymentscaler-sample"}
```
