build-proto:
	protoc ./proto/feed-generator.proto --micro_out=. --go_out=.

run:
	DB_HOST=localhost DB_USER=postgres DB_NAME=feed-generator-srv go run . 

build:
	docker build . -t eu.gcr.io/kytra-app/feed-generator-srv:0.1.0