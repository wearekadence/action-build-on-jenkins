# syntax=docker/dockerfile:1

FROM golang:1.21-bookworm AS builder

WORKDIR /go/src/app
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/app

FROM golang:1.21-bookworm

COPY --from=builder /go/bin/app /
CMD ["/app"]