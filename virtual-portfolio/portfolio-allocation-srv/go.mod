module github.com/kytra-app/portfolio-allocation-srv

go 1.12

require (
	github.com/fatih/structs v1.1.0
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/golang/protobuf v1.3.2
	github.com/jinzhu/gorm v1.9.11
	github.com/kytra-app/followers-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/helpers/iex-cloud v1.0.0
	github.com/kytra-app/helpers/microgorm v1.0.0
	github.com/kytra-app/helpers/microtime v0.0.0-00010101000000-000000000000
	github.com/kytra-app/helpers/unique v0.0.0-00010101000000-000000000000
	github.com/kytra-app/ledger-srv v1.0.0
	github.com/kytra-app/portfolio-valuation-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/portfolio-value-tracking-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/portfolios-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/stock-earnings-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/stock-quote-srv-v2 v0.0.0-00010101000000-000000000000
	github.com/kytra-app/stock-target-price-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/stocks-srv v1.0.0
	github.com/kytra-app/trades-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/users-srv v0.0.0-00010101000000-000000000000
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins v1.5.1
	github.com/pkg/errors v0.8.1
	github.com/satori/go.uuid v1.2.0
)

replace github.com/kytra-app/helpers/iex-cloud => ../helpers/iex-cloud

replace github.com/kytra-app/helpers/microgorm => ../helpers/microgorm

replace github.com/kytra-app/helpers/microtime => ../helpers/microtime

replace github.com/kytra-app/helpers/passwordhasher => ../helpers/passwordhasher

replace github.com/kytra-app/helpers/worldtradingdata => ../helpers/worldtradingdata

replace github.com/kytra-app/helpers/unique => ../helpers/unique

replace github.com/kytra-app/stocks-srv => ../stocks-srv

replace github.com/kytra-app/insights-srv => ../insights-srv

replace github.com/kytra-app/trades-srv => ../trades-srv

replace github.com/kytra-app/followers-srv => ../followers-srv

replace github.com/kytra-app/portfolios-srv => ../portfolios-srv

replace github.com/kytra-app/stock-earnings-srv => ../stock-earnings-srv

replace github.com/kytra-app/stock-target-price-srv => ../stock-target-price-srv

replace github.com/kytra-app/stock-quote-srv-v2 => ../stock-quote-srv-v2

replace github.com/kytra-app/ledger-srv => ../ledger-srv

replace github.com/kytra-app/portfolio-valuation-srv => ../portfolio-valuation-srv

replace github.com/kytra-app/portfolio-value-tracking-srv => ../portfolio-value-tracking-srv

replace github.com/kytra-app/stock-quote-srv => ../stock-quote-srv

replace github.com/kytra-app/users-srv => ../users-srv
