# AutoScale K8s REST example

A REST api example developed in Go, based on Echo framework. backed by Redis, running on K8s, featured with Ingress and HPA.

## Get Started

```
git clone https://github.com/csimplestring/autoscale-k8s-rest-example.git
minikube start (you need to install minikube before)
make build-docker
```

## Deploy on Minikube

```
cd k8s
kubectl create -f namespaces.yaml
kubectl create -f volumes.yaml
kubectl create -f ingress.yaml
kubectl create -f ingress-controller.yaml
kubectl create -f redis.yaml
kubectl create -f api.yaml
```

Because we are using the local Minikube cluster, it is not possible to expose the Ingress Controller
as a 'LoadBalancer' Service(usually you should) on your cloud provider. 

In order to access the api via Ingress, you have to run: 
```
kubectl describe svc --namespace=kube-system traefik-ingress-service
```
you need to find the exposed port, for example. 31321, use this port as {api-port}

## Test

```
// create
Post http://{minikube-ip}:{api-port}}/api/customer
Host: api.minikube
Content-Type: application/json
{
  "name": "1",
  "address": "b"
}
```

```
// get
GET http://{minikube-ip}:{api-port}}/api/customer/1
Host: api.minikube
Content-Type: application/json
```

## Run CronJob

// this cron job will back up the redis data every 5 minutes
```
kubectl create -f jobs.yaml
```

## Enable HorizontalPodAutoscaler

before this, you need to enable heapster and metric-server in Minikube.
```
	minikube addons enable metrics-server
	minikube addons enable heapster
	kubectl autoscale deployment api --cpu-percent=10 --min=2 --max=5 --namespace=testing
```
which means when the average CPU usage is more than 10%, the api pod will be automatically
created, the maximum pod number will be 5. If the average CPU usage is less than 10% later,
the api pods will be automatically deleted in order to save resources.

