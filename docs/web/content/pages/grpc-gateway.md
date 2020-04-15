---
title: GRPC Gateway
keywords: grpc, grpc-gateway
tags: [grpc, grpc-gateway]
sidebar: home_sidebar
permalink: grpc-gateway.html
summary: 
---

This guide helps with using the grpc gateway with go-micro services.

The [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) is a plugin of [protoc](http://github.com/google/protobuf).
It reads [gRPC](http://github.com/grpc/grpc-common) service definition, and generates a reverse-proxy server which translates 
a RESTful JSON API into gRPC.

We use [go-grpc](https://github.com/micro/go-grpc) to write the backend service. Go-GRPC is a simple wrapper around go-micro and 
the grpc plugins for the client and server. When calling [grpc.NewService](https://godoc.org/github.com/micro/go-grpc#NewService) 
it returns a [micro.Service](https://pkg.go.dev/github.com/micro/go-micro/v2#Service).

## Code

Find the example code at [examples/grpc](https://github.com/micro/examples/tree/master/grpc).

## Prereqs

These are some prerequisites

### Install protobuf

```
mkdir tmp
cd tmp
git clone https://github.com/google/protobuf
cd protobuf
./autogen.sh
./configure
make
make check
sudo make install
```

### Install plugins

```
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/micro/protobuf/protoc-gen-go
```

## Greeter Service

In the example we've created a greeter microservice using the gRPC service. The service is very simple. 

The proto is as follows:

```proto
syntax = "proto3";

package go.micro.srv.greeter;

service Say {
	rpc Hello(Request) returns (Response) {}
}

message Request {
	string name = 1;
}

message Response {
	string msg = 1;
}
```

The service is as follows:

```
package main

import (
	"log"
	"time"

	hello "github.com/micro/examples/greeter/srv/proto/hello"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/service/grpc"

	"golang.org/x/net/context"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {
	service := grpc.NewService(
		micro.Name("go.micro.srv.greeter"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	hello.RegisterSayHandler(service.Server(), new(Say))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## GRPC Gateway

The grpc gateway uses the same proto as the service with the addition of a http option

```proto
syntax = "proto3";

package greeter;

import "google/api/annotations.proto";

service Say {
	rpc Hello(Request) returns (Response) {
		option (google.api.http) = {
			post: "/greeter/hello"
			body: "*"
		};
	}
}

message Request {
	string name = 1;
}

message Response {
	string msg = 1;
}
```

The proto generates a grpc stub and reverse proxy with the following commands

```
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:. \
  path/to/your_service.proto
```

```
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:. \
  path/to/your_service.proto
```

We've created a sample api for the greeter service with the following code. Similar code would be written 
to register other endpoints. Note the gateway requires the endpoint address of the greeter service.

```

package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	hello "github.com/micro/examples/grpc/gateway/proto/hello"
)

var (
	// the go.micro.srv.greeter address
	endpoint = flag.String("endpoint", "localhost:9090", "go.micro.srv.greeter address")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := hello.RegisterSayHandlerFromEndpoint(ctx, mux, *endpoint, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":8080", mux)
}

func main() {
	flag.Parse()

	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
```

## Running the example

Run the greeter service. Specify mdns since we don't need discovery

```
go run examples/grpc/greeter/srv/main.go --registry=mdns --server_address=localhost:9090
```

Run the gateway. It will default to the endpoint localhost:9090 for the greeter service

```
go run examples/grpc/gateway/main.go
```

Curl your request at the gateway (localhost:8080)

```
curl -d '{"name": "john"}' http://localhost:8080/greeter/hello
```

## Limitations

The example grpc gateway requires a service address to be provided whereas our own micro api uses service discovery, dynamic routing 
and load balancing. This makes integration of the grpc gateway a little less versatile.

Visit [github.com/micro/micro](https://github.com/micro/micro) to learn more

{% include links.html %}
