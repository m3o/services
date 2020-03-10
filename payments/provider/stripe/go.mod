module github.com/micro/services/payments/provider/stripe

go 1.13

require (
	github.com/coreos/etcd v3.3.18+incompatible
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v2 v2.2.1-0.20200309204305-241614ff686e
	github.com/micro/services/payments/provider v0.0.0-00010101000000-000000000000
	github.com/stripe/stripe-go v70.2.0+incompatible
)

replace github.com/micro/services/payments/provider => ../
