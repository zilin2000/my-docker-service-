# stage 1: build Go
FROM golang:1.22.3 AS builder

WORKDIR /build

COPY go.mod .
COPY v1/server1.go .

RUN go build -o server1 server1.go


# stage 2: run ./server
FROM ubuntu:22.04

WORKDIR /app

COPY --from=builder /build/server1 .

CMD [ "./server1" ]
