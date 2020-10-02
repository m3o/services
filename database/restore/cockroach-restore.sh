#!/bin/bash
# setup certs
cd ~
mkdir crts
cp /certs/* crts
mv crts/cert.pem crts/client.root.crt
mv crts/key.pem crts/client.root.key

# Get the latest dump file. The sort will sort on the date since s3cmd output looks like below
# 2020-07-28 10:08     33451   s3://staging-db-backups/dumps.tar.gz
latest_file=$(s3cmd ls s3://$S3_BUCKET/ -c /s3-config/s3cfg | sort | tail -n1 | awk '{print $4}')
dump_file_name=$(basename $latest_file)
s3cmd get $latest_file $file_name -c /s3-config/s3cfg
tar -xzf $dump_file_name
while read -r line
do
  # pull db name from name of dump file. dump of db foobar is named foobar-dump.sql
  db=$(echo $line | sed 's|dumps/\(.*\)-dump.sql|\1|g')
  cockroach sql --execute "DROP DATABASE IF EXISTS $db; CREATE DATABASE $db" --certs-dir=crts --host cockroachdb-cluster
  cockroach sql --database $db --certs-dir=crts --host cockroachdb-cluster < $line
done < <(ls dumps/*.sql)
