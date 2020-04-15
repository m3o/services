#!/bin/bash
set -e
set -x

export GOBIN=$(go env GOPATH)/bin
export PATH=$PATH:$GOBIN

go env

go get --tags extended github.com/gohugoio/hugo
go get github.com/micro/micro/cmd/platform

mkdir html
mkdir -p docuapi/microApi/content

cd docuapi/microApi/content
platform doc-gen --path=../
cd ..

hugo -D 
mv public/* ../../html/
ls ../../html/
