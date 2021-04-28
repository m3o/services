# Signup Service

To run the whole onboarding/onboarding flow you will need the following configs:

```
micro config set micro.payments.stripe.api_key ...
micro config set micro.onboarding.sendgrid.api_key ...
micro config set micro.onboarding.sendgrid.template_id ...
micro config set micro.onboarding.plan_id ...
```

This is the Signup service

Generated with

```
micro new --namespace=go.micro --type=service onboarding
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.service.onboarding
- Type: service
- Alias: onboarding

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./onboarding-service
```

Build a docker image
```
make docker
```