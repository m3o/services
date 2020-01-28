build-proto:
	protoc ./proto/portfolios.proto --micro_out=. --go_out=.

run:
	DB_HOST=localhost DB_USER=postgres DB_NAME=portfolios-srv go run . 

build:
	docker build . -t eu.gcr.io/kytra-app/portfolios-srv:0.1.0