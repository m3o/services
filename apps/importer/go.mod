module github.com/micro/services/apps/importer

go 1.13

require (
	github.com/golang/protobuf v1.3.4
	github.com/micro/go-micro/v2 v2.2.1-0.20200317112720-ab7312706317
	github.com/micro/services/apps/service v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/micro/services/apps/service => ../service
