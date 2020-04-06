module github.com/micro/services/platform/api

go 1.13

require (
	github.com/golang/protobuf v1.3.3
	github.com/micro/go-micro/v2 v2.4.1-0.20200406115547-bea7c3f7e720
	github.com/micro/micro/v2 v2.4.0
	github.com/micro/services/platform/service v0.0.0
	github.com/micro/services/users/service v0.0.0
)

replace github.com/micro/services/platform/service => ../service

replace github.com/micro/services/users/service => ../../users/service
