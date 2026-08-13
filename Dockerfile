# syntax=docker/dockerfile:1@sha256:ecfaec9ed6d810b56388c508f4121597bfbba70d41a6dfeee4d8cad5f295fc32

FROM golang:1.21-bookworm@sha256:c6a5b9308b3f3095e8fde83c8bf4d68bd101fce606c1a0a1394522542509dda9 AS builder

WORKDIR /go/src/app
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/app

FROM golang:1.21-bookworm@sha256:c6a5b9308b3f3095e8fde83c8bf4d68bd101fce606c1a0a1394522542509dda9

COPY --from=builder /go/bin/app /
CMD ["/app"]