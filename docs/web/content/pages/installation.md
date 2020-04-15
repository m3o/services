---
title: Install Guide
keywords: install
tags: [install]
sidebar: home_sidebar
permalink: installation.html
summary: 
---

## Framework

Go Micro is an RPC framework for development microservices in Go

### Dependencies

You will need protoc-gen-micro for code generation

- [protoc-gen-micro](https://github.com/micro/protoc-gen-micro)

### Import

Ensure you import go-micro v2

```
import "github.com/micro/go-micro/v2"
```

## Runtime

Micro provides a runtime for accessing and managing microservices

### Install

From source

```
go get github.com/micro/micro/v2
```

Docker image

```
docker pull micro/micro
```

Latest release binaries

```
# MacOS
curl -fsSL https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh | /bin/bash

# Linux
wget -q  https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh -O - | /bin/bash

# Windows
powershell -Command "iwr -useb https://raw.githubusercontent.com/micro/micro/master/scripts/install.ps1 | iex"
```

### Usage

Start the server

```shell
micro server
```

Run the greeter service

```shell
micro run github.com/micro/examples/greeter/srv
```

List services

```shell
micro list services
```

Get Service

```shell
micro get service go.micro.srv.greeter
```

Output

```shell
service  go.micro.srv.greeter

version 2019.11.09.10.34

ID      Address Metadata
go.micro.srv.greeter-e25a5edd-0936-4d32-b4d7-e62bf454d5f7       172.17.0.1:33031        broker=http,protocol=mucp,registry=mdns,server=mucp,transport=http

Endpoint: Say.Hello

Request: {
        name string
}

Response: {
        msg string
}
```

Call service

```shell
micro call go.micro.srv.greeter Say.Hello '{"name": "John"}'
```

Output

```shell
{
	"msg": "Hello John"
}
```

{% include links.html %}
