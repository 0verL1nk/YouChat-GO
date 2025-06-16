# cwgo-template

## 项目介绍

参考云原生最佳实践，对 cwgo 的 http server 标准模板进行了优化

- docker 镜像构建
- 系统配置通过环境变量传入
- 日志输出到控制台

## 使用方式

在原有的生成命令末尾添加 `-template https://github.com/suyiiyii/cwgo-template.git -branch hertz` 参数

```shell
cwgo server -type HTTP ... -template https://github.com/suyiiyii/cwgo-template.git -branch hertz
```

例如，当前目录有 `hello.proto` 文件，则执行以下命令生成代码

```shell
cwgo server -type HTTP -service hello -module hello -idl hello.proto -template https://github.com/suyiiyii/cwgo-template.git -branch hertz
```

## 开发最佳实践

1. 编写 proto 文件，并使用相关注解定义 http 接口参数
2. 使用 `cwgo` 工具生成 hertz 框架脚手架代码（见上文）
3. 修改 `/biz/dal/model/model.go` 文件，将 User 结构体重命名为你的模型名称，并根据实际情况修改字段
4. 在模块根目录下创建 .env 文件，并配置数据库连接信息
5. 在模块根目录下执行 `go run cmd/gorm/main.go` 生成数据库表（会自动创建独立数据库，数据库名为服务名称）
6. 修改 `/biz/dal/model/model.go` 文件，在 Querier 结构体中添加自定义的 sql 语句（可选）
7. 在模块根目录下执行 `go run cmd/gorm_gen/main.go` 生成类型安全的数据库操作代码
8. 在 `/biz/service` 目录下实现业务逻辑
9. 在模块根目录下执行 `go run main.go` 启动服务

## 环境变量覆盖配置文件

默认支持使用环境变量覆盖配置文件中的配置，例如配置文件中的配置为：

```yaml
mysql:
  host: ""
  port: ""
  username: ""
  password: ""
```

可以在环境变量中设置以下变量：

```shell
APP_MYSQL_USERNAME
APP_MYSQL_PASSWORD
APP_MYSQL_HOST
APP_MYSQL_PORT
```

即使用 `APP_` 前缀加上配置的 key 的大写形式，层次使用 `_` 表示，作为环境变量名，即可覆盖配置文件中的配置

## 详细更改

在原来的基础上添加了

- Dockerfile 模板
- gorm 模型模板
- gorm gen 模板

在原来的基础上重构了

- 移除多环境配置文件，改为使用 viper 管理的单一配置文件，内嵌到二进制文件，并支持环境变量覆盖
- 移除文件日志，改为直接打印到控制台
- 数据库使用 mysql 并且添加自动创建数据库功能

## 参考

https://www.cloudwego.io/zh/docs/cwgo/tutorials/templete-extension/

