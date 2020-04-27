module github.com/micro/services/event/api

go 1.13

require (
	github.com/golang/protobuf v1.4.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v2 v2.5.1-0.20200422094434-e25ab9f4ca28
	github.com/micro/services/event/service v0.0.0-00010101000000-000000000000
	github.com/micro/services/project/service v0.0.0-20200427143115-d065db30e6e8
)

replace github.com/micro/services/event/service => ../service
