FROM gsoci.azurecr.io/giantswarm/alpine:3.21.3

RUN apk add --no-cache ca-certificates

ADD ./app-exporter /app-exporter

ENTRYPOINT ["/app-exporter"]
