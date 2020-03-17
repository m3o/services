module github.com/micro/services/home/api

go 1.13

require (
	github.com/golang/protobuf v1.3.4
	github.com/micro/go-micro/v2 v2.2.1-0.20200317112720-ab7312706317
	github.com/micro/services/apps/service v0.0.0-00010101000000-000000000000
	github.com/micro/services/users/service v0.0.0-20200313151537-5407234f5db7
)

replace github.com/micro/services/apps/service => ../../apps/service
