module usage

go 1.13

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/golang/protobuf v1.4.2
	github.com/m3o/services v0.0.0-20200828124952-0a5f4a4c520f
	github.com/micro/go-micro/v3 v3.0.0-beta.0.20200825081046-bf8b3aeac796
	github.com/micro/micro/v3 v3.0.0-beta.3
	google.golang.org/protobuf v1.25.0
)
