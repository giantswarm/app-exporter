FROM alpine:3.14.3

RUN apk add --no-cache ca-certificates

ADD ./app-exporter /app-exporter

ENTRYPOINT ["/app-exporter"]
