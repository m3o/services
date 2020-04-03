module github.com/micro/services/platform/api

go 1.13

require (
	github.com/golang/protobuf v1.3.3
	github.com/micro/go-micro/v2 v2.4.1-0.20200403120726-ed6fe67880a4
	github.com/micro/services/platform/service v0.0.0-20200313185528-4a795857eb73
	github.com/micro/services/users/service v0.0.0-20200402122209-bbd3453477a3
)

replace github.com/micro/services/platform/service => ../service
