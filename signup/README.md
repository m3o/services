# Signup Service

To run the whole signup/onboarding flow you will need the following configs:

```
micro config set micro.payments.stripe.api_key ...
micro config set micro.signup.sendgrid.api_key ...
micro config set micro.signup.sendgrid.template_id ...
micro config set micro.signup.plan_id ...
micro config set micro.signup.key PFTTQOLVJRJTA4JROVETKSCOIU======    # random base32 encoded string
micro config set micro.signup.dkey LB3DM5JWKFHVQTRRINSXSU2DPE======   # random base32 encoded string
```

This is the Signup service

Generated with

```
micro new --namespace=go.micro --type=service signup
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.service.signup
- Type: service
- Alias: signup

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
./signup-service
```

Build a docker image
```
make docker
```