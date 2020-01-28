module github.com/kytra-app/feed-generator-srv

go 1.12

require (
	github.com/kytra-app/feed-items-srv v1.0.0
	github.com/kytra-app/followers-srv v1.0.0
	github.com/kytra-app/posts-srv v1.0.0
	github.com/kytra-app/stocks-srv v1.0.0
	github.com/kytra-app/users-srv v1.0.0
	github.com/micro/go-micro v1.8.1
	github.com/micro/go-plugins v1.2.0
)

replace github.com/kytra-app/helpers/passwordhasher => ../helpers/passwordhasher

replace github.com/kytra-app/helpers/microgorm => ../helpers/microgorm

replace github.com/kytra-app/users-srv => ../users-srv

replace github.com/kytra-app/stocks-srv => ../stocks-srv

replace github.com/kytra-app/followers-srv => ../followers-srv

replace github.com/kytra-app/posts-srv => ../posts-srv

replace github.com/kytra-app/feed-items-srv => ../feed-items-srv

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1
