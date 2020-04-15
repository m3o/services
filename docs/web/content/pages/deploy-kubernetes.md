---
title: Kubernetes Deployment
keywords: kubernetes
tags: [kubernetes]
sidebar: home_sidebar
permalink: deploy-kubernetes.html
summary: 
---


This doc provides a guide to running micro on kubernetes.

## Getting Started

- [Dependencies](#dependencies)
- [Install Micro](#install-micro)
- [Writing a Service](#writing-a-service)
- [Deploying a Service](#deploying-a-service)
- [Healthchecking Sidecar](#healthchecking-sidecar)
  - [Install healthchecker](#install-healthchecker)
  - [Run healtchecker](#run-healtchecker)
  - [K8s Deployment](#k8s-deployment)
- [Micro API](#micro-api)
- [Micro Web](#micro-web)
- [Full Deployment](#full-deployment)

## Dependencies

On kubernetes we would recommend running [etcd](https://github.com/etcd-io/etcd) and [nats](https://github.com/nats-io/nats-server).

- Etcd is used for highly scalable service discovery
- NATS is used for asynchronous messaging

To install etcd ([instructions](https://github.com/helm/charts/tree/master/stable/etcd-operator))

```
helm install --name my-release --set customResources.createEtcdClusterCRD=true stable/etcd-operator

# to later uninstall etcd
helm delete my-release
```

To install nats ([instructions](https://github.com/nats-io/nats-operator))

```
kubectl apply -f https://github.com/nats-io/nats-operator/releases/latest/download/00-prereqs.yaml
kubectl apply -f https://github.com/nats-io/nats-operator/releases/latest/download/10-deployment.yaml
```

You should now have the required dependencies.

## Install Micro


Use the below docker image

```
micro/micro
```

Or use go-micro for development

```
github.com/micro/go-micro
```

## Writing a Service

Write a service as you would any other [go-micro](https://github.com/micro/go-micro) service.

```go
import (
	"github.com/micro/go-micro/v2"
)

func main() {
	service := micro.NewService(
		micro.Name("greeter")
	)
	service.Init()
	service.Run()
}
```

## Deploying a Service

Here's an example k8s deployment for a micro service

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  namespace: default
  name: greeter
spec:
  replicas: 1
  selector:
    matchLabels:
      name: greeter-srv
      micro: service
  template:
    metadata:
      labels:
        name: greeter-srv
        micro: service
    spec:
      containers:
        - name: greeter
          command: [
		"/greeter-srv",
	  ]
          image: yourimage/yourservice
          imagePullPolicy: Always
          ports:
          - containerPort: 8080
            name: greeter-port
          env:
          - name: MICRO_SERVER_ADDRESS
            value: "0.0.0.0:8080"
          - name: MICRO_BROKER
            value: "nats"
          - name: MICRO_BROKER_ADDRESS
            value: "nats-cluster"
          - name: MICRO_REGISTRY
            value: "etcd"
          - name: MICRO_REGISTRY_ADDRESS
            value: "etcd-cluster-client"
```

Deploy with kubectl

```
kubectl apply -f greeter.yaml
```

## Healthchecking Sidecar

The healthchecking sidecar exposes `/health` as a http endpoint and calls the rpc endpoint `Debug.Health` on a service. 
Every go-micro service has a built in Debug.Health endpoint.

### Install healthchecker

```
go get github.com/micro/cmd/health
```

or

```
docker pull microhq/health
```

### Run healtchecker

Run e.g healthcheck greeter service with address localhost:9091

```
health --server_name=greeter --server_address=localhost:9091
```

Call the healthchecker on localhost:8080

```
curl http://localhost:8080/health
```

### K8s Deployment

Add to a kubernetes deployment

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  namespace: default
  name: greeter
spec:
  replicas: 1
  selector:
    matchLabels:
      name: greeter-srv
      micro: service
  template:
    metadata:
      labels:
        name: greeter-srv
        micro: service
    spec:
      containers:
        - name: greeter
          command: [
		"/greeter-srv",
	  ]
          image: yourimage/yourservice
          imagePullPolicy: Always
          ports:
          - containerPort: 8080
            name: greeter-port
          env:
          - name: MICRO_SERVER_ADDRESS
            value: "0.0.0.0:8080"
          - name: MICRO_BROKER
            value: "nats"
          - name: MICRO_BROKER_ADDRESS
            value: "nats-cluster"
          - name: MICRO_REGISTRY
            value: "etcd"
          - name: MICRO_REGISTRY_ADDRESS
            value: "etcd-cluster-client"
        - name: health
          command: [
		"/health",
                "--health_address=0.0.0.0:8081",
		"--server_name=greeter",
		"--server_address=0.0.0.0:8080"
	  ]
          image: microhq/health
          livenessProbe:
            httpGet:
              path: /health
              port: 8081
            initialDelaySeconds: 3
            periodSeconds: 3
```

## Micro API

To deploy the micro api use the following config. Note the ENABLE_ACME env var where you want Let's Encrypt SSL by default.

Create the api service

```
apiVersion: v1
kind: Service
metadata:
  name: micro-api
  namespace: default
  labels:
    name: micro-api
    micro: service
spec:
  ports:
  - name: https
    port: 443
    targetPort: 443
  selector:
    name: micro-api
    micro: service
  type: LoadBalancer
```

Create the deployment

```
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: default
  name: micro-api
  labels:
    micro: service
spec:
  replicas: 3
  selector:
    matchLabels:
      name: micro-api
      micro: service
  strategy:
    rollingUpdate:
      maxSurge: 0
      maxUnavailable: 1
  template:
    metadata:
      labels:
        name: micro-api
        micro: service
    spec:
      containers:
      - name: api
        env:
        - name: MICRO_ENABLE_STATS
          value: "true"
        - name: MICRO_BROKER
          value: "nats"
        - name: MICRO_BROKER_ADDRESS
          value: "nats-cluster"
        - name: MICRO_REGISTRY
          value: "etcd"
        - name: MICRO_REGISTRY_ADDRESS
          value: "etcd-cluster-client"
        - name: MICRO_REGISTER_TTL
          value: "60"
        - name: MICRO_REGISTER_INTERVAL
          value: "30"
        - name: MICRO_ENABLE_ACME
          value: "true"
        args:
        - api
        image: micro/micro
        imagePullPolicy: Always
        ports:
        - containerPort: 443
          name: api-port
```

## Micro Web

To deploy the micro web use the following config. Note the ENABLE_ACME env var where you want Let's Encrypt SSL by default.

Create the service

```
apiVersion: v1
kind: Service
metadata:
  name: micro-web
  namespace: default
  labels:
    name: micro-web
    micro: service
spec:
  ports:
  - name: https
    port: 443
    targetPort: 443
  selector:
    name: micro-web
    micro: service
  type: LoadBalancer
```

Create the deployment

```
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: default
  name: micro-web
  labels:
    micro: service
spec:
  replicas: 1
  selector:
    matchLabels:
      name: micro-web
      micro: service
  template:
    metadata:
      labels:
        name: micro-web
        micro: service
    spec:
      containers:
      - name: web
        env:
        - name: MICRO_BROKER
          value: "nats"
        - name: MICRO_BROKER_ADDRESS
          value: "nats-cluster"
        - name: MICRO_ENABLE_STATS
          value: "true"
        - name: MICRO_REGISTRY
          value: "etcd"
        - name: MICRO_REGISTRY_ADDRESS
          value: "etcd-cluster-client"
        - name: MICRO_ENABLE_ACME
          value: "true"
        - name: MICRO_ACME_PROVIDER
          value: certmagic
        args:
        - web
        image: micro/micro
        imagePullPolicy: Always
        ports:
        - containerPort: 443
          name: web-port
```

## Full Deployment

For the full deployment config used by [micro.mu](https://micro.mu) for the cloud platform see:

[https://github.com/micro/micro/tree/master/network/config/kubernetes](https://github.com/micro/micro/tree/master/network/config/kubernetes)

Some of this config is specific to the needs of Micro itself including the use of CloudFlare.

{% include links.html %}
