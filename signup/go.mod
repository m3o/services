module github.com/micro/services/signup

go 1.13

require (
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro/v2 v2.9.1-0.20200703133825-f99b436ec2fb
	github.com/micro/services/account/invite v0.0.0-20200703161509-85801894d6e0
	github.com/micro/services/payments/provider v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.25.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/services/payments/provider => ../payments/provider

replace github.com/micro/services/account/invite => ../account/invite
