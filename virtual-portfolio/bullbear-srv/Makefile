build-proto:
	protoc ./proto/bullbear.proto --micro_out=. --go_out=.

run:
	DB_HOST=localhost DB_USER=postgres DB_NAME=bullbear-srv go run . 

build:
	docker build . -t eu.gcr.io/kytra-app/bullbear-srv:0.1.0