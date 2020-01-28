build-proto:
	protoc ./proto/posts.proto --micro_out=. --go_out=.

run:
	DB_HOST=localhost DB_USER=postgres DB_NAME=posts-srv go run . 

build:
	docker build . -t eu.gcr.io/kytra-app/posts-srv:0.1.0