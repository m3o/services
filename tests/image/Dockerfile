FROM alpine:latest
RUN apk add make git go gcc libtool musl-dev curl bash

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

RUN apk add ca-certificates && \
    rm -rf /var/cache/apk/* /tmp/* && \
    [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf

RUN         apk add --update ca-certificates openssl tar && \
            wget https://github.com/coreos/etcd/releases/download/v3.4.7/etcd-v3.4.7-linux-amd64.tar.gz && \
            tar xzvf etcd-v3.4.7-linux-amd64.tar.gz && \
            mv etcd-v3.4.7-linux-amd64/etcd* /bin/ && \
            apk del --purge tar openssl && \
            rm -Rf etcd-v3.4.7-linux-amd64* /var/cache/apk/*
VOLUME      /data
EXPOSE      2379 2380 4001 7001
ADD         services/tests/image/run-etcd.sh /bin/run.sh

ENV MICRO_REGISTRY=etcd

# Speeding up tests by predownloading and building dependencies for services used.
COPY . .
RUN pushd alert && go get && popd
RUN pushd billing && go get && popd
RUN pushd builder && go get && popd
RUN pushd customers && go get && popd
RUN pushd emails && go get && popd
RUN pushd invite && go get && popd
RUN pushd namespaces && go get && popd
RUN pushd notifications && go get && popd
RUN pushd payments/provider/stripe && go get && popd
RUN pushd platform && go get && popd
RUN pushd signup && go get && popd
RUN pushd status && go get && popd
RUN pushd subscriptions && go get && popd
COPY ./micro/micro /microserver
ENTRYPOINT ["sh", "/bin/run.sh"]
