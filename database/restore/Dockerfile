FROM cockroachdb/cockroach
RUN apt update && apt install s3cmd -y
COPY cockroach-restore.sh /cockroach-restore.sh
ENTRYPOINT ["/cockroach-restore.sh"]