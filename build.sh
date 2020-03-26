#!/bin/bash
set -e

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64 

if [ -z "$1" ]; then
    SERVICES=($(find . -name main.go | cut -c 3- | rev | cut -c 9- | rev))
else
    SERVICES=($1) # e.g. "foobar barfoo helloworld"
fi

rootDir=$(pwd)

function build {
    dir=$1
    echo Building $dir
    cd $dir

    # build the proto buffers
    #find . -name "*.proto" | xargs --no-run-if-empty protoc --proto_path=. --micro_out=. --go_out=.  

    # build the binaries
    go build -ldflags="-s -w" -o service .
    cp $rootDir/dumb-init/dumb-init dumb-init

    # build the docker image
    tag=docker.pkg.github.com/micro/services/$(echo $dir | tr / -)
    docker build . -t $tag -f $rootDir/.github/workflows/Dockerfile

    if [ -z "$1" ]; then
        # push the docker image
        echo Pushing $tag
        docker push $tag
    then
        echo "Skipping pushing docker images due to lack of credentials"
    fi

    # remove the binaries
    rm service
    rm dumb-init

    # go back to the top level dir
    cd $rootDir
}

# This must always be deployed even if it has not changed
# build "explore/web"

for dir in "${SERVICES[@]}"; do
    build $dir
done
