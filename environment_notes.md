# 环境说明

- 项目语言：Go 1.22+（仓库内使用 1.23）
- Docker 基础镜像：`golang:1.23`
- 容器工作目录：`/app`
- 构建时会把项目根目录的仓库文件复制到 `/app`
- 不需要第三方依赖（仅使用标准库 `net/http`、`crypto/rand`、`sync`）
- 默认启动命令：`/app/server`（监听 `0.0.0.0:8789`，可用 `PORT` 环境变量覆盖）
- 默认验证命令：`go test ./...`
- HTTP 端点：`GET /health`、`POST /shorten`、`GET /r/{code}`
- Dockerfile 会把 `/app` 初始化为 `main` 分支 Git 仓库，并创建一个初始提交

## 手动验证命令

```bash
docker build -t url-shortener-go .
docker run --rm -d -p 8789:8789 --name url-shortener-qc url-shortener-go
curl http://127.0.0.1:8789/health
curl -X POST http://127.0.0.1:8789/shorten \
  -H 'content-type: application/json' \
  -d '{"url":"https://example.com/long"}'
docker stop url-shortener-qc
docker run --rm url-shortener-go go test ./...
docker run --rm url-shortener-go pwd
docker run --rm url-shortener-go git status --short
```
