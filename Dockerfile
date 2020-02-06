FROM micro/go-micro
RUN sed -i -e 's/v[[:digit:]]\..*\//edge\//g' /etc/apk/repositories
RUN apk upgrade --update git
COPY entrypoint.sh /
WORKDIR /
RUN chmod 755 entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]
