FROM scratch
ADD /bin
ENTRYPOINT ["dumb-init"]
