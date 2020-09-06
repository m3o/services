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

.PHONY: images
images:
	docker build . -t rg.fr-par.scw.cloud/build/base:latest -f builder/Dockerfile.base
	docker build . -t rg.fr-par.scw.cloud/build/build:latest -f builder/Dockerfile.build
	docker push rg.fr-par.scw.cloud/build/base:latest
	docker push rg.fr-par.scw.cloud/build/build:latest
