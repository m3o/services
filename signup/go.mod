module github.com/micro/services/signup

go 1.13

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/golang/protobuf v1.4.2
	github.com/google/martian v2.1.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v2 v2.9.1-0.20200618113919-8c7c27c573f5
	github.com/micro/services/payments/provider v0.0.0-20200618133042-550220a6eff2
	google.golang.org/protobuf v1.25.0
)
