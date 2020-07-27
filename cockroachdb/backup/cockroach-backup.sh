#!/bin/bash
# setup certs
cd ~
mkdir crts
cp /certs/* crts
mv crts/cert.pem crts/client.root.crt
mv crts/key.pem crts/client.root.key

# dump the data
mkdir dumps
while read -r line
do
  if [[ $line =~ (database_name|postgres|system|defaultdb) ]]; then
    continue
  fi
  if [[ $line =~ "^--" ]]; then
    continue
  fi
  if [[ $line =~ "^\(" ]]; then
    break
  fi
  cockroach dump $line --certs-dir=crts --host cockroachdb-cluster > dumps/$line-dump.sql

done < <(echo "show databases;" | cockroach sql --certs-dir=crts --host cockroachdb-cluster --format tsv)

# compress
tar -czf dumps.tar.gz dumps

# send to s3
s3cmd put -c /s3-config/s3cfg dumps.tar.gz s3://$S3_BUCKET/