module github.com/kytra-app/post-enhancer-srv

go 1.12

require (
	github.com/golang/protobuf v1.3.2
	github.com/kytra-app/bullbear-srv v1.0.0
	github.com/kytra-app/comments-srv v1.0.0
	github.com/kytra-app/followers-srv v1.0.0
	github.com/kytra-app/helpers/authentication v1.0.0
	github.com/kytra-app/helpers/unique v1.0.0
	github.com/kytra-app/posts-srv v1.0.0
	github.com/kytra-app/stocks-srv v1.0.0
	github.com/kytra-app/users-srv v1.0.0
	github.com/micro/go-micro v1.8.1
	github.com/micro/go-plugins v1.2.0
	github.com/nats-io/nats.go v1.8.2-0.20190607221125-9f4d16fe7c2d // indirect
)

replace github.com/kytra-app/comments-srv => ../comments-srv

replace github.com/kytra-app/followers-srv => ../followers-srv

replace github.com/kytra-app/bullbear-srv => ../bullbear-srv

replace github.com/kytra-app/users-srv => ../users-srv

replace github.com/kytra-app/stocks-srv => ../stocks-srv

replace github.com/kytra-app/posts-srv => ../posts-srv

replace github.com/kytra-app/helpers/unique => ../helpers/unique

replace github.com/kytra-app/helpers/passwordhasher => ../helpers/passwordhasher

replace github.com/kytra-app/helpers/authentication => ../helpers/authentication

replace github.com/kytra-app/helpers/microgorm => ../helpers/microgorm

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1
