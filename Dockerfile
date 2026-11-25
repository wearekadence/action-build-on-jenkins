# syntax=docker/dockerfile:1

FROM golang:1.27-bookworm@sha256:ce59b41415e2524ae99926fe1f37ce48ba8c194dcbd2a56a65436cdcefa2bd34 AS builder

WORKDIR /go/src/app
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/app

FROM golang:1.27-bookworm@sha256:ce59b41415e2524ae99926fe1f37ce48ba8c194dcbd2a56a65436cdcefa2bd34

COPY --from=builder /go/bin/app /
CMD ["/app"]