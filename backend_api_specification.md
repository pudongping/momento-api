# 后端 API 接口开发规范文档

## 1. 概述 (Overview)

本主要基于前端 Mock 数据逻辑、现有页面交互逻辑以及已定稿的数据库 Schema 整理而成。后端开发需严格遵循本各接口的定义、参数要求及业务逻辑。

### 1.1 通用约定
- **协议**: HTTP / HTTPS
- **数据格式**: JSON
- **字符编码**: UTF-8
- **时间处理**: 所有时间字段（如 `created_at`, `updated_at`, `transaction_time` 等）均使用 **10位秒级时间戳** (Integer)。
- **ID 类型**:
    - `user_id`: 字符串类型 (String)，对应数据库 `VARCHAR` 或 `BIGINT` (Snowflake ID)。
    - `book_id`, `tag_id`, `transaction_id`: 数值类型 (Integer/Long)。
- **布尔值映射**: 
    - 数据库中通常使用 `TINYINT` 表示布尔状态：`1` 代表是/启用 (True/Yes)，`2` 代表否/禁用 (False/No)。
    - 接口返回时请注意保持一致性，或根据前端 Mock 数据格式（部分 Mock 直接使用了布尔值，需确认转换，建议后端统一返回数值状态码或布尔值，本文档以数据库定义的数值 `1/2` 为主，若 Mock 用布尔值则注明）。
- **响应结构 (Response Structure)**:
    ```json
    {
      "code": 200,          // 状态码，200 表示成功
      "msg": "success", // 提示信息
      "data": { ... }       // 业务数据
    }
    ```

---

## 2. 交易记录模块 (Transactions)

### 2.1 获取交易记录列表
- **接口地址**: `GET /transactions/list`
- **描述**: 获取交易流水列表，支持多条件筛选及分页（分页功能采用游标分页，第一页时 `last_transaction_id` 为 0，后续请求需传入上一页最后一条记录 transaction_id 以获取下一页数据，编写接口时，后端默认 page 永远为 1）。
- **请求参数 (Query Parameters)**:

| 参数名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `book_id` | Integer | 是 | 账本 ID |
| `type` | String | 否 | 交易类型: `expense` (支出), `income` (收入) |
| `tag_id` | Integer | 否 | 标签 ID |
| `start_date` | Integer | 否 | 开始时间戳 (秒) |
| `end_date` | Integer | 否 | 结束时间戳 (秒) |
| `page` | Integer | 否 | 页码，默认 1 |
| `per_page` | Integer | 否 | 每页数量，默认 20 |
| `last_transaction_id` | Integer | 否 | 上一页最后一条记录 ID (用于游标分页优化，必填) |

- **响应数据 (data)**:
    ```json
    {
      "list": [
        {
          "transaction_id": 1001,
          "book_id": 1,
          "user_id": "123456789",
          "type": "expense",
          "amount": 25.50,
          "tag_id": 10,
          "tag_name": "餐饮",
          "tag_color": "#FF5733",
          "tag_icon": "food.png",
          "remark": "午餐",
          "transaction_time": 1700000000,
          "created_at": 1700000000,
          "updated_at": 1700000000,
          "is_auto_generated": 2  // 1-是(自动生成) 2-否(手动)
        }
      ],
      "has_more": true,  // 是否有更多数据
      "total": 100,      // 总记录数
      "page": 1,
      "per_page": 20
    }
    ```

### 2.2 添加交易记录 (含周期性记账创建)
- **接口地址**: `POST /transactions/add`
- **描述**: 该接口为**复合接口**。前端通过 `is_recurring_enabled` 字段控制是创建普通交易还是周期性记账规则。
- **业务逻辑**:
    1.  **普通记账模式** (`is_recurring_enabled` 为 `false` 或未传):
        - 仅使用公共核心参数。
        - 数据写入 `transactions` 表。
        - **注意**: 前端传来的 `created_at` 参数对应数据库中的 `transaction_time` 字段。
    2.  **周期性记账模式** (`is_recurring_enabled` 为 `true`):
        - 使用公共参数 + 周期性规则参数。
        - 数据写入 `recurring_transactions` 表。
        - `name` 字段逻辑：优先使用 `remark`，若为空则自动生成（如 "周期支出"）。

- **请求参数 (Body)**:

