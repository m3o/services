build-proto:
	protoc ./proto/registration.proto --micro_out=. --go_out=.

run:
	go run .

build:
	docker build . -t eu.gcr.io/kytra-app/registration-api:0.1.1

deploy:
	docker push eu.gcr.io/kytra-app/registration-api:0.1.1
	kubectl apply -f deployment/service.yaml