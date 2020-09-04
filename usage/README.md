# Usage Service

Get latest usage stats for different namespaces:

```
$ microadmin usage list
{
	"accounts": [
		{
			"namespace": "splicing-earthlike-salvage",
			"users": "3",
			"services": "2",
			"created": "1599129489"
		},
		{
			"namespace": "snitch-magician-morbidity",
			"users": "1",
			"services": "1",
			"created": "1599129489"
		},
    ]
}
```

Get a list of samples taken in reverse chronological order for a namespace:

```
$ microadmin usage list --namespace=john-secular-carwash
{
	"accounts": [
		{
			"namespace": "john-secular-carwash",
			"users": "1",
			"created": "1598977182"
		},
		{
			"namespace": "john-secular-carwash",
			"users": "1",
			"created": "1598980784"
		},
		{
			"namespace": "john-secular-carwash",
			"users": "1",
			"created": "1598984385"
		},
    ]
}
```

This is the Usage service

Generated with

```
micro new usage
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- Alias: usage

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
./usage
```

Build a docker image
```
make docker
```