module github.com/micro/services/signup

go 1.13

require (
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro/v2 v2.9.1-0.20200630164038-dcf01ebbf033
	github.com/micro/services/payments/provider v0.0.0-20200618133042-550220a6eff2
	google.golang.org/protobuf v1.25.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
