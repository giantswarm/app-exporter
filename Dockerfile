FROM gsoci.azurecr.io/giantswarm/alpine:3.24.2

RUN apk add --no-cache ca-certificates

ARG TARGETARCH
COPY app-exporter-linux-${TARGETARCH} /app-exporter

ENTRYPOINT ["/app-exporter"]
