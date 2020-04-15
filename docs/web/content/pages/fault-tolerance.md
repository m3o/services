---
title: Fault Tolerance
keywords: fault tolerance
tags: [fault tolerance]
sidebar: home_sidebar
permalink: fault-tolerance.html
summary: 
---

In a distributed systems world failures are occurring all the time. Micro attempts to address this with some best practices for fault tolerance. 
This document describes some of the ways in which they can be configured.

## Heartbeating

Heartbeating is a mechanism of refreshing registration in service discovery.

### Rationale

Services register with service discovery on startup and deregister on shutdown. Sometimes these services may unexpectedly die, 
be killed forcefully or face transient network issues. In these cases stale nodes will be left in service discovery. It would be 
ideal if services were automatically removed.

### Solution

Micro supports the option of a register TTL and register interval for this exact reason. TTL specifies how long a registration should 
exist in discovery after which it expires and is removed. Interval is the time at which a service should re-register to preserve 
it's registration in service discovery.

These are options made available in go-micro and as flags in the micro toolkit

### Usage

For the micro toolkit simply use the built in flags to set ttl and interval.

```
micro --register_ttl=30 --register_interval=15 api
```

This example shows that we're setting a ttl of 30 seconds with a re-register interval of 15 seconds.

For go-micro, when declaring a micro service you can pass in the options as time.Duration

```
service := micro.NewService(
        micro.Name("com.example.srv.foo"),
        micro.RegisterTTL(time.Second*30),
        micro.RegisterInterval(time.Second*15),
)
```

## Load Balancing

Load balancing is a way of spreading request load or maintaining high availability

### Rationale

There are limitations to availability and scaling with any single process application. If an application dies for any reason 
you are no longer able to serve requests. If enough request load is sent to the application it may begin to respond slowly 
or not at all. Sending requests across multiple copies of an application would solve these problems.

### Solution

Micro does client side load balancing via the [selector](https://pkg.go.dev/github.com/micro/go-micro/v2/selector#Selector) interface 
to spread requests across any number of nodes of a service. When a service is started it registers with service discovery as a 
service node with a unique address and id. When making a request the micro client uses the selector to decide which node to make 
the request to. The selector uses the service registry to find the nodes of a service, then use a load balancing strategy such as 
random hash or round robin to select a node to send the request to.

### Usage

Client side load balancing is built into the go-micro client. This is done automatically.

## Retries

Retries are a method of retrying a request if it was unsuccessful

### Rationale

A request may fail for any number of reasons; network error, request load, application death. In these instances it would be ideal 
if we could retry the request against a different copy of the application to receive a successful response.

### Solution

The micro client includes a mechanism for retrying requests. The selector (mentioned above) returns a Next function which when executed 
uses a load balancing strategy to return a node from the list. The Next function can be executed multiple times, returning a new node 
based on the load balancing strategy. With retries set, if a request fails, the Next function will be executed and the request 
will be retried against a new node.

### Usage

Retries can be set as a flag or option to the client. It defaults to 1, meaning 1 attempt at a request.

To change via flag

```
micro --client_retries=3
```

To set as an option

```
client.Init(
	client.Retries(3),
)
```

## Caching Discovery

A discovery cache is a client side cache of service discovery information

### Rationale

Service discovery is the core dependency of microservices but can also become a single point of failure if not architected correctly. Every 
discovery system has its own properties for scaling and high availability. In the event discovery goes down it makes the rest of the system 
unusable as services do not know how to resolve names to addresses. Discovery can also be a bottleneck if a lookup is done for every 
request in the system.

### Solution

Client side caching is a method of eliminating service discovery as a bottleneck and single point of failure. Micro includes a selector 
(client side load balancer) which maintains an in memory cache of service discovery information relevant to it. If a cache miss occurs 
the selector will use the service registry to do a lookup and cache the data. The cache is also periodically TTLed out to ensure 
stale data does not live on.

### Usage

The caching selector can be set by flag or when creating a new service.

As a flag with the toolkit do the following

```
micro --selector=cache api
```

Go-micro services can also use the same flag if the Init method is called. Alternatively the selector can be set in code.

```
import (
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/selector/cache"
)

service := micro.NewService(
	micro.Name("com.example.srv.foo"),
)

service.Client().Init(cache.NewSelector())
```

{% include links.html %}
