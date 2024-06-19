FROM  golang:1.22.3 AS builder

WORKDIR /build

# copy all files we need
COPY go.mod .
COPY v2/server2.go .

# build it
RUN go build -o server2 server2.go

# stage2
FROM ubuntu:22.04

WORKDIR /app

# copy server2
COPY --from=builder /build/server2 .

CMD [ "./server2" ]
