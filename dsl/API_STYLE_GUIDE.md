# .api 文件编写规范

为保证代码风格一致性，所有新建的 `.api` 文件必须严格遵循以下规范。

## 1. 结构体定义 (Types)

- **命名规范**: 使用 `大驼峰` 命名法 (e.g., `TagListReq`, `UserInfoResp`)。
- **字段标签**:
  - 必须包含 `json` 标签。
  - **可选参数**：必须在 `json` 标签中添加 `,optional`，并添加 `valid` 标签用于参数校验（配合 `govalidator`）。
  - 示例：`Type string `json:"type,optional" valid:"type"` // expense-支出 income-收入`
- **注释风格**:
  - 字段注释应写在行尾，使用 `//`。
  - 枚举值或特殊含义需详细说明（如：`// 1-系统标签 2-用户自定义标签`）。
- **类型选择**:
  - ID 类字段（如 `user_id`）如果涉及雪花算法或大整数，建议在响应体中使用 `string` 类型以避免前端精度丢失；请求体根据实际情况（如查询参数通常为 string）。
  - 数据库中的 `int` 状态字段对应 Go 的 `int64` 或 `int`。

## 2. 文件结构

- **无需头部 info**: 单个业务模块的 `.api` 文件（如 `dsl/tag/tag.api`）不需要包含 `syntax = "v1"` 和 `info` 块（这些通常在主入口 api 文件中定义，或者如果当前文件被 import，则尽量保持简洁，仅保留 `type` 定义）。

## 3. 服务定义 (Service) - 如果包含

- **注解**: 必须包含 `@doc` 和 `@handler`。
- **Doc**: 简明扼要地描述接口功能。
- **Handler**: 命名遵循 `小驼峰` (e.g., `tagList`)。
- **路由**: 使用全小写，单词间用 `/` 分隔 (e.g., `/tags/list`)。

## 4. 示例模版

```go
type (
    // 接口名Req
    ExampleListReq {
        Keyword string `json:"keyword,optional" valid:"keyword"` // 搜索关键字
        Type int64 `json:"type,optional" valid:"type"`           // 1-类型A 2-类型B
    }
    // 接口名Resp
    ExampleListResp {
        Id int64 `json:"id"`
        Title string `json:"title"`
        Status int64 `json:"status"` // 1-正常 2-禁用
    }
)

// 如果该文件包含 service 定义
// @server(
//     group: example
//     prefix: /api/v1
// )
// service momento-api {
//     @doc "获取示例列表"
//     @handler exampleList
//     get /example/list (ExampleListReq) returns (ExampleListResp)
// }
```

## 5. 关键差异点总结 (对比默认 goctl)

1.  **Valid 标签**: 显式添加 `valid` 标签用于自定义校验。
2.  **Optional**: 可选参数必须显式标记 `optional`。
3.  **行尾注释**: 推荐在字段后直接添加注释，而不是字段上方。
4.  **User ID 类型**: 响应体中 User ID 倾向于使用 `string` 类型。
