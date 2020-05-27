module github.com/micro/services/m3o/api

go 1.13

require (
	github.com/golang/protobuf v1.4.1
	github.com/micro/go-micro/v2 v2.7.1-0.20200523154723-bd049a51e637
	github.com/micro/services/kubernetes/service v0.0.0-20200505140906-ca5cb95fe360
	github.com/micro/services/payments/provider v0.0.0-20200505140906-ca5cb95fe360
	github.com/micro/services/projects/environments v0.0.0-20200505140906-ca5cb95fe360
	github.com/micro/services/projects/invite v0.0.0-20200507152129-b87672dd87ae
	github.com/micro/services/projects/service v0.0.0-20200505140906-ca5cb95fe360
	github.com/micro/services/secrets/service v0.0.0-20200511130034-26a9657c0c03
	github.com/micro/services/users/service v0.0.0-20200505140906-ca5cb95fe360
)

replace github.com/micro/services/users/service => ../../users/service

replace github.com/micro/services/projects/service => ../../projects/service

replace github.com/micro/services/kubernetes/service => ../../kubernetes/service
