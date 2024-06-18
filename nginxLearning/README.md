# nginx learning

> date: 2024.6.13  
> author: zilin xu  
> resouce: 狂神说


## intro

假设我们有很多client和server，我们需要一个中间人做`反向代理`,像是一个虚拟的服务器接受和转发我们client的请求并发送给对应的server；我们还需要通过这个中间人去分配服务器的权重，比如说server a可以搭载64G。 server b可以搭载6G，我们当然希望更多的请求可以打到server a。所以我们需要中间人去做`负载均衡`。

This is **nginx**

这世上所有的架构问题：没有什么是加一层解决不了的问题。

## 正向代理vs反向代理

正向代理： 比如VPN，他是代理客户端的请求。

反向代理： 代理我们的服务器的，比如baidu在上海深圳北京都有服务器，但你访问永远是www.baidu.com.

## 负载均衡

轮询：依次循环  

加权轮询

iphash： 

## 动静分离

有些静态的请求直接通过nginx返回

## install nginx and basic commands

```bash
➜  nginxLearning git:(main) ✗ brew install nginx
==> Downloading https://ghcr.io/v2/homebrew/core/nginx/manifests/1.27.0
############################################################################### 100.0%
==> Fetching nginx
==> Downloading https://ghcr.io/v2/homebrew/core/nginx/blobs/sha256:64223e300749be61e1
############################################################################### 100.0%
==> Pouring nginx--1.27.0.arm64_sonoma.bottle.tar.gz
==> Caveats
Docroot is: /opt/homebrew/var/www

The default port has been set in /opt/homebrew/etc/nginx/nginx.conf to 8080 so that
nginx can run without sudo.

nginx will load all files in /opt/homebrew/etc/nginx/servers/.

To start nginx now and restart at login:
  brew services start nginx
Or, if you don't want/need a background service you can just run:
  /opt/homebrew/opt/nginx/bin/nginx -g daemon\ off\;
==> Summary
🍺  /opt/homebrew/Cellar/nginx/1.27.0: 27 files, 2.4MB
==> Running `brew cleanup nginx`...
Disable this behaviour by setting HOMEBREW_NO_INSTALL_CLEANUP.
Hide these hints with HOMEBREW_NO_ENV_HINTS (see `man brew`).
```

### start and stop the nginx service

```bash
brew services start nginx
## then you can visit http://localhost:8080
brew services stop nginx
```

**nginx path on my Mac**

`/opt/homebrew/etc/nginx/`

```bash
cd /opt/homebrew/etc/nginx/
code .
```

### edit `.conf` file and reload it

change the listen port into `8081`

run `brew services reload nginx`

## implement a simple 负载均衡

implement two simple server using `JS`

```bash
mkdir serverJS
cd serverJS
touch server1.js server2.js
```

edit two `.js` file

```javascript
const http = require('http');

const server = http.createServer((req, res) => {
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('Response from Server 1\n');
});

server.listen(3001, () => {
    console.log('Server 1 is running on http://localhost:3001');
});

const http = require('http');

const server = http.createServer((req, res) => {
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('Response from Server 2\n');
});

server.listen(3002, () => {
    console.log('Server 2 is running on http://localhost:3002');
});
```

start two servers

```bash
node server1.js 
node server2.js 
```

edit `nginx.conf`

```go
upstream backend {
        server localhost:3001;
        server localhost:3002;
    }

    server {
        listen       8080;
        server_name  localhost;

        location / {
            proxy_pass http://backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
```

reload and start nginx


```bash
brew services reload nginx
brew services restart nginx
```


## implement data control 9:1 

add `weight` in upstream

```golang
    upstream backend {
        server localhost:3001 weight=9;
        server localhost:3002 weight=1;
    }
```


好的，我很高兴能帮助你理解 `nginx.conf` 文件的语法和配置原理。

### Nginx 配置文件结构

Nginx 配置文件主要由以下几部分组成：

