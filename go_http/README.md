# go http server


## create `.go` files
```plaintext
GO_HTTP_Project/
│
├── v1/
│   └── server1.go
│
├── v2/
│   └── server2.go
│
└── go.mod
```

after running two `go` files, run following commands

```sh
➜  greeting git:(main) curl localhost:8090/v2
hello, this is v2
➜  greeting git:(main) ✗ curl localhost:8080/v1
hello, this is v1
```

## build docker containers

### edit `Dockerfile`

```go
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
```

### build image and run containers

```sh
docker build -t jfrog.wosai-inc.com/go_httpserver2:v0.0.1 -f v
2/Dockerfile .

docker run -p 8090:8090 jfrog.wosai-inc.com/go_httpserver2:v0.0.1
start listening on port 8090...
```

access two servers using `curl`

```sh
➜  v2 git:(main) ✗ curl localhost:8080/v1
hello, this is v1
➜  v2 git:(main) ✗ curl localhost:8090/v2                                       
hello, this is v2
```




# update on 6.19 afternoon

我们需要用两个ip代表两个容器，且他们共享一个端口号。

**edit code**

```golang
package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello, this is v1\n")
}

func main() {
	port := 8080
	fmt.Printf("start listening on port %d...\n", port)
	http.HandleFunc("/v1", hello)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
```

```golang
package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello, this is v2\n")
}

func main() {
	port := 8080
	fmt.Printf("start listening on port %d...\n", port)
	http.HandleFunc("/v1", hello)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
```

就其实两个服务器唯一不一样的点在于他们输出的内容不一样.

**build two images and run containers**

**get two container ip address**

```sh
➜  v2 git:(main) ✗ docker ps
CONTAINER ID   IMAGE                                      COMMAND       CREATED          STATUS          PORTS     NAMES
95138b5f3d20   jfrog.wosai-inc.com/go_httpserver:v0.0.2   "./server2"   32 seconds ago   Up 32 seconds             vigorous_kepler
d99bfa9da899   jfrog.wosai-inc.com/go_httpserver:v0.0.1   "./server1"   42 seconds ago   Up 42 seconds             trusting_wright

➜  v2 git:(main) ✗ docker inspect \
  -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' d99
172.17.0.2
➜  v2 git:(main) ✗ docker inspect \
  -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' 951
172.17.0.3
```

**edit nginx.conf**

```go
events {}

http {
    upstream backend {
        server 172.17.0.2:8080 weight=9;
        server 172.17.0.3:8080 weight=1;
    }

    server {
        listen 80;

        location / {
            proxy_pass http://backend;
        }
    }
}
```

**run nginx container by using 挂载本地conf**

```sh
docker run --rm -v /Users/xuzilin/Desktop/sqbLearning/go_http/nginx/nginx.conf:/etc/nginx/nginx.conf nginx:latest
```

**进入nginx容器**

```sh
➜  go_http git:(main) ✗ docker ps
CONTAINER ID   IMAGE                                      COMMAND                   CREATED          STATUS          PORTS     NAMES
25494b0271a9   nginx:latest                               "/docker-entrypoint.…"   11 seconds ago   Up 11 seconds   80/tcp    compassionate_boyd
4eebc01ad20d   jfrog.wosai-inc.com/go_httpserver:v0.0.2   "./server2"               6 minutes ago    Up 6 minutes              mystifying_mestorf
198b1699a02d   jfrog.wosai-inc.com/go_httpserver:v0.0.1   "./server1"               6 minutes ago    Up 6 minutes              nervous_bouman
➜  go_http git:(main) ✗ docker exec -it 254 bash
```
然后我们就成功实现了用nginx容器做负载均衡啦！

```sh

➜  go_http git:(main) ✗ docker exec -it 254 bash
root@25494b0271a9:/# curl http://localhostL8080/v1
curl: (6) Could not resolve host: localhostL8080
root@25494b0271a9:/# curl http://localhost:80/v1
hello, this is v1
root@25494b0271a9:/# curl http://localhost:80/v1
hello, this is v1
root@25494b0271a9:/# curl http://localhost:80/v1
hello, this is v1
root@25494b0271a9:/# curl http://localhost:80/v1
hello, this is v1
root@25494b0271a9:/# curl http://localhost:80/v1
hello, this is v1
root@25494b0271a9:/# curl http://localhost:80/v1
hello, this is v2
root@25494b0271a9:/# curl http://localhost:80/v1
hello, this is v1
root@25494b0271a9:/# curl http://localhost:80/v1
hello, this is v1
```