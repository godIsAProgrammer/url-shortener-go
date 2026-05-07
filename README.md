# URL 短链服务

这是一个小型 Go HTTP 服务，把长 URL 收敛为 8 字符的短码并提供 302 重定向。
项目刻意保持轻量，只用 Go 标准库 `net/http` 和 `crypto/rand`，方便在最小容器
环境中运行，也方便在容器里验证短链生成和跳转流程。

## 端点

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/health` | 健康检查，返回 `{"ok": true}` |
| POST | `/shorten` | 请求体 `{"url": "https://..."}`，返回 `{"code": ..., "short_path": "/r/...", "original_url": ...}` |
| GET | `/r/{code}` | 根据短码 302 跳转到原始 URL，不存在时返回 404 |

短码生成使用 `crypto/rand` 取 8 位 base62，存储是进程内 `map[string]string`，
容器重启数据会丢失。

## 本地运行

```bash
go run .
# 默认监听 0.0.0.0:8789，可通过 PORT 环境变量覆盖
```

测试与生成示例：

```bash
go test ./...

curl http://127.0.0.1:8789/health
curl -X POST http://127.0.0.1:8789/shorten \
  -H 'content-type: application/json' \
  -d '{"url":"https://example.com/very/long/path"}'
curl -i http://127.0.0.1:8789/r/<code>
```

## Docker 环境

确保 Docker Desktop 已启动。

在项目根目录构建镜像：

```bash
docker build -t url-shortener-go .
```

启动 HTTP 服务：

```bash
docker run --rm -p 8789:8789 url-shortener-go
```

服务启动后，在另一个终端验证健康检查：

```bash
curl http://127.0.0.1:8789/health
```

预期响应：

```json
{"ok":true}
```

运行测试请使用显式命令：

```bash
docker run --rm url-shortener-go go test ./...
```

验证容器工作目录：

```bash
docker run --rm url-shortener-go pwd
```

预期输出为：

```text
/app
```

验证容器内初始仓库是否为干净 Git 工作区：

```bash
docker run --rm url-shortener-go git status --short
```

预期没有任何输出。
