module github.com/micro/services/account/web

go 1.13

require (
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro/v2 v2.5.1-0.20200421145440-d7ecb58f6cf6
	github.com/micro/services/login/service v0.0.0-20200313083714-e72c0c76aa9a
	github.com/micro/services/teams/invites v0.0.0-20200421101014-4b009b48a425
	github.com/micro/services/users/service v0.0.0-20200421073553-26a9ccb4988a
)

replace github.com/micro/services/teams/invites => ../../teams/invites
