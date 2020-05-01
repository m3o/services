module github.com/micro/services/events/api

go 1.13

require (
	github.com/golang/protobuf v1.4.0
	github.com/micro/go-micro/v2 v2.5.1-0.20200430215125-359b8bc50305
	github.com/micro/services/events/service v0.0.0-00010101000000-000000000000
	github.com/micro/services/project/service v0.0.0-20200427143115-d065db30e6e8
)

replace github.com/micro/services/events/service => ../service
