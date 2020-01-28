module github.com/kytra-app/stock-target-price-srv

go 1.13

require (
	github.com/fatih/structs v1.1.0
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/golang/protobuf v1.3.2
	github.com/jinzhu/gorm v1.9.11
	github.com/kytra-app/helpers/iex-cloud v1.0.0
	github.com/kytra-app/helpers/microgorm v1.0.0
	github.com/kytra-app/helpers/microtime v0.0.0-00010101000000-000000000000
	github.com/kytra-app/stocks-srv v1.0.0
	github.com/kytra-app/trades-srv v0.0.0-00010101000000-000000000000
	github.com/micro/go-micro v1.16.0

	github.com/micro/go-plugins v1.5.1
	github.com/micro/micro v1.16.0
	github.com/pkg/errors v0.8.1
	github.com/satori/go.uuid v1.2.0
	golang.org/dl v0.0.0-20191111193948-37d848e6a9e1 // indirect
)

replace github.com/kytra-app/helpers/iex-cloud => ../helpers/iex-cloud

replace github.com/kytra-app/helpers/microgorm => ../helpers/microgorm

replace github.com/kytra-app/helpers/microtime => ../helpers/microtime

replace github.com/kytra-app/stocks-srv => ../stocks-srv

replace github.com/kytra-app/trades-srv => ../trades-srv

replace github.com/kytra-app/portfolios-srv => ../portfolios-srv

replace github.com/kytra-app/ledger-srv => ../ledger-srv
