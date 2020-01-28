build-proto:
	protoc ./proto/comments.proto --micro_out=. --go_out=.

run:
	DB_HOST=localhost DB_USER=postgres DB_NAME=comments-srv go run . 

build:
	docker build . -t eu.gcr.io/kytra-app/comments-srv:0.1.0