FROM alpine:latest as dumb-init
RUN apk add --no-cache build-base git bash
RUN git clone https://github.com/Yelp/dumb-init.git
WORKDIR /dumb-init
RUN make

FROM micro/go-micro
COPY --from=dumb-init /dumb-init/dumb-init /bin/dumb-init
RUN sed -i -e 's/v[[:digit:]]\..*\//edge\//g' /etc/apk/repositories
RUN apk upgrade --update git
COPY entrypoint.sh /
WORKDIR /
RUN chmod 755 entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]
