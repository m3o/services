module github.com/micro/services/login/service

go 1.13

require (
	github.com/golang/protobuf v1.3.4
	github.com/micro/go-micro/v2 v2.2.1-0.20200311230942-1ca4619506bd
	github.com/micro/services/users/service v0.0.0-20200311145701-949f1a383199
	golang.org/x/crypto v0.0.0-20200221231518-2aa609cf4a9d
)

replace github.com/micro/services/users/service => ../../users/service
