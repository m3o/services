#!/bin/bash
set -e

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64 

SERVICES=($1) # e.g. "foobar barfoo helloworld"

rootDir=$(pwd)

function build {
    dir=$1
    echo Building $dir
    cd $dir

    # build the proto buffers
    #find . -name "*.proto" | xargs --no-run-if-empty protoc --proto_path=. --micro_out=. --go_out=.  

    if [ "$dir" == "explore/web" ]; then
        go generate
    fi

    # build the binaries
    go build -ldflags="-s -w" -o service .
    cp $rootDir/dumb-init/dumb-init dumb-init

    # build the docker image
    tag=docker.pkg.github.com/micro/services/$(echo $dir | tr / -)
    docker build . -t $tag -f $rootDir/.github/workflows/Dockerfile

    # push the docker image
    echo Pushing $tag
    docker push $tag

    # remove the binaries
    rm service
    rm dumb-init

    # go back to the top level dir
    cd $rootDir
}

# This must always be deployed even if it has not changed
build "explore/web"

for dir in "${SERVICES[@]}"; do
    build $dir
done
