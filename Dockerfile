FROM golang:1.23

# 静态二进制免去 alpine 上的 glibc 依赖问题。
ENV CGO_ENABLED=0

# URL 短链服务的源码、go.mod 和测试都在 /app 下运行。
WORKDIR /app

# 质检构建上下文为 Dockerfile + repo/,这里只生成运行所需二进制,不执行测试或 Git 初始化。
COPY repo/ .

RUN go build -o /app/server .

# 暴露 HTTP 服务端口，质检和评审可以通过 -p 端口映射访问 /health 与 /shorten。
EXPOSE 8789

# 容器默认启动短链服务，可通过 PORT 环境变量覆盖端口。
CMD ["/app/server"]
