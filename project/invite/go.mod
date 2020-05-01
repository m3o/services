module github.com/micro/services/project/invite

go 1.13

require (
	github.com/golang/protobuf v1.3.5
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro/v2 v2.5.1-0.20200430215125-359b8bc50305
	github.com/micro/services/project/service v0.0.0-20200421073553-26a9ccb4988a
	github.com/micro/services/users/service v0.0.0-20200421073553-26a9ccb4988a
)

replace github.com/micro/services/project/service => ../service
