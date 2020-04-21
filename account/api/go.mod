module github.com/micro/services/account/api

go 1.13

require (
	github.com/golang/protobuf v1.3.5
	github.com/micro/go-micro/v2 v2.5.1-0.20200420135429-7c31edd5f845
	github.com/micro/services/account/invite v0.0.0-20200421094732-38d776e22810
	github.com/micro/services/payments/provider v0.0.0-20200421094732-38d776e22810
	github.com/micro/services/teams/invites v0.0.0-20200421113609-61eca69e93f6
	github.com/micro/services/teams/service v0.0.0-20200421073553-26a9ccb4988a
	github.com/micro/services/users/service v0.0.0-20200421094732-38d776e22810
)

replace github.com/micro/services/payments/provider => ../../payments/provider
