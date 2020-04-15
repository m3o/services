---
title: Enterprise Toolkit
keywords: enterprise
tags: [micro, enterprise]
sidebar: home_sidebar
permalink: /enterprise.html
summary: 
---

Micro Enterprise is a turn-key solution for building production ready cloud-native applications.

# Overview

Micro Enterprise is an enterprise version of the micro toolkit. It includes additional features which are absolutely critical for a 
stable and performant production grade microservice solution. This version of Micro requires a 
[Micro Enterprise License Agreement](https://github.com/micro/enterprise/blob/master/LICENSE) commercial subscription.

## Getting Started

- [Create Subscription](#create-subscription)
- [Manage Subscription](#manage-subscription)
- [Generate Token](#generate-token)
- [Manage Tokens](#manage-tokens)
- [Generate License](generate-license)
- [Manage Licenses](manage-licenses)
- [Use License](use-license)

## Create Subscription

Micro Enterprise is distributed under a commercial license. While the source is available on GitHub you'll need to purchase a 
subscription to use it. Head to the [pricing](https://micro.mu/pricing/) to subscribe using a credit card.

### Pricing

Micro Enterprise is priced on a per service basis. Each service registered under service discovery is billed as a unique service. 
Instances of a service count under a single service so you may run as many instances as you like.

### Example

Services "user", "order", "customer", "account", "login" are each unique services, billed at $25/service/month. This means a 
monthly total of $125/month. Each service can have unlimited instances. 

## Manage Subscription

Browse to [https://micro.chargebeeportal.com/](https://micro.chargebeeportal.com/) to manage your billing and subscription.

## Generate Token

An API token is required to manage license and subscriptions. The API token should be used with any request.

Generate the token using your subscription email address

```
micro token generate --email=user@example.com
```

You will be sent an OTP pass to enter on the prompt

```
micro token generate --email=user@example.com
Enter OTP: 632968
Your token (set as MICRO_TOKEN_KEY env var or X-Micro-Token http header):
1Qs2t4ueli9xVQz4HM0wfL3WXcPKNKMJYvB3SCOYAIsyuUBUsXIq2Yf55eFMyViBjnczXotTJ57Jb8tNqghg66b8ITZVVxVrR286F9HCaA8JcppTUjpi2UbHsAPaPFd37NiRe6hyZ6TvmujxxW82Lj9ZFZ59bIHwqJ7Z4atOxlT3P2NqvdXbNM7mvoYuxhPA32RzZ
```

Export the token. It will be required for any `micro license` command.

```
export MICRO_TOKEN_KEY=1Qs2t4ueli9xVQz4HM0wfL3WXcPKNKMJYvB3SCOYAIsyuUBUsXIq2Yf55eFMyViBjnczXotTJ57Jb8tNqghg66b8ITZVVxVrR286F9HCaA8JcppTUjpi2UbHsAPaPFd37NiRe6hyZ6TvmujxxW82Lj9ZFZ59bIHwqJ7Z4atOxlT3P2NqvdXbNM7mvoYuxhPA32RzZ
```

## Manage Tokens

List your tokens

```
micro token list

[
	{
		"id": "a586f641-be91-11e8-8287-1079ddac194b",
		"expires": 1538244139,
		"claims": {
			"email": "user@example.com"
		},
		"key": "1Qs2t4ueli9xVQz4HM0wfL3WXcPKNKMJYvB3SCOYAIsyuUBUsXIq2Yf55eFMyViBjnczXotTJ57Jb8tNqghg66b8ITZVVxVrR286F9HCaA8JcppTUjpi2UbHsAPaPFd37NiRe6hyZ6TvmujxxW82Lj9ZFZ59bIHwqJ7Z4atOxlT3P2NqvdXbNM7mvoYuxhPA32RzZ"
	}
]
```

Revoke your tokens

```
micro token revoke --token=1Qs2t4ueli9xVQz4HM0wfL3WXcPKNKMJYvB3SCOYAIsyuUBUsXIq2Yf55eFMyViBjnczXotTJ57Jb8tNqghg66b8ITZVVxVrR286F9HCaA8JcppTUjpi2UbHsAPaPFd37NiRe6hyZ6TvmujxxW82Lj9ZFZ59bIHwqJ7Z4atOxlT3P2NqvdXbNM7mvoYuxhPA32RzZ
```

## Generate License

Every micro enterprise app needs to use a license. You can generate as many licenses as you need or just use one.

Before doing so you need to get your subscription id to generate the license for.

```
micro license subscriptions
[
	{
		"id": "1mkVvwvQunJ6KT1uG",
		"email": "user@example.com",
		"type": "service",
		"status": "active",
		"created": 1528838195,
		"quantity": 5
	}
]
```

Now generate the license

```
micro license generate --subscription=1mkVvwvQunJ6KT1uG
Your license (set as MICRO_LICENSE_KEY env var or X-Micro-License http header):
FLMQn8cQJgXfSNMFgGxYKUyxjSOyZ5t6IREePHYczRuKTUN22qmPVXOmVYD7vZPX8bNwfuJHT95Mw2wtdNbBatIKDrlWn9YFKTwDi8k6BG22Rnvp6edZGSsjfBJFj79LvVDbl93LUPQqwmMs7UdSY8VNOdf61zU6yxK9La7GPSeXNsRMTIQCJNygsE5fm7GeiJTd3s4rfskambmBJARrHeS70T73Jl4mB3hSbra1J5Ww51itrtomoq8mC5TX1yRDL4Us6sKMd4aRDiTQzDTPdSMEVlVrYtbtUQN8CHKjhHfQxKsVOE8BUUDggugmrpJzYH8mmzlH31VppV1CqjWlXtbTtZ1cm8kljjeCmAs
```

## Manage Licenses

List your licenses

```
micro license list
[
	{
		"id": "5afa98c9-bf09-11e8-b8bd-b54ebe697adf",
		"version": "20180905",
		"created": 1537690753,
		"key": "FLMQn8cQJgXfSNMFgGxYKUyxjSOyZ5t6IREePHYczRuKTUN22qmPVXOmVYD7vZPX8bNwfuJHT95Mw2wtdNbBatIKDrlWn9YFKTwDi8k6BG22Rnvp6edZGSsjfBJFj79LvVDbl93LUPQqwmMs7UdSY8VNOdf61zU6yxK9La7GPSeXNsRMTIQCJNygsE5fm7GeiJTd3s4rfskambmBJARrHeS70T73Jl4mB3hSbra1J5Ww51itrtomoq8mC5TX1yRDL4Us6sKMd4aRDiTQzDTPdSMEVlVrYtbtUQN8CHKjhHfQxKsVOE8BUUDggugmrpJzYH8mmzlH31VppV1CqjWlXtbTtZ1cm8kljjeCmAs",
		"subscription": {
			"id": "1mkVvwvQunJ6KT1uG",
			"email": "user@example.com",
			"type": "service",
			"status": "active",
			"created": 1528838195,
			"quantity": 5
		}
	}
]
```

Revoke your license

```
micro license revoke --license=FLMQn8cQJgXfSNMFgGxYKUyxjSOyZ5t6IREePHYczRuKTUN22qmPVXOmVYD7vZPX8bNwfuJHT95Mw2wtdNbBatIKDrlWn9YFKTwDi8k6BG22Rnvp6edZGSsjfBJFj79LvVDbl93LUPQqwmMs7UdSY8VNOdf61zU6yxK9La7GPSeXNsRMTIQCJNygsE5fm7GeiJTd3s4rfskambmBJARrHeS70T73Jl4mB3hSbra1J5Ww51itrtomoq8mC5TX1yRDL4Us6sKMd4aRDiTQzDTPdSMEVlVrYtbtUQN8CHKjhHfQxKsVOE8BUUDggugmrpJzYH8mmzlH31VppV1CqjWlXtbTtZ1cm8kljjeCmAs
Licence revoked
```

## Use License

An enterprise license should be used with any service using Micro Enterprise.

You can either specify it as the environment variable `MICRO_LICENSE_KEY` or as an option to the service like so

```
package main

import (
	"github.com/micro/go-micro/v2"
	emicro "github.com/micro/enterprise/go/micro"
)

func main() {
	service := emicro.NewService(
		micro.Name("foo"),
		emicro.License("the key"),
	)
	service.Run()
}
```

The service will periodically ping the micro license api to track usage and update the subscription quantity. 

{% include links.html %}
