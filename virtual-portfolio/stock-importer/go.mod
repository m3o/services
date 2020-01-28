module github.com/kytra-app/stock-importer

go 1.12

require (
	github.com/kytra-app/helpers/photos v1.0.0
	github.com/kytra-app/stocks-srv v1.0.0
	github.com/marten-seemann/qtls v0.3.1 // indirect
	github.com/micro/go-micro v1.8.1
	github.com/micro/go-plugins v1.1.1
)

replace github.com/kytra-app/stocks-srv => ../stocks-srv

replace github.com/kytra-app/helpers/photos => ../helpers/photos
