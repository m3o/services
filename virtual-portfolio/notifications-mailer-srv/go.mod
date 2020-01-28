module github.com/kytra-app/notifications-mailer-srv

go 1.12

require (
	github.com/golang/protobuf v1.3.2
	github.com/kytra-app/helpers/iex-cloud v1.0.0
	github.com/kytra-app/helpers/mailer v0.0.0-00010101000000-000000000000
	github.com/kytra-app/helpers/photos v0.0.0-00010101000000-000000000000
	github.com/kytra-app/helpers/unique v0.0.0-00010101000000-000000000000
	github.com/kytra-app/notifications-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/stocks-srv v1.0.0
	github.com/kytra-app/users-srv v1.0.0
	github.com/micro/go-micro v1.8.1
	github.com/micro/go-plugins v1.2.0
	github.com/robfig/cron v1.2.0 // indirect
	github.com/robfig/cron/v3 v3.0.0
)

replace github.com/kytra-app/stocks-srv => ../stocks-srv

replace github.com/kytra-app/notifications-srv => ../notifications-srv

replace github.com/kytra-app/comments-srv => ../comments-srv

replace github.com/kytra-app/users-srv => ../users-srv

replace github.com/kytra-app/posts-srv => ../posts-srv

replace github.com/kytra-app/feed-items-srv => ../feed-items-srv

replace github.com/kytra-app/followers-srv => ../followers-srv

replace github.com/kytra-app/helpers/iex-cloud => ../helpers/iex-cloud

replace github.com/kytra-app/helpers/unique => ../helpers/unique

replace github.com/kytra-app/helpers/mailer => ../helpers/mailer

replace github.com/kytra-app/helpers/microgorm => ../helpers/microgorm

replace github.com/kytra-app/helpers/passwordhasher => ../helpers/passwordhasher

replace github.com/kytra-app/helpers/textenhancer => ../helpers/textenhancer

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1

replace github.com/kytra-app/helpers/photos => ../helpers/photos
