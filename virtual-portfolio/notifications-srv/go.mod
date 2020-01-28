module github.com/kytra-app/notifications-srv

go 1.12

require (
	github.com/go-ozzo/ozzo-validation v3.5.0+incompatible
	github.com/golang/protobuf v1.3.2
	github.com/jinzhu/gorm v1.9.10
	github.com/kytra-app/comments-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/feed-items-srv v1.0.0
	github.com/kytra-app/followers-srv v1.0.0
	github.com/kytra-app/helpers/microgorm v1.0.0
	github.com/kytra-app/helpers/textenhancer v0.0.0-00010101000000-000000000000
	github.com/kytra-app/posts-srv v1.0.0
	github.com/kytra-app/push-notifications-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/stocks-srv v1.0.0
	github.com/kytra-app/users-srv v1.0.0
	github.com/micro/go-micro v1.8.1
	github.com/micro/go-plugins v1.2.0
	github.com/nats-io/nats.go v1.8.2-0.20190607221125-9f4d16fe7c2d // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271
	go.etcd.io/etcd v3.3.13+incompatible
)

replace github.com/kytra-app/helpers/passwordhasher => ../helpers/passwordhasher

replace github.com/kytra-app/helpers/sms => ../helpers/sms

replace github.com/kytra-app/helpers/microgorm => ../helpers/microgorm

replace github.com/kytra-app/helpers/textenhancer => ../helpers/textenhancer

replace github.com/kytra-app/users-srv => ../users-srv

replace github.com/kytra-app/stocks-srv => ../stocks-srv

replace github.com/kytra-app/followers-srv => ../followers-srv

replace github.com/kytra-app/push-notifications-srv => ../push-notifications-srv

replace github.com/kytra-app/posts-srv => ../posts-srv

replace github.com/kytra-app/comments-srv => ../comments-srv

replace github.com/kytra-app/feed-items-srv => ../feed-items-srv

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1
