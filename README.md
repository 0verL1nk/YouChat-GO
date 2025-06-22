# YouChat Backend GO 💬

## 项目介绍
YouChat Backend GO 是一个基于 <mcurl name="Hertz" url="https://github.com/cloudwego/hertz/">Hertz</mcurl> 框架开发的聊天应用后端服务。本项目采用现代化的技术栈和最佳实践，提供高性能、可扩展的即时通讯解决方案。🚀

## 主要特性 ✨
- 基于高性能的 Hertz 框架构建
- WebSocket 实时通讯支持 ⚡
- 完整的用户认证系统 🔐
- 群组聊天功能 👥
- 消息持久化存储 💾
- 分布式架构支持 🌐
- 完整的日志系统 📝
- 支持文件上传和管理 📤
- 集成了多种中间件（pprof、cors、recovery、access_log、gzip等）

## 技术栈 🛠️
- 框架：Hertz
- 数据库：MySQL 🗄️
- 缓存：Redis 💨
- 消息队列：Kafka 📨
- 对象存储：OSS ☁️
- 认证：JWT 🔑

## 目录结构 📂
```
├── biz/                    # 业务逻辑目录
│   ├── cerrors/           # 自定义错误定义
│   ├── chttp/             # HTTP 相关工具
│   ├── dal/               # 数据访问层
│   ├── handler/           # 请求处理器
│   ├── jwt/               # JWT 认证相关
│   ├── logger/            # 日志处理
│   ├── router/            # 路由配置
│   ├── service/           # 业务服务层
│   ├── socket/            # WebSocket 实现
│   └── utils/             # 工具函数
├── cmd/                    # 命令行工具
├── conf/                   # 配置文件
├── hertz_gen/             # Hertz 生成的代码
├── idl/                   # 接口定义文件
└── template/              # 代码模板
```

## 环境要求 📋
- Go 1.23.6+
- MySQL
- Redis
- Kafka
- 对象存储服务（如阿里云 OSS）

## 快速开始 ▶️

### 1. 配置环境
```bash
# 复制配置文件示例并修改
cp conf/conf.yaml.example conf/conf.yaml
```

### 2. 配置数据库
在 `.env` 文件中配置数据库连接信息：
```env
APP_MYSQL_HOST=localhost
APP_MYSQL_PORT=3306
APP_MYSQL_USERNAME=your_username
APP_MYSQL_PASSWORD=your_password
```

### 3. 初始化数据库
```bash
go run cmd/gorm/main.go
```

### 4. 生成数据访问层代码
```bash
go run cmd/gorm_gen/main.go
```

### 5. 构建和运行
```bash
# 方式1：直接运行
go run main.go

# 方式2：构建后运行
sh build.sh
sh output/bootstrap.sh
```

### 6. Docker 部署 🐳
```bash
# 使用 docker-compose 启动所有服务
docker-compose up -d
```

## 已实现功能
- [x] 用户注册与登录
- [x] 群组创建与加入
- [x] 消息发送（单聊、群聊）
- [x] WebSocket实时通信

## 开发计划
- [ ] 用户搜索功能
- [ ] 好友添加功能
- [ ] 群组管理（删除、修改）
- [ ] 消息撤回功能
- [ ] 通知系统
- [ ] 文件上传/下载
- [ ] 其他类型消息的收发

## 开发指南 🧑‍💻

### 1. 生成 API 代码
```bash
make gen svc=your_service_name
```

### 2. 项目结构说明
- `biz/handler`: 请求处理、参数验证和响应返回
- `biz/service`: 实际业务逻辑实现
- `biz/dal`: 数据访问层，包含数据库操作
- `biz/router`: 路由和中间件注册
- `biz/utils`: 通用工具方法

### 3. 测试 ✅
项目包含完整的单元测试框架，可以运行以下命令执行测试：
```bash
go test ./...
```

## 贡献指南 🤝
1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 Pull Request

## 许可证 📄
本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解更多详细信息。

## 联系方式 📧
如有任何问题或建议，请提交 issue 或 pull request。

---
**注意**: 在使用本项目之前，请确保已经正确配置了所有必要的环境变量和依赖服务。
        