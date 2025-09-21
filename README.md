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
export IMG=deployment-scaler:dev
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

1. 起動しているPodを確認

```
❯ kubectl get po
NAME                     READY   STATUS    RESTARTS   AGE
server-5b997465b-6z42v   1/1     Running   0          11m
server-5b997465b-fk67z   1/1     Running   0          11m
server-5b997465b-wsrhx   1/1     Running   0          11m
```

2. Podを１つ削除する

```
❯ kubectl delete po server-5b997465b-6z42v
pod "server-5b997465b-6z42v" deleted
```

3. 再度起動しているPodを確認

```
❯ kubectl get po
NAME                     READY   STATUS    RESTARTS   AGE
server-5b997465b-fk67z   1/1     Running   0          11m
server-5b997465b-pm9sf   1/1     Running   0          2s
server-5b997465b-wsrhx   1/1     Running   0          11m
```
