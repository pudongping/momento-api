# CLAUDE.md

此文件为 Claude Code (claude.ai/code) 在本仓库中工作时提供指导。

## 项目概述

**momento-api** 是一个基于 go-zero 框架构建的微信小程序后端 API。这是一个个人记账应用（“时光小账本”），用于管理用户、标签和节日。项目使用了：
- **框架**: go-zero (rest framework)
- **语言**: Go 1.25.5
- **数据库**: MySQL
- **缓存**: Redis
- **代码生成**: goctl (go-zero 的代码生成器)

## 架构

项目遵循 go-zero 标准的 REST API 架构：

```
├── dsl/                 # API 定义文件 (.api 文件, goctl DSL)
│   ├── user/           # 用户相关 API 定义
│   ├── tag/            # 标签相关 API 定义
│   ├── festival/       # 节日相关 API 定义
│   └── miniapp.api     # 主 API 入口 (导入所有子模块)
├── internal/
│   ├── config/         # 配置结构 (从 etc/momentoapi.yaml 加载)
│   ├── handler/        # HTTP 请求处理器 (从 .api 文件自动生成)
│   ├── logic/          # 业务逻辑层
│   ├── svc/            # 服务上下文 (依赖注入容器)
│   ├── middleware/     # HTTP 中件间 (身份验证检查等)
│   ├── constant/       # Redis 常量和其他常量
│   ├── requests/       # 请求验证结构
│   ├── types/          # 类型定义
│   └── dao/            # 数据访问对象
├── model/              # 数据库模型 (从 MySQL schema 自动生成)
├── coreKit/            # go-zero 项目的共享工具库
│   ├── errcode/        # 自定义错误定义
│   ├── httpRest/       # HTTP 助手 (错误处理器, CORS)
│   ├── responses/      # 响应格式化工具
│   ├── jwtToken/       # JWT token 工具
│   ├── validator/      # 请求验证
│   ├── helpers/        # 通用助手
│   ├── middleware/     # 中间件模板
│   └── ctxData/        # 上下文数据工具
├── goctlTemplates/     # 自定义 goctl 模板 (覆盖默认生成)
├── etc/                # 配置文件 (YAML)
├── sql/                # 数据库 schema 和迁移脚本
├── local_run.sh        # 用于代码生成的开发辅助脚本
└── momentoapi.go       # 主入口点
```

### 关键架构模式

1. **Request → Handler → Logic → Model**: 每个 API 端点遵循此流程：
   - Handler 接收 HTTP 请求，进行验证，调用 Logic
   - Logic 包含业务规则并调用 Model
   - Model 处理数据库操作

2. **Service Context**: `internal/svc/serviceContext.go` 初始化所有依赖项（MySQL, Redis, Model, Middleware）并在整个应用程序中注入它们。

3. **代码生成工作流**:
   - `.api` 文件使用 goctl DSL 定义 API 契约
   - goctl 从 `.api` 文件生成 Handler 存根和路由
   - 使用 goctl 从 MySQL schema 生成数据库模型
   - `goctlTemplates/` 中的自定义模板确保生成的代码遵循项目规范

## 开发命令

### 代码生成
```bash
# 从 dsl/ 目录中的所有 .api 文件生成 Go 代码
make api
# 或手动执行：
./local_run.sh genapi

# 从 MySQL 表生成数据库模型
./local_run.sh model <table_name>
# 示例：./local_run.sh model users

# 格式化 API 定义文件
docker run --rm -it -v `pwd`:/app kevinwan/goctl:1.9.2 api format --dir ./dsl/<filename>.api

# 生成 markdown API 文档
./local_run.sh mddoc

# 初始化 goctl 模板 (一次性设置)
./local_run.sh tplinit

# 直接运行任何 goctl 命令
./local_run.sh goctl <args>
```

### 构建与运行
```bash
# 构建应用程序
go build -o momento-api momentoapi.go

# 运行应用程序
./momento-api -f etc/momentoapi.yaml
```

### 数据库
- Schema 文件位于 `sql/` 目录
- 数据库连接在 `etc/momentoapi.yaml` 中配置
- 使用 `./local_run.sh model <table_name>` 自动生成 Model

## 重要模式与规范

### .api 文件结构 (DSL)

新的 `.api` 文件必须遵循 `dsl/API_STYLE_GUIDE.md` 中的 API 风格指南：

