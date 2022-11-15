# finalizers

Small piece of go code to remove finalizers from pods

https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers/

https://kubebyexample.com/learning-paths/operator-framework/kubernetes-api-fundamentals/finalizers

kubectl run nginx --image=nginx --dry-run=client -o yaml

```bash

# start a cluster
minikube start --cpus=3 --memory=6144 --kubernetes-version=v1.23.2

# deploy a pod with a finalizer
k apply -f ./pod.yaml

# delete the pod and see that it stays in a terminating state
k delete -f ./pod.yaml
k get po -w 

# run the app to remove the finalizer
go run main.go

```

