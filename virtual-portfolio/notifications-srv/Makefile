build-proto:
	protoc ./proto/notifications.proto --micro_out=. --go_out=.

run:
	DB_HOST=localhost DB_USER=postgres DB_NAME=notifications-srv go run . 

build:
	docker build . -t eu.gcr.io/kytra-app/notifications-srv:0.1.0