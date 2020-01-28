module github.com/kytra-app/portfolio-valuation-srv

go 1.12

require (
	github.com/golang/groupcache v0.0.0-20191002201903-404acd9df4cc // indirect
	github.com/golang/protobuf v1.3.2
	github.com/kytra-app/ledger-srv v1.0.0
	github.com/kytra-app/stock-quote-srv v1.0.0
	github.com/kytra-app/stocks-srv v1.0.0
	github.com/kytra-app/trades-srv v0.0.0-00010101000000-000000000000
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins v1.5.1
)

replace github.com/kytra-app/helpers/iex-cloud => ../helpers/iex-cloud

replace github.com/kytra-app/helpers/microgorm => ../helpers/microgorm

replace github.com/kytra-app/helpers/microtime => ../helpers/microtime

replace github.com/kytra-app/portfolios-srv => ../portfolios-srv

replace github.com/kytra-app/stocks-srv => ../stocks-srv

replace github.com/kytra-app/stock-quote-srv => ../stock-quote-srv

replace github.com/kytra-app/trades-srv => ../trades-srv

replace github.com/kytra-app/ledger-srv => ../ledger-srv
