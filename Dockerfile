# syntax=docker/dockerfile:1@sha256:4edf897a3ffa55b89f906fc8cc78afdb3f1834cc9c7083565e611a8a7d5fe99e

FROM golang:1.27-bookworm@sha256:648f440f42a0958804efb24df176f806f9d353b41f1c0627f666428e40310f6b AS builder

WORKDIR /go/src/app
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/app

FROM golang:1.27-bookworm@sha256:648f440f42a0958804efb24df176f806f9d353b41f1c0627f666428e40310f6b

COPY --from=builder /go/bin/app /
CMD ["/app"]