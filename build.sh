#!/bin/bash

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64 

SERVICES=($1) # e.g. "foobar barfoo helloworld"

rootDir=$(pwd)
for dir in "${SERVICES[@]}"; do
    echo Building $dir
    cd $dir

    # build the binaries
    go build -ldflags="-s -w" -o app .
    cp $rootDir/dumb-init/dumb-init dumb-init

    # build the docker image
    tag=docker.pkg.github.com/micro/services/$(echo $dir | tr / -)
    docker build . -t $tag -f .github/workflows/Dockerfile

    # push the docker image
    echo Pushing $tag
    docker push $tag

    # remove the binaries
    rm $dir/app
    rm $dir/dumb-init

    # go back to the top level dir
    cd $rootDir
done
