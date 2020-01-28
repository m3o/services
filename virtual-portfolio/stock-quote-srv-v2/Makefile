build-proto:
	protoc ./proto/stock-quote.proto --micro_out=. --go_out=.

run:
	DB_HOST=localhost DB_USER=postgres DB_NAME=stock-quote-srv go run . 

build:
	docker build . -t eu.gcr.io/kytra-app/stock-quote-srv:0.1.0