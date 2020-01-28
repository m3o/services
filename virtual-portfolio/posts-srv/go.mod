module github.com/kytra-app/posts-srv

go 1.12

require (
	github.com/fatih/structs v1.1.0
	github.com/go-ozzo/ozzo-validation v3.5.0+incompatible
	github.com/golang/protobuf v1.3.2
	github.com/jinzhu/gorm v1.9.10
	github.com/kytra-app/helpers/microgorm v0.0.0-00010101000000-000000000000
	github.com/micro/go-micro v1.8.1
	github.com/micro/go-plugins v1.1.1
	github.com/pkg/errors v0.8.1
	github.com/satori/go.uuid v1.2.0
)

replace github.com/kytra-app/helpers/microgorm => ../helpers/microgorm
