# Cockroach restore

This directory defines a Kubernetes pod to import a dump from s3 in to the Cockroach cluster. 

The pod will:
1. download the latest dump file from a given S3 bucket
2. decompress the dumps
3. import to the cockroach cluster `cockroach sql --execute="IMPORT PGDUMP '<dump file>'"` 

## Setup
The pod depends on:
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
- A config map `s3-config` that defines the key `bucket` which tells the script which bucket to download from. `kubectl create configmap s3-config --from-literal=bucket=<BUCKET_NAME>`

