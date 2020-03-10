
GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=. proto/notes/notes.proto
	protoc --proto_path=. --micro_out=. --go_out=. proto/sprints/sprints.proto
    

.PHONY: build
build: proto

	go build -o distributed-srv *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t distributed-srv:latest
