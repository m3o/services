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

## S3 bucket versioning and lifecycle 
For added peace of mind, you should probably enable versioning on the S3 bucket to ensure you don't lose data due to an errant upload. However, if you do this you should add some lifecycle rules to prevent your storage usage exploding. See [lifecycle.xml](lifecycle.xml) for an example policy which expires (deletes) all non current backup files after 1 day. You can apply the polciy with the below:

`s3cmd -c s3cfg setlifecycle lifecycle.xml s3://<BUCKET_NAME>`