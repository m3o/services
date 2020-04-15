---
title: Hello World
keywords: go-micro, framework, helloworld
tags: [go-micro, framework, helloworld]
sidebar: home_sidebar
permalink: "/go-helloworld.html"
summary: A hello world example using Go Micro
---

# Overview

This is a helloworld example of using Go Micro. We'll run through writing the whole thing.

## Writing a service

Writing a service is incredibly simple with Go Micro. It provides you the framework to move fast without needing to understand 
everything at first. Below is a simple greeter service example that we'll run through.

Find the full example code at [examples/service](https://github.com/micro/examples/tree/master/service).

### Service Proto

One of the key requirements of microservices is strongly defined interfaces. Micro uses protobuf to achieve this.

Here we define the Greeter handler with the method Hello. It takes a Request and Response both with one string arguments.

```proto
syntax = "proto3";

service Greeter {
	rpc Hello(Request) returns (Response) {}
}

message Request {
	string name = 1;
}

message Response {
	string greeting = 2;
}
```

### Generate Proto

After writing the proto definition we must compile it using protoc with the micro plugin.

```shell
protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. path/to/greeter.proto
```

### Implement Service

Now that we've defined the service interface we need to implement the service.

Below is the code for the greeter service. It does the following:

1. Implements the interface defined for the Greeter handler
2. Initialises a [micro.Service](https://pkg.go.dev/github.com/micro/go-micro/v2#Service)
3. Registers the Greeter handler
4. Runs the service

```go
package main

import (
	"context"
	"fmt"

	micro "github.com/micro/go-micro/v2"
	proto "github.com/micro/examples/service/proto"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}

func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("greeter"),
	)

	// Init will parse the command line flags.
	service.Init()

	// Register handler
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
```

### Run Service

Now run the example using Go. If you're working with the example code do the following:

```
go run examples/service/main.go
```

This should output something like the folloing

```
2019-11-13 21:39:38.327452 I | Transport [http] Listening on [::]:32945
2019-11-13 21:39:38.327548 I | Broker [http] Connected to [::]:38955
2019-11-13 21:39:38.328095 I | Registry [mdns] Registering node: greeter-39373107-5ae7-42a2-b5e2-ebebf7eafbd9
```

### Write Client

Once we've got a service we actually need a way to query it. This is at the heart of microservices as we not 
only serve but also consume other services. Below is the client code to query the greeter service. 

The generated proto includes a greeter client to reduce boilerplate code.

```go
package main

import (
	"context"
	"fmt"

	micro "github.com/micro/go-micro/v2"
	proto "github.com/micro/examples/service/proto"
)

func main() {
	// Create a new service
	service := micro.NewService(micro.Name("greeter.client"))
	// Initialise the client and parse command line flags
	service.Init()

	// Create new greeter client
	greeter := proto.NewGreeterService("greeter", service.Client())

	// Call the greeter
	rsp, err := greeter.Hello(context.TODO(), &proto.Request{Name: "John"})
	if err != nil {
		fmt.Println(err)
	}

	// Print response
	fmt.Println(rsp.Greeting)
}
```

Now run the client

```shell
go run client.go
```

And the output should simply print the response

```
Hello John
```

{% include links.html %}
