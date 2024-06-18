# docker project learning

> date: 2024.6.18

## 一些补充

发现了之前写的dockerfile有些问题，在定义的时候尽量不要让docker去拿上一级目录的文件，运行docker build命令在根目录中。

**修改后的Dockerfile**

```go
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

```

## build image 

in root dir, run the following command, 

**note: naming template**

`docker build -t jfrog.wosai-inc/greeting_server:v0.0.1 -f server/Dockerfile .`

## list image

```bash
➜  greeting git:(main) ✗ docker image ls -a
REPOSITORY                            TAG       IMAGE ID       CREATED          SIZE
jfrog.wosai-inc.com/greeting_client   v0.0.1    10196d0854ef   10 minutes ago   19.7MB
jfrog.wosai-inc.com/greeting_server   v0.0.1    2005cb0ca364   2 hours ago      67.6MB
```

## how to make containers communicate 

using `docker network` 

create a network 

`docker network create my_network`

list the network

```bash
➜  greeting git:(main) ✗ docker network ls
NETWORK ID     NAME               DRIVER    SCOPE
e52a1f553de6   bridge             bridge    local
280eee1f35ef   greeting_default   bridge    local
d6bbf19c55e2   host               host      local
27c31a014e15   minikube           bridge    local
f79f52d623f6   none               null      local
➜  greeting git:(main) ✗ docker network create my_network
2b45b019f0d2eba483b38fc5c985953dc6764aef70d72d078a16154fe898ee94
➜  greeting git:(main) ✗ docker network ls               
NETWORK ID     NAME               DRIVER    SCOPE
e52a1f553de6   bridge             bridge    local
280eee1f35ef   greeting_default   bridge    local
d6bbf19c55e2   host               host      local
27c31a014e15   minikube           bridge    local
2b45b019f0d2   my_network         bridge    local
f79f52d623f6   none               null      local
```

run the server container and **connect the network**

```bash
➜  greeting git:(main) ✗ docker run --network my_network --name greeting_server jfrog.wosai-inc.com/greeting_server:v0.0.1

2024/06/18 06:24:15 server listening at [::]:50051

```
**edit `main.go` under client!!!**

```golang

const (
	address     = "greeting_server:50051"
	defaultName = "world"
)
```

rebuild the client image

run the client container

```bash
➜  greeting git:(main) ✗ docker run --network my_network --name greeting_client jfrog.wosai-inc.com/greeting_client:v0.0.1

2024/06/18 06:25:30 Greeting: Hello xzl
```

