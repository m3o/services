build-proto:
	protoc ./proto/followers.proto --micro_out=. --go_out=.

run:
	DB_HOST=localhost DB_USER=postgres DB_NAME=followers-srv go run . 

build:
	docker build . -t eu.gcr.io/kytra-app/followers-srv:0.1.0