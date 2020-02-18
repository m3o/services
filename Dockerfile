# Dumb Init
FROM alpine:latest as dumb-init
RUN apk add --no-cache build-base git bash
RUN git clone https://github.com/Yelp/dumb-init.git
WORKDIR /dumb-init
RUN make

# Copy the services
FROM golang:1.13-alpine as builder
COPY . /services
WORKDIR services
RUN mkdir builds

# Build the services
RUN for dir in */; do \
cd $dir && \
(CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build . && mv ./${dir::-1} ../builds/${dir::-1} || true) && \
cd ..; \
done

# Copy result to new image 
FROM scratch
COPY --from=dumb-init /dumb-init/dumb-init /bin/dumb-init
COPY --from=builder services/builds /bin