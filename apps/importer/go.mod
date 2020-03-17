module github.com/micro/services/apps/importer

go 1.13

require (
	github.com/go-delve/delve v1.4.0
	github.com/golang/protobuf v1.3.4
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v2 v2.2.1-0.20200314211841-0449138f61f5
	github.com/micro/services/apps/service v0.0.0-00010101000000-000000000000
	google.golang.org/appengine v1.6.1
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/micro/services/apps/service => ../service
