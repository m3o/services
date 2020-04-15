---
title: Internals
keywords: go-micro
tags: [go-micro]
sidebar: home_sidebar
permalink: /go-internals.html
summary: A description of the go-micro internal features
---

Go Micro is a framework for microservices development. It provides pluggable abstractions via the use of Go interfaces.

<p align="center">
  <img src="images/go-micro.svg" />
</p>

## Features

Go Micro abstracts away the details of distributed systems. Here are the main features.

- **Service Discovery** - Automatic registration and name resolution with service discovery
- **Load Balancing** - Smart client side load balancing of services built on discovery
- **Synchronous Comms** - RPC based communication with support for bidirectional streaming
- **Asynchronous Comms** - PubSub interface built in for event driven architectures
- **Message Encoding** - Dynamic encoding based on content-type with protobuf and json out of the box
- **Service Interface** - All features are packaged in a simple high level interface for developing microservices

Go Micro supports a services programming model for building distributed systems.

## Plugins

Go Micro makes use of the Go interface for it's abstractions. Because of this the underlying implementation can be swapped out.

We provide sane defaults out of the box:

- mDNS for zeroconf service discovery
- Random hashed client side load balancing
- json and protobuf encoding based on Content-Type
- RPC over http synchronous communication
- HTTP point-to-point pubsub for asynchronous communication

## Packages

Go micro is composed of a number of packages:

- **transport** for sync messaging
- **broker** for async messaging
- **codec** for message encoding
- **registry** for service discovery
- **selector** for load balancing
- **client** for making requests
- **server** for handling requests
- **store** for key-value storage
- **config** for dynamic configuration

Further details are below

### Registry

The registry provides a service discovery mechanism to resolve names to addresses. It can be backed by consul, etcd, zookeeper, dns, gossip, etc. 
Services should register using the registry on startup and deregister on shutdown. Services can optionally provide an expiry TTL and reregister 
on an interval to ensure liveness and that the service is cleaned up if it dies.

### Selector

The selector is a load balancing abstraction which builds on the registry. It allows services to be "filtered" using filter functions and "selected" 
using a choice of algorithms such as random, roundrobin, leastconn, etc. The selector is leveraged by the Client when making requests. The client 
will use the selector rather than the registry as it provides that built in mechanism of load balancing. 

### Transport

The transport is the interface for synchronous request/response communication between services. It's akin to the golang net package but provides 
a higher level abstraction which allows us to switch out communication mechanisms e.g http, rabbitmq, websockets, NATS. The transport also 
supports bidirectional streaming. This is powerful for client side push to the server.

### Broker

The broker provides an interface to a message broker for asynchronous pub/sub communication. This is one of the fundamental requirements of an event 
driven architecture and microservices. By default we use an inbox style point to point HTTP system to minimise the number of dependencies required 
to get started. However there are many message broker implementations available in go-plugins e.g RabbitMQ, NATS, NSQ, Google Cloud Pub Sub.

### Codec

The codec is used for encoding and decoding messages before transporting them across the wire. This could be json, protobuf, bson, msgpack, etc. 
Where this differs from most other codecs is that we actually support the RPC format here as well. So we have JSON-RPC, PROTO-RPC, BSON-RPC, etc. 
It separates encoding from the client/server and provides a powerful method for integrating other systems such as gRPC, Vanadium, etc.

### Server

The server is the building block for writing a service. Here you can name your service, register request handlers, add middeware, etc. The service 
builds on the above packages to provide a unified interface for serving requests. The built in server is an RPC system. In the future there maybe 
other implementations. The server also allows you to define multiple codecs to serve different encoded messages.

### Client

The client provides an interface to make requests to services. Again like the server, it builds on the other packages to provide a unified interface 
for finding services by name using the registry, load balancing using the selector, making synchronous requests with the transport and asynchronous 
messaging using the broker. 

The  above components are combined at the top-level of micro as a **Service**.

We additionally provide some other components for distributed systems development:

### Store

The store is a simple key-value storage interface used to abstract away lightweight data storage. We're not trying to implement a full blown sql dialect 
or storage, just simply the ability to hold state that would otherwise be lost in stateless services. They store interface will become a building block 
for all forms of storage in the future.

### Config

Config is an interface for dynamic config loading from any number of sources which can be combined and merged. Most systems actively require configuration 
that changes independent of the code. Having a config interface which can dynamically load these values as needed is powerful. It supports 
many different configuration formats also.


{% include links.html %}
