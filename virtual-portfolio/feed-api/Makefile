build-proto:
	protoc ./proto/feed.proto --micro_out=. --go_out=.

run:
	go run .

build:
	docker build . -t eu.gcr.io/kytra-app/feed-api:0.1.0

deploy:
	docker push eu.gcr.io/kytra-app/feed-api:0.1.0
	kubectl apply -f deployment/service.yaml