| 参数分类 | 参数名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- | :--- |
| **公共参数** | `book_id` | Integer | 是 | 账本 ID |
| | `type` | String | 是 | `expense` / `income` |
| | `amount` | Decimal | 是 | 金额 |
| | `tag_id` | Integer | 是 | 标签 ID |
| | `remark` | String | 否 | 备注 (去除首尾空格) |
| | `created_at` | Integer | 是 | **交易发生时间** (对应 DB `transaction_time`) |
| **周期参数** | `is_recurring_enabled` | Boolean | 否 | **模式开关**，`true` 表示周期记账 |
| (当开关开启时) | `name` | String | 否 | 规则名称 (前端逻辑：`remark` \|\| "周期{type}") |
| | `recurring_type` | String | 是 | `daily`, `weekly`, `monthly`, `quarterly`, `yearly` |
| | `recurring_hour` | Integer | 是 | 执行小时 (0-23) |
| | `recurring_minute` | Integer | 是 | 执行分钟 (0-59) |
| **类型特定** | `recurring_weekday` | Integer | 条件 | **Weekly必填**: 0-6 (0为周日) |
| | `recurring_month` | Integer | 条件 | **Monthly必填**: 1-12 |
| | `recurring_day` | Integer | 条件 | **Monthly必填**: 1-31 |

- **请求示例**:

**1. 普通记账**:
```json
{
  "book_id": 1,
  "type": "expense",
  "amount": 25.5,
  "tag_id": 101,
  "remark": "午餐",
  "created_at": 1705286400
}
```

**2. 周期性记账 (每月15号)**:
```json
{
  "book_id": 1,
  "type": "expense",
  "amount": 3000,
  "tag_id": 201,
  "remark": "房租",
  "created_at": 1705286400,
  "is_recurring_enabled": true,
  "name": "房租",
  "recurring_type": "monthly",
  "recurring_hour": 9,
  "recurring_minute": 0,
  "recurring_month": 1,
  "recurring_day": 15
}
```

- **响应数据 (data)**:
    - 返回创建成功的记录对象 (结构同 list 中的 item)。



### 2.3 更新交易记录
- **接口地址**: `PUT /transactions/update`
- **描述**: 更新已存在的普通交易记录。
- **请求参数 (Body)**:

| 参数名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `transaction_id` | Integer | 是 | 交易记录 ID |
| `type` | String | 否 | `expense` / `income` |
| `amount` | Decimal | 否 | 金额 |
| `tag_id` | Integer | 否 | 标签 ID |
| `remark` | String | 否 | 备注 |
| `transaction_time` | Integer | 否 | 交易时间 |

- **响应数据 (data)**: 更新后的交易记录对象。

### 2.4 删除交易记录
- **接口地址**: `DELETE /transactions/delete`
- **请求参数 (Body)**:
    - `transaction_id`: Integer, 必填
- **响应数据 (data)**: 

### 2.5 获取交易统计
- **接口地址**: `GET /transactions/stats`
- **描述**: 根据筛选条件计算收支总额及各标签占比。
- **请求参数**: 同 `GET /transactions/list` (通常不传分页参数)。
- **响应数据 (data)**:
    ```json
    {
      "total_income": 5000.00,
      "total_expense": 3000.00,
      "balance": 2000.00,
      "tag_stats": [
        {
          "tag_id": 10,
          "tag_name": "餐饮",
          "tag_color": "#FF5733",
          "tag_icon": "food.png",
          "count": 15,       // 笔数
          "amount": 1500.00  // 总金额
        }
      ]
    }
    ```

---

## 3. 周期性记账模块 (Recurring Transactions)

### 3.1 获取周期性记账列表
- **接口地址**: `GET /recurring/list`
- **描述**: 获取用户设置的自动记账规则列表，周期性记账列表不需要分页。
- **请求参数 (Query Parameters)**:
    - `book_id`: Integer, 必填
- **响应数据 (data)**:
    ```json
    {
      "list": [
        {
          "recurring_id": 1,
          "name": "每月房租",
          "type": "expense",
          "amount": 3000.00,
          "recurring_type": "monthly",
          "recurring_day": 15,
          "recurring_hour": 9,
          "recurring_minute": 0,
          "is_recurring_enabled": 1, // 1-启用 2-禁用
          "next_execution_time": 1702602000 // (可选) 计算出的下一次执行时间
        }
      ]
    }
    ```

### 3.2 删除周期性记账规则
- **接口地址**: `DELETE /recurring/delete`
- **请求参数 (Body)**:
    - `recurring_id`: Integer, 必填
- **响应数据**: 

---

## 4. 账本模块 (Account Books)

