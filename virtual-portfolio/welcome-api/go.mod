module github.com/kytra-app/welcome-api

go 1.12

require (
	github.com/dnaeon/go-vcr v1.0.1 // indirect
	github.com/dongri/phonenumber v0.0.0-20191029000444-38c6f1b163b5
	github.com/golang/groupcache v0.0.0-20191002201903-404acd9df4cc // indirect
	github.com/golang/protobuf v1.3.2
	github.com/kytra-app/helpers/authentication v1.0.0
	github.com/kytra-app/helpers/microgorm v1.0.0
	github.com/kytra-app/helpers/photos v0.0.0-00010101000000-000000000000
	github.com/kytra-app/helpers/sms v0.0.0-00010101000000-000000000000
	github.com/kytra-app/ledger-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/portfolios-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/sms-verification-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/users-srv v1.0.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins v1.5.1
	github.com/pkg/errors v0.8.1
	go.etcd.io/etcd v3.3.13+incompatible
)

replace github.com/kytra-app/helpers/authentication => ../helpers/authentication

replace github.com/kytra-app/helpers/passwordhasher => ../helpers/passwordhasher

replace github.com/kytra-app/helpers/microgorm => ../helpers/microgorm

replace github.com/kytra-app/helpers/microtime => ../helpers/microtime

replace github.com/kytra-app/helpers/iex-cloud => ../helpers/iex-cloud

replace github.com/kytra-app/helpers/sms => ../helpers/sms

replace github.com/kytra-app/helpers/photos => ../helpers/photos

replace github.com/kytra-app/users-srv => ../users-srv

replace github.com/kytra-app/portfolios-srv => ../portfolios-srv

replace github.com/kytra-app/sms-verification-srv => ../sms-verification-srv

replace github.com/kytra-app/ledger-srv => ../ledger-srv

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1
