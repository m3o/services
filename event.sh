#!/bin/bash

EVENT=$1
COMMIT=$2
REPO=$3
URL=https://micro.mu/platform/v1/github/events

curl $URL -X POST -d @$HOME/files.json \
-H "Content-Type: application/json" \
-H "Micro-Event: $EVENT" \
-H "X-Github-Sha: $COMMIT" \
-H "X-Github-Repo: $REPO"