1. **全局块（Global block）**：包括影响整个 Nginx 服务器的指令，如用户权限、进程数量等。
2. **事件块（Events block）**：包括处理连接的指令，如最大连接数等。
3. **HTTP 块（HTTP block）**：包括处理 HTTP 请求的指令，是最重要的一部分，包含多个子块（如 `server` 块、`location` 块等）。

### 示例配置文件分析

```nginx
#user  nobody;
worker_processes  1;

#error_log  logs/error.log;
#error_log  logs/error.log  notice;
#error_log  logs/error.log  info;

#pid        logs/nginx.pid;

events {
    worker_connections  1024;
}

http {
    include       mime.types;
    default_type  application/octet-stream;

    #log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
    #                  '$status $body_bytes_sent "$http_referer" '
    #                  '"$http_user_agent" "$http_x_forwarded_for"';

    #access_log  logs/access.log  main;

    sendfile        on;
    #tcp_nopush     on;

    #keepalive_timeout  0;
    keepalive_timeout  65;

    #gzip  on;

    upstream backend {
        server localhost:3001 weight=9;
        server localhost:3002 weight=1;
    }

    server {
        listen       8080;
        server_name  localhost;

        location / {
            proxy_pass http://backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
        }
    }

    include servers/*;
}
```

### 主要配置项讲解

#### 全局块

```nginx
worker_processes  1;
```

- `worker_processes`：指定 Nginx 工作进程的数量。一般设置为 CPU 核心数以充分利用多核。

#### 事件块

```nginx
events {
    worker_connections  1024;
}
```

- `worker_connections`：每个工作进程的最大连接数。这个值决定了 Nginx 可以处理的最大并发连接数。

#### HTTP 块

```nginx
http {
    include       mime.types;
    default_type  application/octet-stream;
```

- `include mime.types`：包含文件，定义了 MIME 类型。
- `default_type`：设置默认的 MIME 类型。

```nginx
    upstream backend {
        server localhost:3001 weight=9;
        server localhost:3002 weight=1;
    }
```

- `upstream`：定义一组后端服务器，用于负载均衡。
  - `server`：定义后端服务器地址和权重（`weight`）。权重越大，服务器分配到的请求越多。

```nginx
    server {
        listen       8080;
        server_name  localhost;
```

- `server` 块：定义一个虚拟服务器。
  - `listen`：定义服务器监听的端口。
  - `server_name`：定义服务器名称。

```nginx
        location / {
            proxy_pass http://backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
```

- `location` 块：定义请求的路由规则。
  - `proxy_pass`：将请求转发到 `upstream` 定义的后端服务器。
  - `proxy_set_header`：设置转发请求时的头信息。

```nginx
        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
        }
```

- `error_page`：定义错误页面。
  - `location = /50x.html`：定义特定错误页面的位置。

### 总结

通过配置 `nginx.conf` 文件，我们可以灵活地控制流量分发、负载均衡以及错误处理。理解每个配置项的作用，有助于我们根据需求进行定制和优化。

### update on 6.18

我们需要升级一下我们的项目，用官方nginx镜像并挂载本地的配置文件（.conf）

> nginx PATH on my mac : /opt/homebrew/etc/nginx  
> JS services on my mac: /Users/xuzilin/Desktop/sqbLearning/nginxLearning/serverJS

**steps**

pull official `nginx image`

`docker pull nginx`

edit `.conf` file, Docker 容器中的 Nginx 无法连接到本地主机上的 Node.js 服务。原因是 Docker 容器中的 localhost 指向的是容器本身，而不是主机系统。为了让容器内的 Nginx 访问主机上的 Node.js 服务，你需要使用主机网络模式或者特定的 IP 地址来实现。

```nginx
upstream backend {
    server host.docker.internal:3001 weight=9;
    server host.docker.internal:3002 weight=1;
}
```

run two servers

```sh
node server1.js
node server2.js
```

run nginx container

```sh
docker run --name my-nginx -v /opt/homebrew/etc/nginx/nginx.conf:/etc/nginx/nginx.conf:ro -p 8080:8080 -d nginx
```


