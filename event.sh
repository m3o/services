#!/bin/bash

# event is build.{started, failed, finished}
EVENT=$1
# build is the build number to be used at $REPO/actions/runs/$BUILD
BUILD=$2
# the url to send events to
# URL=https://micro.mu/platform/v1/github/events
URL=https://b24e2441.ngrok.io/platform/v1/github/events

curl $URL -X POST -d @$HOME/changes.json \
-H "Content-Type: application/json" \
-H "X-Github-Build: $BUILD" \
-H "Micro-Event: $EVENT"