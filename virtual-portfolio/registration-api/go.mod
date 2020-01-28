module github.com/kytra-app/registration-api

go 1.12

require (
	github.com/golang/protobuf v1.3.2
	github.com/kytra-app/helpers/authentication v1.0.0
	github.com/kytra-app/ledger-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/portfolios-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/users-srv v1.0.0
	github.com/marten-seemann/qtls v0.3.1 // indirect
	github.com/micro/go-micro v1.8.1
	github.com/micro/go-plugins v1.1.1
)

replace github.com/kytra-app/helpers/authentication => ../helpers/authentication

replace github.com/kytra-app/helpers/passwordhasher => ../helpers/passwordhasher

replace github.com/kytra-app/helpers/microgorm => ../helpers/microgorm

replace github.com/kytra-app/users-srv => ../users-srv

replace github.com/kytra-app/portfolios-srv => ../portfolios-srv

replace github.com/kytra-app/ledger-srv => ../ledger-srv

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1
