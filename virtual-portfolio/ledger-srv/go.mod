module github.com/kytra-app/ledger-srv

go 1.12

require (
	github.com/go-ozzo/ozzo-validation v3.5.0+incompatible
	github.com/golang/protobuf v1.3.2
	github.com/jinzhu/gorm v1.9.10
	github.com/kytra-app/helpers/microgorm v0.0.0-00010101000000-000000000000
	github.com/kytra-app/helpers/microtime v0.0.0-00010101000000-000000000000
	github.com/kytra-app/portfolios-srv v0.0.0-00010101000000-000000000000
	github.com/micro/go-micro v1.8.1
	github.com/micro/go-plugins v1.1.1
	github.com/pkg/errors v0.8.1
	github.com/satori/go.uuid v1.2.0
	github.com/shopspring/decimal v0.0.0-20190905144223-a36b5d85f337 // indirect
)

replace github.com/kytra-app/helpers/microgorm => ../helpers/microgorm

replace github.com/kytra-app/helpers/microtime => ../helpers/microtime

replace github.com/kytra-app/portfolios-srv => ../portfolios-srv
