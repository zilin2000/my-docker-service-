# update: 6.6 -> multiple stage

# stage1: build Go
FROM golang:1.16 AS builder 

# make a new working directory called /build 
WORKDIR /build

# copy needed files into image
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY server/main.go .
COPY proto ./proto

# go build: name called server
RUN go build -o server main.go

# stage 2: 把build好的东西都放到/app里面去运行
FROM ubuntu:18.04

WORKDIR /app  

COPY --from=builder /build/server .
COPY --from=builder /build/proto ./proto

CMD [ "./server" ]

 