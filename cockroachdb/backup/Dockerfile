FROM cockroachdb/cockroach
RUN apt update && apt install s3cmd -y
COPY cockroach-backup.sh /cockroach-backup.sh
ENTRYPOINT ["/cockroach-backup.sh"]