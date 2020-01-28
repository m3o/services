build-proto:
	protoc ./proto/trades.proto --micro_out=. --go_out=.

run:
	DB_HOST=localhost DB_USER=postgres DB_NAME=trades-srv go run . 

build:
	docker build . -t eu.gcr.io/kytra-app/trades-srv:0.1.0