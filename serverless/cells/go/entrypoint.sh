#!/bin/dumb-init /bin/sh

set -x  
set -e

function join_by { local IFS="$1"; shift; echo "$*"; }

# In the form of: github.com/crufter/reponame/folder1/folder2/folder2/folder3
P=$1
arrPATH=(${P//\// })
DIR=${arrPATH[2]}
SOURCE=${arrPATH[@]:3}
SOURCE=${SOURCE// /\/}
REPO=${arrPATH[@]: 0:3}
REPO=${REPO// /\/}

echo $REPO $SOURCE

# clone the repo
echo "Cloning $REPO"
git clone https://$REPO

cd $DIR
# go to source
cd $SOURCE

# run the source
echo "Running service"
go run .