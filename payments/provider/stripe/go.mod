module github.com/micro/services/payments/provider/stripe

go 1.13

require (
	github.com/micro/go-micro/v2 v2.9.1-0.20200703133825-f99b436ec2fb
	github.com/micro/services/payments/provider v0.0.0-00010101000000-000000000000
	github.com/micro/services/users/service v0.0.0-20200313083714-e72c0c76aa9a
	github.com/stripe/stripe-go v70.2.0+incompatible
)

replace github.com/micro/services/payments/provider => ../

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
