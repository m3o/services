module github.com/kytra-app/stock-news-srv

go 1.12

require (
	contrib.go.opencensus.io/exporter/ocagent v0.5.1 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.12.4 // indirect
	github.com/Azure/azure-pipeline-go v0.2.2 // indirect
	github.com/Azure/azure-storage-blob-go v0.7.0 // indirect
	github.com/Azure/go-autorest v12.3.0+incompatible // indirect
	github.com/GoogleCloudPlatform/cloudsql-proxy v0.0.0-20190725230627-253d1edd4416 // indirect
	github.com/RoaringBitmap/roaring v0.4.18 // indirect
	github.com/anacrolix/tagflag v1.0.0 // indirect
	github.com/dnaeon/go-vcr v1.0.1 // indirect
	github.com/fatih/structs v1.1.0
	github.com/glycerine/go-unsnap-stream v0.0.0-20190730064659-98d31706395a // indirect
	github.com/glycerine/goconvey v0.0.0-20190410193231-58a59202ab31 // indirect
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/golang/protobuf v1.3.2
	github.com/google/pprof v0.0.0-20190723021845-34ac40c74b70 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.5 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jinzhu/gorm v1.9.11
	github.com/kytra-app/helpers/microgorm v1.0.0
	github.com/kytra-app/helpers/microtime v0.0.0-00010101000000-000000000000
	github.com/kytra-app/helpers/news v0.0.0-00010101000000-000000000000
	github.com/kytra-app/insights-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/stock-movers-srv v0.0.0-00010101000000-000000000000
	github.com/kytra-app/stocks-srv v1.0.0
	github.com/kytra-app/trades-srv v0.0.0-00010101000000-000000000000
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins v1.5.1
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/pkg/errors v0.8.1
	github.com/rcrowley/go-metrics v0.0.0-20190706150252-9beb055b7962 // indirect
	github.com/robfig/cron/v3 v3.0.0
	github.com/rogpeppe/fastuuid v1.2.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/smartystreets/assertions v1.0.1 // indirect
	github.com/smartystreets/goconvey v0.0.0-20190710185942-9d28bd7c0945 // indirect
	golang.org/x/exp v0.0.0-20190718202018-cfdd5522f6f6 // indirect
	golang.org/x/image v0.0.0-20190729225735-1bd0cf576493 // indirect
	golang.org/x/mobile v0.0.0-20190719004257-d2bd2a29d028 // indirect
	gopkg.in/jcmturner/gokrb5.v7 v7.3.0 // indirect
	pack.ag/amqp v0.12.0 // indirect
)

replace github.com/kytra-app/helpers/news => ../helpers/news

replace github.com/kytra-app/helpers/microgorm => ../helpers/microgorm

replace github.com/kytra-app/helpers/microtime => ../helpers/microtime

replace github.com/kytra-app/helpers/iex-cloud => ../helpers/iex-cloud

replace github.com/kytra-app/helpers/passwordhasher => ../helpers/passwordhasher

replace github.com/kytra-app/stocks-srv => ../stocks-srv

replace github.com/kytra-app/trades-srv => ../trades-srv

replace github.com/kytra-app/portfolios-srv => ../portfolios-srv

replace github.com/kytra-app/portfolio-valuation-srv => ../portfolio-valuation-srv

replace github.com/kytra-app/followers-srv => ../followers-srv

replace github.com/kytra-app/users-srv => ../users-srv

replace github.com/kytra-app/stock-earnings-srv => ../stock-earnings-srv

replace github.com/kytra-app/stock-quote-srv => ../stock-quote-srv

replace github.com/kytra-app/stock-movers-srv => ../stock-movers-srv

replace github.com/kytra-app/ledger-srv => ../ledger-srv

replace github.com/kytra-app/insights-srv => ../insights-srv
