# GoValidator 使用指南与规则速查

本项目使用了 `github.com/thedevsaddam/govalidator` 进行参数校验。`coreKit/validator` 包对其进行了封装，提供了便捷的结构体和请求参数验证方法。

## 1. 基本使用方式

在本项目中，主要通过 `coreKit/validator` 包提供的辅助函数进行验证。

### 1.1 结构体验证 (推荐)

用于验证 JSON Body 或已绑定到结构体的数据。

```go
import (
    "github.com/thedevsaddam/govalidator"
    "github.com/pudongping/momento-api/coreKit/validator"
)

// 1. 定义结构体标签
type LoginReq struct {
    Code string `json:"code" valid:"code"` // 注意：valid 标签对应规则中的字段名
}

// 2. 定义验证规则与自定义消息
func LoginRequestCheck(data interface{}) map[string][]string {
    rules := govalidator.MapData{
        "code": []string{"required", "min:2", "max:150"},
    }

    messages := govalidator.MapData{
        "code": []string{
            "required:code为必填项",
            "min:code长度不能小于2",
        },
    }

    return validator.ValidateStruct(data, rules, messages)
}

// 3. 在 Handler 中调用
// if msg, ok := validator.CallValidate(&req, requests.LoginRequestCheck); !ok {
//     responses.ToParamValidateResponse(r, w, nil, msg)
//     return
// }
```

### 1.2 HTTP Request 验证

用于验证 `form-data`, `x-www-form-urlencoded` 和 `query` 参数。

```go
rules := govalidator.MapData{
    "username": []string{"required", "between:3,8"},
    "email":    []string{"required", "min:4", "max:20", "email"},
}

// validator.Validate(r, rules, messages)
```

---

## 2. 验证规则速查表

以下规则整理自 `govalidator` 官方文档。规则可组合使用，例如 `required|min:3`。

| 规则 | 格式示例 | 说明 |
| :--- | :--- | :--- |
| **必须性** | | |
| `required` | `required` | 字段必须存在且非空。<br>空值定义：`nil`、空字符串 `""`、长度为0的 map/slice、数字 `0`。 |
| **字符串/字符类型** | | |
| `alpha` | `alpha` | 仅包含字母字符 (a-z, A-Z)。 |
| `alpha_dash` | `alpha_dash` | 仅包含字母、数字、破折号 `-` 和下划线 `_`。 |
| `alpha_space` | `alpha_space` | 仅包含字母、数字、破折号 `-`、下划线 `_` 和空格。 |
| `alpha_num` | `alpha_num` | 仅包含字母和数字。 |
| `numeric` | `numeric` | 必须完全由数字字符组成 (用于字符串类型的数字校验)。 |
| `bool` | `bool` | 必须能转换为布尔值。接受：`true`, `false`, `1`, `0`, `"1"`, `"0"`。 |
| `json` | `json` | 必须是有效的 JSON 字符串。 |
| **数值与范围** | | |
| `numeric_between` | `numeric_between:18,65` | 字段值必须是数值，且介于指定范围（含边界）。<br>支持整数与浮点数，如 `35`、`55.5` 均合法。<br>可省略任意边界：仅设最小值 `numeric_between:18,`（≥18）；仅设最大值 `numeric_between:,65`（≤65）；也可同时省略左右边界写成 `numeric_between:,` 表示任意数值均通过。 |
| `digits` | `digits:4` | 必须是数字且长度严格等于指定值 (如 `1234` 通过，`123` 失败)。 |
| `digits_between` | `digits_between:3,5` | 必须是数字且长度在指定范围内。 |
| `float` | `float` | 必须是有效的浮点数。 |
| `between` | `between:3,10` | 验证长度或范围。<br>- 字符串：字符长度<br>- 数组/切片/Map：元素个数<br>- 数字：数值大小范围 (两个整数或浮点数之间) |
| `min` | `min:3` | 最小长度或最小值。<br>- 字符串：最小字符数<br>- 切片/Map：最小元素数<br>- 数字：最小值 (支持整数或浮点数) |
| `max` | `max:10` | 最大长度或最大值。<br>- 字符串：最大字符数<br>- 切片/Map：最大元素数<br>- 数字：最大值 (支持整数或浮点数) |
| `len` | `len:4` | 精确长度或值。<br>- 字符串：精确字符数<br>- 切片/Map：精确元素数<br>- 数字：精确值 (支持整数或浮点数) |
| **枚举与集合** | | |
| `in` | `in:foo,bar,baz` | 字段值必须是列表中的一个。 |
| `not_in` | `not_in:admin,root` | 字段值不能是列表中的任何一个。 |
| **网络与格式** | | |
| `email` | `email` | 必须是有效的电子邮件地址。 |
| `url` | `url` | 必须是有效的 URL。 |
| `ip` | `ip` | 必须是有效的 IP 地址 (v4 或 v6)。 |
| `ip_v4` | `ip_v4` | 必须是有效的 IPv4 地址。 |
| `ip_v6` | `ip_v6` | 必须是有效的 IPv6 地址。 |
| `mac_address` | `mac_address` | 必须是有效的 MAC 地址。 |
| `uuid` | `uuid` | 必须是有效的 UUID。 |
| `uuid_v3` / `v4` / `v5`| `uuid_v4` | 必须是特定版本的 UUID。 |
| `coordinate` | `coordinate` | 必须是有效的坐标值。 |
| `lat` | `lat` | 必须是有效的纬度。 |
| `lon` | `lon` | 必须是有效的经度。 |
| `css_color` | `css_color` | 必须是有效的 CSS 颜色 (hex, rgb, rgba, hsl, hsla)。 |
| `credit_card` | `credit_card` | 必须是有效的信用卡号 (Visa, MasterCard, Amex 等)。 |
| **日期时间** | | |
| `date` | `date` | 必须是 `yyyy-mm-dd` 或 `yyyy/mm/dd` 格式的日期。 |
| `date:custom` | `date:dd-mm-yyyy` | 必须符合指定的日期格式。 |
| **文件验证** | | |
| `size` | `size:1024` | 文件大小 (字节)，仅适用于 `form-data`。 |
| `ext` | `ext:jpg,png` | 文件扩展名必须在列表中。 |
| `mime` | `mime:image/jpeg` | 文件 MIME 类型必须在列表中。 |
| **高级** | | |
| `regex` | `regex:^[a-z]+$` | 必须匹配指定的正则表达式。 |

## 3. 本项目自定义规则

除了官方规则外，本项目在 `coreKit/validator/rulesCnLength.go` 中扩展了以下规则：

| 规则 | 格式示例 | 说明 |
| :--- | :--- | :--- |
| `min_cn` | `min_cn:2` | 中文字符最小长度 (一个汉字算1个长度)。 |
| `max_cn` | `max_cn:10` | 中文字符最大长度 (一个汉字算1个长度)。 |