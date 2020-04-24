---
title: Getting Started with Micro
keywords: micro
tags: [micro]
sidebar: home_sidebar
permalink: /getting-started.html
summary: A getting started guide for Micro
---

# Getting Started with Micro

> 

## What is Micro?

Micro is a system for building and managing microservices.

Key components of Micro are as follows:
* **Runtime**: a runtime environment for managing services including auth, config, discovery, networking, storage. 
* **Framework**: a Go framework for writing services to be run in the runtime.
* **Clients**: multi-language clients to enable other programs to access Micro services.

## Install

Using Go:

```sh
go install github.com/micro/micro/v2
```

Or by downloading the binary

```sh
# MacOS
curl -fsSL https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh | /bin/bash

# Linux
wget -q  https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh -O - | /bin/bash

# Windows
powershell -Command "iwr -useb https://raw.githubusercontent.com/micro/micro/master/scripts/install.ps1 | iex"
```

## Running a service

Before diving into writing a service, let's run an existing one, because it's just a few commands away!


First, we have to start the `micro server`. The command to do that is:

```sh
micro server
```

To talk to this server, we just have to tell Micro CLI to address our server instead of using the default implementations - micro can work without a server too, but let's ignore that for now.

The following command tells the CLI to talk to our server:

```
micro env set server
```

Great! We are ready to roll. Just to verify that everything is in order, let's see what services are running:

```
$ micro list services
go.micro.api
go.micro.auth
go.micro.bot
go.micro.broker
go.micro.config
go.micro.debug
go.micro.network
go.micro.proxy
go.micro.registry
go.micro.router
go.micro.runtime
go.micro.server
go.micro.tunnel
go.micro.web
```

All those services are ones started by our `micro server`. This is pretty cool, but still it's not something we launched! Let's start a service for which existence we can actually take credit for.

If we go to [github.com/micro/services](https://github.com/micro/services), we see a bunch of services written by micro authors. One of them is the `helloworld`. Try our luck, shall we?

The command to run services is `micro run`. This command may take a while as it checks out
the repository from GitHub. (@todo this actually fails currently, fix)

```
micro run github.com/micro/services/helloworld
```


If we take a look at the running `micro server`, we should see something like

```
Creating service helloworld version latest source /tmp/github.com-micro-services/helloworld
Processing create event helloworld:latest
```

We can also have a look at logs of the service to verify it's running.
The default log location is the temporary directory of your system:

```
tail /tmp/micro/logs/helloworld.log
```

```
Starting [service] go.micro.service.helloworld
Server [grpc] Listening on [::]:41601
```

So since our service is running happily, let's try to call it! That's what services are for.

## Calling a service

We have a couple of options to call a service running on our `micro server`.

### From CLI

The easiest is perhaps with the CLI:

```sh
$ micro call go.micro.service.helloworld Helloworld.Call '{"name":"Jane"}'
{
	"msg": "Hello Jane"
}

```

That worked! If we wonder what endpoints a service has, the best place to look for is its [proto folder](https://github.com/micro/services/blob/master/helloworld/proto/helloworld/helloworld.proto). There are other tools in the making too, like [explore](https://web.micro.mu/explore/?go#helloworld), but they are still a bit experimental.

### With Go Micro

Let's write the most minimal service we can have that calls an other service.
The program below should output `Response:  Hello John`.

```go
package main

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/v2"
	proto "github.com/micro/services/helloworld/proto"
)

func main() {
	service := micro.NewService()
	service.Init()

	client := proto.NewHelloworldService("go.micro.service.helloworld", service.Client())

	rsp, err := client.Call(context.Background(), &proto.Request{
		Name: "John",
	})
	if err != nil {
		fmt.Println("call err: ", err)
		return
	}

	fmt.Println("Response: ", rsp.Msg)

}
```

### From other languages

In the [clients repo](https://github.com/micro/clients) there are Micro clients for various languages and frameworks. They are designed to connect easily to the live Micro environment or your local one, but more about environments later.

## Writing a service

To scaffold a new service, the `micro new` command can be used. It should output something
reasonably similar to the following:

```sh
$ micro new foobar
Creating service go.micro.service.foobar in foobar

.
├── main.go
├── generate.go
├── plugin.go
├── handler
│   └── foobar.go
├── subscriber
│   └── foobar.go
├── proto/foobar
│   └── foobar.proto
├── Dockerfile
├── Makefile
├── README.md
├── .gitignore
└── go.mod


download protobuf for micro:

brew install protobuf
go get -u github.com/golang/protobuf/proto
go get -u github.com/golang/protobuf/protoc-gen-go
go get github.com/micro/micro/v2/cmd/protoc-gen-micro@master

compile the proto file foobar.proto:

cd foobar
protoc --proto_path=.:$GOPATH/src --go_out=. --micro_out=. proto/foobar/foobar.proto
```

As can be seen from the output above, before building the first service, the following tools must be installed:
* [protoc](http://google.github.io/proto-lens/installing-protoc.html)
* [protobuf/proto](github.com/golang/protobuf/protoc-gen-go)
* [protoc-gen-micro](github.com/golang/protobuf/protoc-gen-go)

They are all needed to translate proto files to actual Go code.
Protos exist to provide a language agnostic way to describe service endpoints, their input and output types, and to have an efficient serialization format at hand.

Currently Micro is  Go focused (apart from the [clients](#-from-other-languages) mentioned before), but this will change soon.

So once all tools are installed, being inside the service root, we can issue the following command to generate the Go code from the protos:

```
protoc --proto_path=.:$GOPATH/src --go_out=. --micro_out=. proto/foobar/foobar.proto
```

The generated code must be committed to source control, to enable other services to import the proto when making service calls (see previous section [Calling a service](#-calling-a-service).

At this point, we know how to write a service, run it, and call other services too.
We have everything at our fingertips, but there are still some missing pieces to write applications. One of such pieces is the store interface, which helps with persistent data storage even without a database.

## Storage

Micro includes a persistent store service for storing key-value data.

### Go Micro interfaces in general

Both Micro (the server/CLI) and Go Micro (the framework) is very centered around precisely defined interfaces and their different implementations. What does this mean?

Let's take our current case of the [store interface](https://github.com/micro/go-micro/blob/master/store/store.go). It's aimed to enable service writers data storage with a couple of different implementations:
* in memory
* local file (default when running `micro server`)
* couchbase backed

Similarly, the [runtime](https://github.com/micro/go-micro/blob/master/runtime/runtime.go) interface, that represents something that runs processes, has a couple of implementations:
* local, which just runs actual binaries - aimed at local usage (who would have guessed)
* kubernetes, aimed for beefier production settings
