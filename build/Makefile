.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=. proto/build.proto

.PHONY: build
build: proto
	go build -o build-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t build-service:latest
