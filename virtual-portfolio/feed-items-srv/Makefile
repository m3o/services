build-proto:
	protoc ./proto/feeditems.proto --micro_out=. --go_out=.

run:
	DB_HOST=localhost DB_USER=postgres DB_NAME=feed-items-srv go run . 

build:
	docker build . -t eu.gcr.io/kytra-app/feeditems-srv:0.1.0