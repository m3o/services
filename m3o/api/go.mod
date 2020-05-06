module github.com/micro/services/m3o/api

go 1.13

require (
	github.com/golang/protobuf v1.4.1
	github.com/micro/go-micro/v2 v2.6.1-0.20200506095814-5387f73b5d8e
	github.com/micro/services/kubernetes/service v0.0.0-20200505140906-ca5cb95fe360
	github.com/micro/services/payments/provider v0.0.0-20200505140906-ca5cb95fe360
	github.com/micro/services/projects/environments v0.0.0-20200505140906-ca5cb95fe360
	github.com/micro/services/projects/service v0.0.0-20200505140906-ca5cb95fe360
	github.com/micro/services/users/service v0.0.0-20200505140906-ca5cb95fe360
)

replace github.com/micro/services/projects/service => ../../projects/service
