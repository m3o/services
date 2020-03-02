FROM alpine:latest
COPY ./dumb-init /bin/dumb-init
COPY . .
RUN ["chmod", "+x", "./app"]
ENTRYPOINT ["dumb-init", "./app"]