## 项目规范

- `.api` 文件创建规范，详见 [API_STYLE_GUIDE.md](./dsl/API_STYLE_GUIDE.md)
- 参数需要验证时，将验证方法定义在 `./internal/requests` 目录下，相应的验证规则可详见 `coreKit/validator/README.md` 文件。需要特别注意的是，验证字符串长度时，需要使用 `min_cn` 和 `max_cn` 自定义验证规则方法。