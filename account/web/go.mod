module github.com/micro/services/account/web

go 1.13

require (
	github.com/micro/go-micro/v2 v2.2.1-0.20200313093044-fbde872e7f02
	github.com/micro/services/login/service v0.0.0-00010101000000-000000000000
	github.com/micro/services/users/service v0.0.0-20200311145701-949f1a383199
)

replace github.com/micro/services/login/service => ../../login/service

replace github.com/micro/services/users/service => ../../users/service
