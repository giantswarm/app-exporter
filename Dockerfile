FROM alpine:3.10

RUN apk add --no-cache ca-certificates

ADD ./app-collector /app-collector

ENTRYPOINT ["/app-collector"]
