#!/bin/bash

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64 

REPO=($1) # e.g. "micro/services"
SERVICES=($2) # e.g. "foobar barfoo helloworld"

echo Repo: $REPO
echo Services: $SERVICES

for dir in "${SERVICES[@]}"; do
    for path in $(find $dir -name "main.go"); do
        # transform the name to one suitable for a docker image
        # e.g. "helloworld/api" => "helloworld-api"
        name=$(echo $(dirname $path) | tr / -)
        echo Building $name

        # build the go binary
        go build -ldflags="-s -w" -o tmp/app $path

        # build the docker image
        tag=docker.pkg.github.com/$REPO/$name
        docker build tmp -t $tag -f .github/workflows/Dockerfile
        docker push $tag

        # remove the binary
        rm tmp/app
    done
done