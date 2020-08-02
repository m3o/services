# Notifications Service

This is the Notifications service. Much like Github notifications or similar, you can subscribe to get notifications for a particular resource and all subsequent notifications for that resource will result in a notification being generated for you. 

Typical flow is
- `Notifications.Subscribe` - to register interest in something
- `Notifications.Notify` - when something happens
- `Notifications.List` - to get all the notifications
- `Notifications.MarkAsRead` - to mark notifications read
- `Notifications.Unsubscribe` - to remove interest in something

Generated with

```
micro new --namespace=go.micro --type=service notifications
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.service.notifications
- Type: service
- Alias: notifications

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
./notifications-service
```

Build a docker image
```
make docker
```