build-proto:
	protoc ./proto/ledger.proto --micro_out=. --go_out=.

run:
	DB_HOST=localhost DB_USER=postgres DB_NAME=ledger-srv go run . 

build:
	docker build . -t eu.gcr.io/kytra-app/ledger-srv:0.1.0