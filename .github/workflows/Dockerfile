FROM alpine:latest
COPY ./dumb-init /bin/dumb-init
COPY . .
ENTRYPOINT ["dumb-init", "./service"]
