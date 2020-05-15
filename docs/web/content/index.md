---
title: Docs
keywords: homepage
tags: [introduction]
sidebar: home_sidebar
permalink: index.html
---

Micro addresses the key requirements for building distributed systems in the Cloud and beyond. It leverages the microservices
architecture pattern and provides a set of services which act as the building blocks of a platform. Micro deals
with the complexity of distributed systems and provides simpler programmable abstractions to build on.

<img src="https://micro.mu/images/runtime10.svg?load" />

Technology is constantly evolving. The infrastructure stack is always changing. Micro is a platform which
addresses these issues with a pluggable foundation and strongly defined apis to build on. Plug into any stack or cloud.

Follow us on [Twitter](https://twitter.com/microhq) or join the [Slack](https://slack.micro.mu) community.

## Features

The runtime is composed of the following features:


### Clients

Clients are entrypoints into the system. They enable access to your services through known access points.

- **api:** An api gateway which acts as a single entry point for the frontend with dynamic request routing using service discovery. 

- **bot:** A slack bot which enables you to query and interact with Micro directly from within slack. It's great for ChatOps.

- **cli:** Access services via the terminal. Every good developer tool needs a CLI as a defacto standard for operating a system. 

- **proxy:** An identity away proxy built which allows you to access remote environments without painful configuration or vpn.

- **web:** The web dashboard allows you to explore your services, describe their endpoints, the request and response formats and even
query them directly.

### Services

Services are the core services that makeup the runtime. They provide a programmable abstraction layer for distributed systems infrastructure.

- **auth:** Authentication and authorization is a core requirement for any production ready platform. Micro builds in an auth service 
for managing service to service and user to service authentication.

- **broker:** A message broker allowing for async messaging. Microservices are event driven architectures and should provide messaging as a first
class citizen. Notify other services of events without needing to worry about a response.

- **config:** Manage dynamic config in a centralised location for your services to access. Has the ability to load config from multiple 
sources and enables you to update config without needing to restart services.

- **debug:** Built in aggregation of stats, logs and tracing info for debugging. The debug service scrapes all services for their info to 
help understand the overall scope of the system from one location. 

- **network:** A drop in service to service networking solution. Offload service discovery, load balancing and fault tolerance to the network.
The micro network dynamically builds a latency based routing table based on the local registry. It includes support for multi-cloud networking.

- **registry:** The registry provides service discovery to locate other services, store feature rich metadata and endpoint information. It's a
service explorer which lets you centrally and dynamically store this info at runtime.

- **runtime:** A service runtime which manages the lifecycle of your service, from source to running. The runtime service can natively locally 
or on kubernetes, providing a seamless abstraction across both.

- **store:** State is a fundamental requirement of any system. We provide a key-value store to provide simple storage of state which can be shared
between services or offload long term to keep microservices stateless and horizontally scalable.

### Framework

To write applications which run on Micro you can use the framework Go Micro.

- **go-micro:** Leverage the powerful [Go Micro](https://github.com/micro/go-micro) framework to develop microservices easily and quickly.
Go Micro abstracts away the complexity of distributed systems and provides simpler abstractions to build highly scalable microservices.

## Getting Started

- **Framework** - Start writing services using the [go-micro](framework.html) framework.
- **Runtime** - Manage your services using the [micro](runtime.html) runtime.
- **Platform** - Host and collaborate on your services on the [platform](platform.html).
- **Examples** - Learn by example with real code [examples](https://github.com/micro/examples)
- **Community** - Join the the community and start collaborating on [slack](https://slack.micro.mu)

{% include links.html %}
