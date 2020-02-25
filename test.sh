#!/bin/bash

# Ensure monorepo conforms to our requirements
#
# All services should include a README.md
# All services should include a main.go
# All services should include a go.mod
# No binaries should be present anywhere

set -e

function fatal {
	PWD=`pwd`

	echo $@ && exit 1
		
}

function checkFiles() {
	# should exist
	for file in README.md go.mod main.go; do
		if [ ! -f $file ]; then
			fatal "$1 does not include $file"
		fi
	done

	# should NOT exist e.g binaries
	for file in ${1} ${1}-srv; do
		if [ -f $file ] && ! grep -Fxq $file .gitignore; then
			fatal "$1 should not include $file"
		fi
	done
}

# Check the service repos
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

# Check for outlier binaries
find . -type f -size +1M | grep -v \.git | while read file; do
	if ! grep $(basename -- $file) $(dirname -- $file)/.gitignore; then
		fatal "$file is larger than 1M"
	fi
done
