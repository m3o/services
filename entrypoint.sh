#!/bin/dumb-init /bin/sh

set -x
set -e

SOURCE=$1
REPO=github.com/micro/services

# clone the repo
echo "Cloning $REPO"
git clone --no-checkout https://$REPO

# cd into source
cd services

# make a sparse checkout
git sparse-checkout init --cone

# set the repo to checkout
echo "Checking out $SOURCE"
git sparse-checkout set $SOURCE

# go to source
cd $SOURCE

# run the source
echo "Running service"
go run .
