#!/bin/bash

PWD=`pwd`

cd web/content && bundle exec jekyll build -d ../html && cd $PWD