1. **命名**: 类型使用 PascalCase (例如 `TagListReq`, `UserInfoResp`)
2. **字段**:
   - 必须带有 `json` 标签
   - 可选字段在 json 标签中需要 `,optional`，并带有用于验证的 `valid` 标签
   - 字段注释在同一行使用 `//`
   - 示例：`Type string \`json:"type,optional" valid:"type"\` // expense-支出 income-收入`
3. **ID**: 在响应体中为 ID 字段使用 `string` 类型，以避免前端精度丢失
4. **Handler 命名**: 使用 camelCase (例如 `tagList`, `userInfo`)
5. **路由**: 使用小写字母并以 `/` 分隔 (例如 `/tags/list`, `/user/info`)

结构示例：
```go
// 在模块化 .api 文件中 (例如 dsl/tag/tag.api) - 无需 syntax/info
type (
    TagListReq {
        Keyword string `json:"keyword,optional" valid:"keyword"` // 搜索关键词
        Type int64 `json:"type,optional" valid:"type"` // 1-系统 2-自定义
    }
    TagListResp {
        Id string `json:"id"`
        Title string `json:"title"`
    }
)

@server(
    group: tag
    tags: "标签相关"
    jwt: JWTAuth
    middleware: AuthCheckMiddleware
)
service momentoapi {
    @doc "获取标签列表"
    @handler tagList
    get /tags/list (TagListReq) returns ([]TagListResp)
}
```

### 身份验证

- 基于 JWT 的身份验证
- Token 在 `coreKit/jwtToken/` 中处理
- `internal/middleware/authCheckMiddleware.go` 中的 `AuthCheckMiddleware` 验证 Token
- Token claims 存储用户信息，包括用户 ID

### 错误处理

- 使用 `coreKit/errcode` 包定义自定义错误
- 所有错误都转换为 `{"code": ..., "msg": "..."}` JSON 响应
- `momentoapi.go` 中的主错误处理器将错误转换为 HTTP 200 及相应的错误代码

### 验证

- 在结构体字段上使用 `coreKit/validator` 和 `valid` 标签
- 在 `internal/requests/` 结构体中定义验证规则
- 使用 govalidator 库实现验证
- 参数需要验证时，将验证方法定义在 `./internal/requests` 目录下，相应的验证规则可详见 `coreKit/validator/README.md` 文件。需要特别注意的是，验证字符串长度时，需要使用 `min_cn` 和 `max_cn` 自定义验证规则方法。


### 数据库模型

使用以下命令从 MySQL schema 自动生成 Model：
```bash
./local_run.sh model <table_name>
```

生成的 Model 使用：
- Redis 缓存 (当使用 `--cache=true` 时)
- 缓存键前缀：`momento_api:cache:`
- 标准 CRUD 方法

### 配置

配置从 `etc/momentoapi.yaml` 加载：
```yaml
Name: momentoapi
Host: 0.0.0.0
Port: 8888

Mysql:
  DataSource: "root:password@tcp(127.0.0.1:3306)/momento"

Redis:
  Host: 127.0.0.1:6379

JWTAuth:
  AccessSecret: "your-secret"
  AccessExpire: 604800  # 7 天（以秒为单位）

WXMiniProgram:
  AppID: "your-app-id"
  AppSecret: "your-app-secret"
```

### HTTP Header

CORS 配置中处理自定义 Header：
- `X-Request-ID`: 用于幂等的请求 ID
- `X-Device-ID`: 设备标识符
- `X-User-ID`: 用户 ID (经过身份验证时)

## 添加新功能的流程

1. **定义 API 契约**: 按照风格指南在 `dsl/` 中创建或更新 `.api` 文件
2. **生成代码**: 运行 `make api` 自动生成 Handler 和路由
3. **实现 Logic**: 在 `internal/logic/` 目录中编写业务逻辑（需要注意业务逻辑解耦和代码适当的复用）
4. **创建/更新 Model**: 如果需要，使用 `./local_run.sh model <table_name>` 生成数据库模型
5. **添加验证**: 在 `.api` 文件中使用 `valid` 标签定义验证规则

## 特别注意事项

- **goctl 版本**: 1.9.2 (在 `local_run.sh` 中配置)
- **自定义模板**: 位于 `goctlTemplates/1.9.2/` - 这些模板覆盖了 goctl 的默认设置，以确保模块名称替换
- **模块命名**: `github.com/pudongping/momento-api` - 所有导入必须使用此路径
- **coreKit 复用性**: coreKit 旨在可复制到其他 go-zero 项目中；请勿在此处添加特定于业务的代码
- **local_run.sh 中的用户路径**: 为 `/Users/pudongping` 配置 - 如果项目位置发生变化，可能需要调整
