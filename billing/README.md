# Billing Service


Caveats:

Since currently we are doing manual reviews of subscription changes, instead of updating the users subscription automatically as their usage changes, we record the maximum value of what they used and we keep that for the month so we can do the subscription change with `billing apply` at any time during that month.

Of course the later we do the apply the less they pay for that month so there is a balance there, as at the beginning of the next month the max value will be reset. Combined with the fact that our month != their subscription month this gets a bit complicated and the only accurate solution will be updating their subscription frequently and letting Stripe handle pro rata.

```
micro new billing
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- Alias: billing

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
./billing
```

Build a docker image
```
make docker
```