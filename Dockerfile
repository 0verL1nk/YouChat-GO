
ARG GO_VERSION=1.23.6
FROM golang:${GO_VERSION} AS build
WORKDIR /src

# 复制依赖
COPY go.mod go.sum ./
RUN go mod download -x

# 复制全部源代码（包含 conf.yaml 或 conf.example.yaml）
COPY . .
COPY ./conf/conf.example.yaml ./conf.yaml
# 构建二进制
ARG TARGETARCH
RUN CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server .
RUN CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/migrator ./cmd/gorm/main.go

FROM debian:12-slim AS final

# Create a non-privileged user that the app will run under.
# See https://docs.docker.com/go/dockerfile-user-best-practices/
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
RUN mkdir -p /app/logs && chown -R ${UID}:${UID} /app
USER appuser

# Copy the executable from the "build" stage.
COPY --from=build /bin/server /bin/
COPY --from=build /bin/migrator /bin/

# Expose the port that the application listens on.
EXPOSE 8080

# What the container should run when it is started.
ENTRYPOINT [ "/bin/server" ]