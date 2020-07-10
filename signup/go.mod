module github.com/micro/services/signup

go 1.13

require (
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro/v2 v2.9.1-0.20200709192134-3480e0a64e21
	github.com/micro/services/account/invite v0.0.0-00010101000000-000000000000
	github.com/micro/services/kubernetes/service v0.0.0-00010101000000-000000000000
	github.com/micro/services/payments/provider v0.0.0-00010101000000-000000000000
	github.com/sethvargo/go-diceware v0.2.0
	google.golang.org/protobuf v1.25.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/services/payments/provider => ../payments/provider

replace github.com/micro/services/account/invite => ../account/invite

replace github.com/micro/services/kubernetes/service => ../kubernetes/service
