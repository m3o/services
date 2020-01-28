module github.com/kytra-app/investors-api

go 1.12

require (
	github.com/dnaeon/go-vcr v1.0.1 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/kytra-app/followers-srv v1.0.0
	github.com/kytra-app/helpers/authentication v1.0.0
	github.com/kytra-app/helpers/iex-cloud v1.0.0
	github.com/kytra-app/helpers/photos v1.0.0
	github.com/kytra-app/portfolios-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/post-enhancer-srv v1.0.0
	github.com/kytra-app/posts-srv v1.0.0
	github.com/kytra-app/stock-quote-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/stocks-srv v1.0.0
	github.com/kytra-app/trades-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/users-srv v1.0.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins v1.5.1

)

replace github.com/kytra-app/posts-srv => ../posts-srv

replace github.com/kytra-app/comments-srv => ../comments-srv

replace github.com/kytra-app/bullbear-srv => ../bullbear-srv

replace github.com/kytra-app/ledger-srv => ../ledger-srv

replace github.com/kytra-app/stocks-srv => ../stocks-srv

replace github.com/kytra-app/helpers/unique => ../helpers/unique

replace github.com/kytra-app/helpers/microtime => ../helpers/microtime

replace github.com/kytra-app/followers-srv => ../followers-srv

replace github.com/kytra-app/users-srv => ../users-srv

replace github.com/kytra-app/portfolios-srv => ../portfolios-srv

replace github.com/kytra-app/trades-srv => ../trades-srv

replace github.com/kytra-app/stock-quote-srv => ../stock-quote-srv

replace github.com/kytra-app/post-enhancer-srv => ../post-enhancer-srv

replace github.com/kytra-app/helpers/iex-cloud => ../helpers/iex-cloud

replace github.com/kytra-app/helpers/authentication => ../helpers/authentication

replace github.com/kytra-app/helpers/passwordhasher => ../helpers/passwordhasher

replace github.com/kytra-app/helpers/microgorm => ../helpers/microgorm

replace github.com/kytra-app/helpers/photos => ../helpers/photos

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1
