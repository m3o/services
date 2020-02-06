#!/bin/sh

set -x
set -e

SOURCE=$1
REPO=github.com/micro/services

# clone the repo
echo "Downloading $SOURCE"
git clone --no-checkout https://$REPO

# cd into source
cd services

# make a sparse checkout
git sparse-checkout init --cone

# set the repo to checkout
git sparse-checkout set $SOURCE

# go to source
cd $SOURCE

# run the source
go run .