### 4.1 获取账本列表
- **接口地址**: `GET /accountBooks/list`
- **描述**: 获取当前用户参与的所有账本（包括自己创建的和加入的）。
- **响应数据 (data)**:
    ```json
    [
      {
        "book_id": 1,
        "name": "家庭账本",
        "creator_user_id": "12345",
        "is_creator": 1,   // 1-是创建者 2-不是
        "is_default": 1,   // 1-是默认账本 2-不是
        "member_count": 2, // 成员数量
        "created_at": 1700000000
      }
    ]
    ```

### 4.2 创建账本
- **接口地址**: `POST /accountBooks/create`
- **请求参数**:
    - `name`: String, 必填 (账本名称、min_cn:2, max_cn:8)
- **业务逻辑**: 创建账本的同时，将当前用户作为 Creator 加入 `account_book_users` 表，并设置 `is_creator=1`。
- **响应数据**: 新创建的账本对象。

### 4.3 删除账本
- **接口地址**: `DELETE /accountBooks/delete`
- **请求参数**: `book_id`
- **业务逻辑**: 仅 Creator 可删除。需级联删除该账本下的所有交易记录、周期规则及成员关系。

### 4.4 邀请用户加入账本
- **接口地址**: `POST /accountBooks/invite`
- **请求参数**:
    - `book_id`: Integer
    - `target_uid`: String (被邀请人的 User ID) 被邀请人必须一定存在于系统中
- **业务逻辑**: 创建一条邀请记录，状态为 `pending`。

### 4.5 获取邀请列表
- **接口地址**: `GET /accountBooks/invitations`
- **描述**: 获取当前用户收到的邀请。
- **响应数据**:
    ```json
    [
      {
        "invitation_id": 101,
        "book_id": 3,
        "book_name": "聚餐账本",
        "inviter_uid": "222",
        "inviter_nickname": "老王",
        "target_uid": "123456789012345678",
        "status": "pending", // pending, accepted, rejected
        "created_at": 1700000000
      }
    ]
    ```

### 4.6 处理邀请 (接受/拒绝)
- **接受**: `POST /accountBooks/accept`, 参数 `{ "invitation_id": 1 }`
- **拒绝**: `POST /accountBooks/reject`, 参数 `{ "invitation_id": 1 }`

### 4.7 退出账本
- **接口地址**: `POST /accountBooks/exit`
- **参数**: `book_id`
- **业务逻辑**: 非 Creator 可退出。删除 `account_book_users` 关联。

### 4.8 设置默认账本
- **接口地址**: `PUT /accountBooks/setDefault`
- **参数**: `book_id`
- **业务逻辑**: 将用户其他账本的 `is_default` 设为 2，当前账本设为 1。

### 4.9 获取账本成员
- **接口地址**: `GET /accountBooks/members`
- **参数**: `book_id`
- **响应数据**:
    ```json
    [
      {
        "user_id": "123",
        "nickname": "张三",
        "avatar": "http://...",
        "status": "joined" // joined, waiting
      }
    ]
    ```

### 4.10 移除成员
- **接口地址**: `POST /accountBooks/removeMember`
- **请求参数 (Body)**:
    - `book_id`: Integer, 必填
    - `target_uid`: String, 必填 (被移除成员的用户ID)
- **请求示例**:
    ```json
    {
      "book_id": 1,
      "target_uid": "223456789012345678"
    }
    ```
- **权限**: 仅 Creator 可操作。
- **业务逻辑**: 从 `account_book_users` 表中移除该用户，使其无法再访问此账本。
- **响应数据**: 

---

## 5. 文件上传模块 (File Upload)

### 5.1 上传文件
- **接口地址**: `POST /upload/file`
- **Content-Type**: `multipart/form-data`
- **请求参数**:
    - `file`: Binary (文件流)
    - `file_type`: String (`image`, `video`, `excel`, etc.)
    - `business_type`: String (`user_avatar`, `background`, `transaction_image` 等，用于区分存储路径)
- **响应数据 (data)**:
    ```json
      {
        "file_id": 1234567890,
        "user_id": "123456789012345678",
        "relative_path": "/uploads/user_avatar/1704067200_abc123.jpg",
        "absolute_url": "https://cdn.example.com/uploads/user_avatar/1704067200_abc123.jpg",
        "file_size": 102400,
        "file_type": "image",
        "business_type": "user_avatar",
        "upload_time": 1704067200
      }
    ```
- **业务逻辑**:
    1. 校验文件大小和格式。
    2. 保存文件到服务器存储或 OSS。
    3. 记录文件信息到 `upload_files` 表。
    4. 返回可访问的 URL。
