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

# Speeding up tests by predownloading and building dependencies for services used.
RUN mkdir services
COPY services/go.* services/
RUN cd services && go mod download
COPY . .
RUN go install services/alert/main.go
RUN go install services/billing/main.go
RUN go install services/builder/main.go
RUN go install services/customers/main.go
RUN go install services/emails/main.go
RUN go install services/invite/main.go
RUN go install services/namespaces/main.go
RUN go install services/payments/provider/stripe/main.go
RUN go install services/platform/main.go
RUN go install services/signup/main.go
RUN go install services/subscriptions/main.go
COPY ./micro/micro /microserver
ENTRYPOINT ["sh", "/bin/run.sh"]
