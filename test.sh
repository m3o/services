#!/bin/bash

# Ensure monorepo conforms to our requirements
#
# All services should include a README.md
# All services should include a main.go
# All services should include a go.mod and go.sum
# No binaries should be present anywhere

set -e

function fatal {
	PWD=`pwd`

	echo $@ && exit 1
		
}

function checkFiles() {
	# should exist
	for file in README.md go.mod go.sum main.go; do
		if [ ! -f $file ]; then
			fatal "$1 does not include $file"
		fi
	done

	# should NOT exist e.g binaries
	for file in ${1} ${1}-srv; do
		if [ -f $file ]; then
			fatal "$1 should not include $file"
		fi
	done
}


ls | while read service; do
	if [ ! -d $service ]; then
		continue
	fi

	echo "Checking $service"
	# push into service repo
	pushd $service &> /dev/null
	# check the files
	checkFiles $service
	# exit the service
	popd &>/dev/null
done
