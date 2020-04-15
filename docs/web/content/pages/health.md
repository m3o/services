---
title: Service Health
keywords: health
tags: [health]
sidebar: home_sidebar
permalink: health.html
summary: Micro health is a healthchecking sidecar
---

The service healthchecker is a http sidecar which runs at http://localhost:8088/health and queries your service health via rpc.

<img src="images/health.svg" />

## Overview

The micro health checker provides a http endpoint at http://localhost:8088/health which queries a backend service rpc endpoint Debug.Health 
when the http endpoint is called. This is useful inside environments like kubernetes where you need to run some sort of external healthcheck.

## Run Health

Start the health server (Runs on port :8088) and specify the service/address to check

```shell
micro health --check_service=go.micro.srv.greeter --check_address=localhost:9090
```

Now query http://localhost:8088. It should return 200 OK.

## Check Health

Additonally you can check health from the command line as an RPC call

```
micro health check go.micro.srv.greeter
```

{% include links.html %}
