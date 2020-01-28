build-proto:
	protoc ./proto/portfolio_valuation.proto --micro_out=. --go_out=.

run:
	DB_HOST=localhost DB_USER=postgres DB_NAME=portfolio-valuation-srv go run . 

build:
	docker build . -t eu.gcr.io/kytra-app/portfolio-valuation-srv:0.1.0