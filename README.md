# Services

This is the home of the Micro Services Platform aka (M3O).

## Overview

This repository serves as the home for the M3O platform and related services. Those invited to develop the platform will be added to the team and have the ability to create and modify services here. We'll build additional value add services here to be consumed by users of the platform.

## Design

All services are Micro services written using the Micro framework without exception.

- Services speak to each other via RPC
- Messages are used for async eventing
- Infrastructure usage occurs only through Micro

## Dependencies

We depend on Micro run using the "platform" profile. It runs on kubernetes with shared resources. 
Some of which are listed below.

- Etcd
- NATS
- CockroachDB

## Naming

Directories are the domain boundary for a specific concern e.g user, account, payment. They act as the 
alias for the otherwise fully qualified domain name "go.micro.service.alias". Services should follow 
this naming convention and focus on single word naming where possible.

## Structure

Services should be generated using the `micro new` command using the alias e.g `micro new account`. 
The internal structure is defined by our new template generator. Extending this should follow 
a further convention as follows:

```
user/
    api/	# api routes
    web/	# web html
    client/	# generated clients
    service/	# core service types
    handler/	# request handlers
    subscriber/	# message subscribers
    proto/	# proto generated code
    main.go	# service main
    user.mu	# mu definition
    README.md	# readme
```

## Contribution

Please sign-off contributions with DCO sign-off

```
git commit --signoff 'Signed-off-by: John Doe <john@example.com>`
```

## License

See [LICENSE](LICENSE) which makes use of [Polyform Strict](https://polyformproject.org/licenses/strict/1.0.0/). 
For commercial use please email [contact@m3o.com](mailto:contact@m3o.com). 
