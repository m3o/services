module github.com/micro/services/account/api

go 1.13

require (
	github.com/golang/protobuf v1.3.5
	github.com/micro/go-micro/v2 v2.4.1-0.20200406124040-c4442a753323
	github.com/micro/services/account/invite v0.0.0-00010101000000-000000000000
	github.com/micro/services/payments/provider v0.0.0-20200331171103-a3eba43a815a
	github.com/micro/services/users/service v0.0.0-20200401191043-bafc59c2e760
)

replace github.com/micro/services/users/service => ../../users/service

replace github.com/micro/services/account/invite => ../invite
