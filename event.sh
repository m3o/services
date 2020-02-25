#!/bin/bash

# event is build.{started, failed, finished}
EVENT=$1
# build is the build number to be used at $REPO/actions/runs/$BUILD
BUILD=$2
# commit is a git hash to be used as $REPO/commit/$COMMIT
COMMIT=$3
# repo is a full github.com/micro/services url
REPO=$4
# the url to send events to
URL=https://micro.mu/platform/v1/github/events

curl $URL -X POST -d @$HOME/services.json \
-H "Content-Type: application/json" \
-H "Micro-Event: $EVENT" \
-H "X-Github-Build: $BUILD" \
-H "X-Github-Commit: $COMMIT" \
-H "X-Github-Repo: $REPO"
