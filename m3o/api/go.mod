module github.com/micro/services/m3o/api

go 1.13

require (
	github.com/golang/protobuf v1.3.5
	github.com/micro/go-micro/v2 v2.5.1-0.20200430232517-e8105d22adc6
	github.com/micro/services/project/service v0.0.0-00010101000000-000000000000
	github.com/micro/services/users/service v0.0.0-20200421152545-96775626d99a
)

replace github.com/micro/services/project/service => ../../project/service
