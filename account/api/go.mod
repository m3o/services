module github.com/micro/services/account/api

go 1.13

require (
	github.com/golang/protobuf v1.3.4
	github.com/micro/go-micro/v2 v2.4.1-0.20200331151804-26747906947f
	github.com/micro/services/login/service v0.0.0-20200313083714-e72c0c76aa9a
	github.com/micro/services/payments/provider v0.0.0-20200331171103-a3eba43a815a
	github.com/micro/services/users/service v0.0.0-20200313083714-e72c0c76aa9a
)

replace github.com/micro/go-micro/v2 => ../../../go-micro

replace github.com/micro/services/users/service => ../../users/service
