
GOPATH:=$(shell go env GOPATH)
MODIFY=Mproto/imports/registry.proto=github.com/micro/micro/v3/proto/registry

.PHONY: proto
proto:
    
	protoc --proto_path=. --micro_out=${MODIFY}:. --go_out=${MODIFY}:. proto/explore/explore.proto
    

.PHONY: build
build: proto

	go build -o explore *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t explore	:latest
