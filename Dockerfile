FROM scratch
ADD binaries.tar.gz /bin
ENTRYPOINT ["dumb-init"]
