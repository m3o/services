---
title: Docker Deployment
keywords: docker
tags: [docker]
sidebar: home_sidebar
permalink: deploy-docker.html
summary: 
---

Micro easily runs inside docker containers

## Prebuilt Images

Prebuilt images are available on [Docker Hub](https://hub.docker.com/r/micro/)

### Install Micro

```
docker pull micro/micro
```

## Compose

Run a local deployment using docker compose

```
consul:
  command: -server -bootstrap -rejoin 
  image: progrium/consul:latest
  hostname: "registry"
  ports:
  - "8300:8300"
  - "8400:8400"
  - "8500:8500"
  - "8600:53/udp"
api:
  command: --registry_address=registry:8500 --register_interval=5 --register_ttl=10 api
  build: .
  links:
  - consul
  ports:
  - "8080:8080"
proxy:
  command: --registry_address=registry:8500 --register_interval=5 --register_ttl=10 proxy
  build: .
  links:
  - consul
  ports:
  - "8081:8081"
web:
  command: --registry_address=registry:8500 --register_interval=5 --register_ttl=10 web
  build: .
  links:
  - consul
  ports:
  - "8082:8082"
```

## Build from scratch

A Dockerfile is included in the repo

```
FROM alpine:3.2
RUN apk add --update ca-certificates && \
    rm -rf /var/cache/apk/* /tmp/*
ADD micro /micro
WORKDIR /
ENTRYPOINT [ "/micro" ]
```

### Build Binary

```
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -i -o micro ./main.go 
```

### Build Image

```
docker build -t micro .
```

{% include links.html %}
