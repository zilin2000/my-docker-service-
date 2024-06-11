# update: 6.6 -> multiple stage

# stage1: build Go
FROM golang:1.16 AS builder 

# make a new working directory called /build 
WORKDIR /build

# 做上一次一样的操作
COPY ../go.mod .
COPY ../go.sum .
RUN go mod download

COPY ../server/main.go .
 
COPY ../proto ./proto


# go build: name called server
RUN go build -o server main.go

# stage 2: 把build好的东西都放到/app里面去运行
FROM alpine:latest

WORKDIR /app  

COPY --from=builder /build/server .
COPY --from=builder /build/proto ./proto

CMD [ "./server" ]

 