#!/bin/bash

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64 

SERVICES=($1) # e.g. "foobar barfoo helloworld"

for dir in "${SERVICES[@]}"; do
    echo Building $dir

    # build the binaries
    go build -ldflags="-s -w" -o $dir/app $dir/main.go
    cp dumb-init/dumb-init $dir/dumb-init

    # build the docker image
    tag=docker.pkg.github.com/micro/services/$(echo $dir | tr / -)
    docker build $dir -t $tag -f .github/workflows/Dockerfile

    # push the docker image
    echo Pushing $tag
    docker push $tag

    # remove the binaries
    rm $dir/app
    rm $dir/dumb-init
done
