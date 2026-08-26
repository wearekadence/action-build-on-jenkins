# syntax=docker/dockerfile:1

FROM golang:1.27-bookworm@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452 AS builder

WORKDIR /go/src/app
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/app

FROM golang:1.27-bookworm@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452

COPY --from=builder /go/bin/app /
CMD ["/app"]