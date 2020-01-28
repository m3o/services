module github.com/kytra-app/followers-srv

go 1.12

require (
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a // indirect
	github.com/go-ozzo/ozzo-validation v3.5.0+incompatible
	github.com/golang/groupcache v0.0.0-20191002201903-404acd9df4cc // indirect
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/jinzhu/gorm v1.9.10
	github.com/kytra-app/helpers/microgorm v1.0.0
	github.com/kytra-app/helpers/microtime v0.0.0-00010101000000-000000000000
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins v1.2.0
	github.com/nats-io/nats-server/v2 v2.1.0 // indirect
	github.com/pkg/errors v0.8.1
)

replace github.com/kytra-app/helpers/microtime => ../helpers/microtime

replace github.com/kytra-app/helpers/microgorm => ../helpers/microgorm
