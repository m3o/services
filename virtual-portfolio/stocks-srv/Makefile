build-proto:
	protoc ./proto/stocks.proto --micro_out=. --go_out=.

run:
	DB_HOST=localhost DB_USER=postgres DB_NAME=stocks-srv go run . 

build:
	docker build . -t eu.gcr.io/kytra-app/stocks-srv:0.1.0