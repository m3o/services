#!/bin/bash
set -e

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64 



# Might not always have services passed down -
# Github Actions needs GITHUB_TOKEN and for PR forks we don't have that.
if [ -n "$PR_NUMBER" ] && [ -z "$1" ]; then
    echo "Getting files from github api to detect files changed "
    SERVICES=($(find . -name main.go | cut -c 3- | rev | cut -c 9- | rev))
    URL="https://api.github.com/repos/$GITHUB_REPOSITORY/pulls/$PR_NUMBER/files"
    FILES=($(curl -s -X GET -G $URL | jq -r '.[] | .filename'))
else
    SERVICES=($1) # e.g. "foobar barfoo helloworld"
fi

rootDir=$(pwd)

PARAMS="$1"
function containsElement () {
  # If file change was passed down, this function always returns true
  if [ -n "$PARAMS" ]; then
    return 0;
  fi
  local e match="$1"
  shift
  for e; do [[ "$e" =~ ^$match ]] && return 0; done
  return 1
}

function build {
    dir=$1
    EXIT_CODE=0
    # We don't want to fail the whole script if contains fails
    containsElement $dir "${FILES[@]}" || EXIT_CODE=$?
    if [ $EXIT_CODE -eq 0 ]; then
        echo Building $dir
    else
        echo Skipping $dir
        return 0
    fi
    
    cd $dir

    # build the proto buffers
    #find . -name "*.proto" | xargs --no-run-if-empty protoc --proto_path=. --micro_out=. --go_out=.  

    # build the binaries
    go build -ldflags="-s -w" -o service .
    cp $rootDir/dumb-init/dumb-init dumb-init

    # build the docker image
    tag=docker.pkg.github.com/micro/services/$(echo $dir | tr / -)
    docker build . -t $tag -f $rootDir/.github/workflows/Dockerfile

    if [ -n "$1" ] && [ "$BRANCH" = "refs/heads/master" ]; then
        # push the docker image
        echo Pushing $tag
        docker push $tag
    else
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
