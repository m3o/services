module github.com/kytra-app/sms-verification-srv

go 1.12

require (
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a // indirect
	github.com/go-ozzo/ozzo-validation v3.5.0+incompatible
	github.com/golang/protobuf v1.3.2
	github.com/jinzhu/gorm v1.9.10
	github.com/kytra-app/helpers/iex-cloud v1.0.0
	github.com/kytra-app/helpers/microgorm v1.0.0
	github.com/kytra-app/helpers/sms v0.0.0-00010101000000-000000000000
	github.com/lucas-clemente/quic-go v0.11.2 // indirect
	github.com/micro/go-micro v1.8.1
	github.com/micro/go-plugins v1.1.1
	github.com/pkg/errors v0.8.1
	github.com/satori/go.uuid v1.2.0
)

replace github.com/kytra-app/helpers/sms => ../helpers/sms

replace github.com/kytra-app/helpers/microgorm => ../helpers/microgorm

replace github.com/kytra-app/helpers/iex-cloud => ../helpers/iex-cloud
