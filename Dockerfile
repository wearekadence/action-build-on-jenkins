# syntax=docker/dockerfile:1@sha256:87999aa3d42bdc6bea60565083ee17e86d1f3339802f543c0d03998580f9cb89

FROM golang:1.21-bookworm@sha256:c6a5b9308b3f3095e8fde83c8bf4d68bd101fce606c1a0a1394522542509dda9 AS builder

WORKDIR /go/src/app
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/app

FROM golang:1.21-bookworm@sha256:c6a5b9308b3f3095e8fde83c8bf4d68bd101fce606c1a0a1394522542509dda9

COPY --from=builder /go/bin/app /
CMD ["/app"]