# CLAUDE.md

本文件旨在为 Claude Code (claude.ai/code) 在此仓库中的工作提供指导。

## 项目概述

**momento-api** 是一个基于 go-zero 框架构建的微信小程序后端 API。这是一个个人记账应用（“时光账记”），用于管理用户、标签和节日。
项目技术栈包括：
- **框架**: go-zero (REST framework)
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
│   ├── transaction/    # 交易相关 API 定义
│   └── miniapp.api     # 主 API 入口 (导入所有子模块)
├── internal/
│   ├── config/         # 配置结构 (从 etc/momentoapi.yaml 加载)
│   ├── handler/        # HTTP 请求处理器 (从 .api 文件自动生成)
│   ├── logic/          # 业务逻辑层
│   ├── svc/            # 服务上下文 (依赖注入容器)
│   ├── service/        # 共享服务层 (业务逻辑复用)
│   ├── middleware/     # HTTP 中间件 (身份验证检查等)
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

1. **Request → Handler → Logic → Model**: 每个 API 接口均遵循以下流程：
   - Handler 接收 HTTP 请求，执行参数验证，随后调用 Logic
   - Logic 包含业务逻辑规则并调用 Model
   - Model 负责处理数据库操作

2. **Service Context**: `internal/svc/serviceContext.go` 负责初始化所有依赖项（MySQL, Redis, Model, Middleware）并将它们注入到整个应用中。

3. **代码生成工作流**:
   - `.api` 文件使用 goctl DSL 定义 API 契约
   - goctl 根据 `.api` 文件生成 Handler 存根和路由
   - goctl 根据 MySQL schema 生成数据库模型
   - `goctlTemplates/` 中的自定义模板确保生成的代码符合项目规范

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

# 生成 Markdown 格式的 API 文档
./local_run.sh mddoc

# 初始化 goctl 模板（一次性操作）
./local_run.sh tplinit

# 直接运行任意 goctl 命令
./local_run.sh goctl <args>
```

### 构建与运行
```bash
# 构建应用
go build -o momento-api momentoapi.go

# 运行应用
./momento-api -f etc/momentoapi.yaml
```

### 数据库
- Schema 文件存放于 `sql/` 目录
- 数据库连接配置位于 `etc/momentoapi.yaml`
- 使用 `./local_run.sh model <table_name>` 自动生成 Model

## 重要模式与规范

### .api 文件结构 (DSL)

新增的 `.api` 文件必须遵循 `dsl/API_STYLE_GUIDE.md` 中定义的 API 风格规范：

1. **命名**: 类型使用 PascalCase (例如 `TagListReq`, `UserInfoResp`)
2. **字段**:
   - 必须包含 `json` 标签
   - 可选字段在 json 标签中需添加 `,optional`，并配合 `valid` 标签进行验证
   - 字段注释应在同一行使用 `//`
   - 示例：`Type string \`json:"type,optional" valid:"type"\` // expense-支出 income-收入`
3. **ID**: 响应体中的 ID 字段应使用 `string` 类型，以防止前端出现精度丢失问题
4. **Handler 命名**: 使用 camelCase (例如 `tagList`, `userInfo`)
5. **路由**: 使用小写字母，并用 `/` 分隔 (例如 `/tags/list`, `/user/info`)

结构示例：
```go
// 在模块化 .api 文件中 (例如 dsl/tag/tag.api) - 仅包含 Type 定义
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

// 注意：所有的 @server 和 service 定义必须统一编写在 dsl/miniapp.api 文件中
```

### 身份验证

- 采用基于 JWT 的身份验证机制
- Token 处理逻辑位于 `coreKit/jwtToken/`
- `internal/middleware/authCheckMiddleware.go` 中的 `AuthCheckMiddleware` 负责验证 Token
- Token claims 中存储了用户信息，包括用户 ID

### 错误处理

- 使用 `coreKit/errcode` 包定义自定义错误
- 所有错误均转换为 `{"code": ..., "msg": "..."}` 格式的 JSON 响应
- `momentoapi.go` 中的主错误处理器将错误转换为 HTTP 200 及相应的错误代码

### 验证

- 在结构体字段上使用 `coreKit/validator` 和 `valid` 标签
- 在 `internal/requests/` 结构体中定义验证规则
- 使用 govalidator 库实现验证
- 参数验证方法应定义在 `./internal/requests` 目录下，具体验证规则请参考 `coreKit/validator/README.md` 文件。特别注意：验证字符串长度时，需使用 `min_cn` 和 `max_cn` 自定义验证规则。


### 数据库模型

使用以下命令从 MySQL schema 自动生成 Model：
```bash
./local_run.sh model <table_name>
```

生成的 Model 具备以下特性：
- 支持 Redis 缓存 (当启用 `--cache=true` 时)
- 缓存 Key 前缀：`momento_api:cache:`
- 提供标准 CRUD 方法

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
```
