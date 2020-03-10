GOPATH:=$(shell go env GOPATH)

.PHONY: build
build: proto
	go build -o account-web *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t account-web:latest
