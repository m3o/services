module github.com/micro/services/m3o/api

go 1.13

require (
	github.com/golang/protobuf v1.3.5
	github.com/micro/go-micro/v2 v2.6.1-0.20200504125053-90dd1f63c853
	github.com/micro/services/payments/provider v0.0.0-20200504210551-837f046f89b2
	github.com/micro/services/projects/service v0.0.0-00010101000000-000000000000
	github.com/micro/services/users/service v0.0.0-20200421152545-96775626d99a
)

replace github.com/micro/services/projects/service => ../../projects/service
