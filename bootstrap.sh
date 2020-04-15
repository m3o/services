# This script assumes configs are already set up.
# Ask for config values as they contain secret, hence can't be part
# of this file.

set -x

micro server &
MICRO_ID=$!

sleep 5

micro run --server platform/api
micro run --server platform/web
micro run --server platform/service
micro run --server users/service
micro run --server account/api 
micro run --server account/web

# Waiting for account & user servicec to start
sleep 5
micro call go.micro.api.account Account.Signup '{"email":"user@micro.mu","password":"local"}'

function cleanup()
{
    echo "Killing micro server"
    kill $MICRO_ID
}

trap cleanup EXIT