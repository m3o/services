#!/bin/bash
set -e
set -x

go env

go get --tags extended github.com/gohugoio/hugo
go get github.com/micro/platform

mkdir html

cd docuapi/microApi/content;
platform doc-gen --path=../
ls
cd ..

hugo -D 
mv public/* ../../html/