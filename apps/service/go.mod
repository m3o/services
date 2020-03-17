module github.com/micro/services/apps/service

go 1.13

require (
	github.com/golang/protobuf v1.3.4
	github.com/micro/go-micro/v2 v2.2.1-0.20200314211841-0449138f61f5
	github.com/micro/services/payments/provider v0.0.0-20200313105341-66ab9d74feb7
)

replace github.com/micro/services/payments/provider => ../../payments/provider
