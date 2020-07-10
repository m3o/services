#!/bin/bash

# event is build.{started, failed, finished}
EVENT=$1
# build is the build number to be used at $REPO/actions/runs/$BUILD
BUILD=$2
# the url to send events to
LIVE=https://api.m3o.com/v1/platform/events
# the staging url to send events to
STAGING=https://api.m3o.dev/v1/platform/events

sendEvent() {
	URL=$1

	echo "sending event $EVENT to $URL"
	curl --connect-timeout 5 --retry 3 -s -S \
	$URL -X POST -d @$HOME/changes.json \
	-H "Content-Type: application/json" \
	-H "X-Github-Build: $BUILD" \
	-H "Micro-Event: $EVENT"
}

## TODO: uncomment when we revive auto updates
## send event to live
#sendEvent $LIVE
## sent event to staging
#sendEvent $STAGING
