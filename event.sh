#!/bin/bash

# event is build.{started, failed, finished}
EVENT=$1
# the url to send events to
URL=https://micro.mu/platform/v1/github/events

curl $URL -X POST -d @$HOME/changes.json \
-H "Content-Type: application/json" \
-H "Micro-Event: $EVENT"