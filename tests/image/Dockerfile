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
RUN cd alert && go get 
RUN cd billing && go get 
RUN cd builder && go get 
RUN cd customers && go get
RUN cd emails && go get
RUN cd invite && go get
RUN cd namespaces && go get
RUN cd notifications && go get
RUN cd payments/provider/stripe && go get
RUN cd platform && go get
RUN cd signup && go get
RUN cd status && go get
RUN cd subscriptions && go get
COPY ./micro/micro /microserver
ENTRYPOINT ["sh", "/bin/run.sh"]
