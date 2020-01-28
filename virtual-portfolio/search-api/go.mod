module github.com/kytra-app/search-api

go 1.12

require (
	github.com/abbot/go-http-auth v0.4.1-0.20181019201920-860ed7f246ff
	github.com/dnaeon/go-vcr v1.0.1 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/kytra-app/followers-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/helpers/authentication v0.0.0-00010101000000-000000000000
	github.com/kytra-app/helpers/iex-cloud v1.0.0
	github.com/kytra-app/helpers/photos v0.0.0-00010101000000-000000000000
	github.com/kytra-app/stocks-srv v1.0.0
	github.com/kytra-app/users-srv v1.0.0
	github.com/micro/go-micro v1.16.0
	github.com/micro/go-plugins v1.5.1
)

replace github.com/kytra-app/stocks-srv => ../stocks-srv

replace github.com/kytra-app/users-srv => ../users-srv

replace github.com/kytra-app/followers-srv => ../followers-srv

replace github.com/kytra-app/helpers/iex-cloud => ../helpers/iex-cloud

replace github.com/kytra-app/helpers/microgorm => ../helpers/microgorm

replace github.com/kytra-app/helpers/authentication => ../helpers/authentication

replace github.com/kytra-app/helpers/passwordhasher => ../helpers/passwordhasher

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1

replace github.com/kytra-app/helpers/photos => ../helpers/photos
