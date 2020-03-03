module github.com/micro/services/users/api

go 1.13

require (
	github.com/golang/protobuf v1.3.2
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v2 v2.2.0
	github.com/micro/services/users/service v0.0.0-00010101000000-000000000000
)

replace github.com/micro/services/users/service => ../service
