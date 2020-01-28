FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates

FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY stock-colors /bin/app
CMD ["/bin/app"]