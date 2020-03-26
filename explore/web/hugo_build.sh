#!/bin/bash
set -e
set -x

go get --tags extended https://github.com/gohugoio/hugo
go get github.com/micro/platform

cd docuapi/microApi/content;
platform doc-gen --path=../
cd ..

hugo -D 
mv public/* ../../html/