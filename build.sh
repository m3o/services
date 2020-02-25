#!/bin/bash

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64 

SERVICES=($1) # e.g. "foobar barfoo helloworld"

echo Services: $SERVICES

for dir in "${SERVICES[@]}"; do
    for path in $(find $dir -name "main.go"); do
        dir=$(dirname $path)
        echo Building $dir

        # build the binaries
        go build -ldflags="-s -w" -o $dir/app $path
        cp dumb-init/dumb-init $dir/dumb-init

        # build the docker image
        tag=docker.pkg.github.com/micro/services/$(echo $dir | tr / -)
        docker build $dir -t $tag -f .github/workflows/Dockerfile

        # push the docker image
        echo Pushing $tag
        # docker push $tag

        # remove the binaries
        rm $dir/app
        rm $dir/dumb-init
    done
done