module github.com/micro/services/account/web

go 1.13

require (
	github.com/dghubble/gologin v2.1.0+incompatible
	github.com/dghubble/gologin/v2 v2.2.0
	github.com/micro/go-micro/v2 v2.2.1-0.20200311230942-1ca4619506bd
	github.com/micro/services/login/service v0.0.0-00010101000000-000000000000
	github.com/micro/services/users/service v0.0.0-20200311145701-949f1a383199
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
)

replace github.com/micro/go-micro/v2 => ../../../go-micro

replace github.com/micro/services/login/service => ../../login/service

replace github.com/micro/services/users/service => ../../users/service
