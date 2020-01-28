build-proto:
	protoc ./proto/search.proto --micro_out=. --go_out=.

run:
	go run .

build:
	docker build . -t eu.gcr.io/kytra-app/search-api:0.1.0

deploy:
	docker push eu.gcr.io/kytra-app/search-api:0.1.0
	kubectl apply -f deployment/service.yaml