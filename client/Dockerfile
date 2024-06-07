# same as Dockerfile in server
FROM golang:1.16 AS builder 

WORKDIR /build
COPY ../go.mod .
COPY ../go.sum .

RUN go mod download

COPY ../client/main.go .
COPY ../proto ./proto


# CMD ["go", "run", "main.go"]
RUN go build -o client main.go


# stage 2
FROM alpine:latest

WORKDIR /app

COPY --from=builder ../build/client .
COPY --from=builder ../build/proto ./proto

CMD [ "./client","xzl" ]