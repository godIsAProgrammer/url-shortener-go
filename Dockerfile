FROM golang:1.23

# 静态二进制免去 alpine 上的 glibc 依赖问题。
ENV CGO_ENABLED=0

# URL 短链服务的源码、go.mod 和测试都在 /app 下运行。
WORKDIR /app

# 复制路由、内存存储和 testing 用例作为任务起始现场。
COPY . /app/

# 先确认所有测试通过、二进制可成功编译，再把这个可工作的项目固化为 Git 初始提交。
RUN go test ./... \
    && go build -o /app/server . \
    && git init -b main \
    && git config user.email "agent@example.invalid" \
    && git config user.name "Agent Fixture" \
    && git add . \
    && git commit -m "Initial url shortener fixture"

# 暴露 HTTP 服务端口，质检和评审可以通过 -p 端口映射访问 /health 与 /shorten。
EXPOSE 8789

# 容器默认启动短链服务，可通过 PORT 环境变量覆盖端口。
CMD ["/app/server"]
