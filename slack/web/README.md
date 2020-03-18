# Slack Service

This is the Slack service (currently just redirect)

Generated with

```
micro new web --namespace=go.micro --alias=slack --type=web
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.web.slack
- Type: web
- Alias: slack

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
./slack-web
```

Build a docker image
```
make docker
```
