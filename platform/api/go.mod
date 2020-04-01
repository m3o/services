module github.com/micro/services/platform/api

go 1.13

require (
	github.com/golang/protobuf v1.3.3
	github.com/micro/go-micro/v2 v2.4.1-0.20200401165638-cd3d704aa5cd
	github.com/micro/services/platform/service v0.0.0-20200313185528-4a795857eb73
	github.com/micro/services/users/service v0.0.0-20200319144224-0f47a73e0f07
)

replace github.com/micro/services/platform/service => ../service
