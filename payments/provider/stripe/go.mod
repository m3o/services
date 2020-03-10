module github.com/micro/services/payments/provider/stripe

go 1.13

replace github.com/micro/services/payments/provider => ../

require (
	github.com/micro/go-micro/v2 v2.2.0
	github.com/micro/services/payments/provider v0.0.0-00010101000000-000000000000
)
