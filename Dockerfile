# syntax=docker/dockerfile:1

FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY **/*.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/build-in-jenkins

FROM scratch

COPY --from=builder /app/build-in-jenkins /app/build-in-jenkins

CMD ["/app/build-in-jenkins"]