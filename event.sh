#!/bin/bash

# event is build.{started, failed, finished}
EVENT=$1
# build is the build number to be used at $REPO/actions/runs/$BUILD
BUILD=$2
# the url to send events to
LIVE=https://web.micro.mu/platform/v1/github/events
# the staging url to send events to
STAGING=https://web.m3o.dev/platform/v1/github/events

sendEvent() {
	URL=$1

	echo "sending event $EVENT to $URL"
	curl --connect-timeout 5 --retry 3 -s -S \
	$URL -X POST -d @$HOME/changes.json \
	-H "Content-Type: application/json" \
	-H "X-Github-Build: $BUILD" \
	-H "Micro-Event: $EVENT"
}

## send event to live
## TODO: only on releases
sendEvent $LIVE
## sent event to staging
sendEvent $STAGING
