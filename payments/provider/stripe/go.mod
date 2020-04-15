module github.com/micro/services/payments/provider/stripe

go 1.13

require (
	github.com/micro/go-micro/v2 v2.4.1-0.20200409084150-c1ad6d6c7c11
	github.com/micro/services/payments/provider v0.0.0-20200327173731-ae7aead341ae
	github.com/micro/services/users/service v0.0.0-20200313083714-e72c0c76aa9a
	github.com/stripe/stripe-go v70.2.0+incompatible
)

replace github.com/micro/services/payments/provider => ../
