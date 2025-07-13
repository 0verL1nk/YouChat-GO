# 构建阶段
FROM golang:1.23.6-alpine AS builder

WORKDIR /app
ENV CGO_ENABLED=0 \
    GOOS=linux

# 复制依赖清单优先利用缓存
COPY go.mod go.sum ./
RUN go mod download -x

# 复制全部源代码
COPY . .
RUN cp conf/conf.example.yaml conf.yaml
# 构建可执行文件
RUN go build -o server . && \
    go build -o mingrate ./cmd/gorm/main.go

# 运行阶段
FROM alpine:3.21

WORKDIR /app
# 从构建阶段拷贝二进制文件
COPY --from=builder /app/server .
COPY --from=builder /app/mingrate .

# 时区环境变量
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8080

# 启动命令（生产环境建议使用具体配置）
CMD ["./server"]
