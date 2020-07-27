# Cockroach backup

This directory defines a Kubernetes cron job to periodically backup the Cockroach cluster. 

The cron job runs the backup script at regular intervals which:
1. runs `cockroach dump` on the DBs
2. compresses the dumps
3. uploads to S3

## Setup
The cron job depends on:
- a Kubernetes secret `s3-config` which is an `s3cmd` config file. Example 
```
[default]
# Object Storage Region FR-PAR
host_base = s3.fr-par.scw.cloud
host_bucket = %(bucket)s.s3.fr-par.scw.cloud
bucket_location = fr-par
use_https = True

# Login credentials
access_key = <ACCESS_KEY>
secret_key = <SECRET>
```
Put this in a file `s3cfg` and run `kubectl create secret generic s3-config --from-file=./s3cfg`
- A config map `s3-config` that defines the key `bucket` which tells the script which bucket to upload to. `kubectl create configmap s3-config --from-literal=bucket=<BUCKET_NAME>`

