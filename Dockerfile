FROM golang:1.16-alpine3.14 as builder
LABEL maintainer="iKOPKACtraxa"
WORKDIR /app
COPY . .
ARG VERSION_HASH="somehash"
RUN CGO_ENABLED=0 go build \
    -ldflags "-X main.release=develop -X main.buildDate=$(date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$VERSION_HASH" \
    -o /bin/rotation \
    ./cmd/rotation

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache bash
WORKDIR /app
COPY . .
COPY --from=builder /bin/rotation /app/bin/rotation
EXPOSE 50051:50051
CMD ["/app/bin/rotation", "--config", "/app/configs/configDocker.yaml"]