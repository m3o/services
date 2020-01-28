build-proto:
	protoc ./proto/users.proto --micro_out=. --go_out=.

run:
	DB_HOST=localhost DB_USER=postgres DB_NAME=user-srv go run . 

build:
	docker build . -t eu.gcr.io/kytra-app/user-srv:0.1.0