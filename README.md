# momento-api

「时光小账本」API —— 每一笔账单，都是生活的旁白。

微信：

- [小程序用户头像昵称获取规则调整公告](https://developers.weixin.qq.com/community/develop/doc/00022c683e8a80b29bed2142b56c01)
- [小程序新的方式获取头像和昵称](https://developers.weixin.qq.com/miniprogram/dev/framework/open-ability/userProfile.html)

- [prf16/go-zero-box](https://github.com/prf16/go-zero-box/blob/main/app/doc/doc/doc.api)
- [go-zero-looklook](https://github.com/Mikaelemmmm/go-zero-looklook/blob/main/pkg/uniqueid/uniqueid.go)

## 1、配置环境

```bash
# GOROOT 设置 go 的版本
go env -w GOROOT='~/go/sdk/go1.25.5'

# GOPROXY
go env -w GOPROXY=https://goproxy.cn,direct
```

## 2、安装 goctl

```bash
# 直接安装
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 为了方便，直接使用 docker 进行安装
docker pull kevinwan/goctl:1.9.2
# 验证
docker run --rm -it -v `pwd`:/app kevinwan/goctl:1.9.2 --help
```

## 3、如果是创建一个全新的项目时

```bash
# 创建一个名称为 miniapp 的 API Rest 服务
goctl api new miniapp --style goZero
# 或者
docker run --rm -it -v `pwd`:/app kevinwan/goctl:1.9.2 api new miniapp --style goZero
```

## 4、也可以先生成 `*.api` 文件

```bash
# 在项目根目录下执行 
goctl api -o ./dsl/miniapp.api
# 或者
docker run --rm -it -v `pwd`:/app kevinwan/goctl:1.9.2 api -o ./dsl/miniapp.api
```

然后执行命令去自动生成 go 文件

```bash
# 在项目根目录下执行
goctl api go -api ./dsl/*.api -dir . --style=goZero
# 或者
docker run --rm -it -v `pwd`:/app kevinwan/goctl:1.9.2 api go -api ./dsl/*.api -dir . --style=goZero
```

api 文件格式化美观

```bash
# 好像只能一个文件一个文件的格式化
docker run --rm -it -v `pwd`:/app kevinwan/goctl:1.9.2 api format --dir ./dsl/miniapp.api
```


```mysql
CREATE TABLE `users` (
  `user_id` bigint(20) NOT NULL COMMENT '用户唯一ID（雪花算法）',
  `openid` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '微信openid',
  `unionid` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '微信unionid',
  `nickname` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户昵称',
  `avatar` varchar(1024) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户头像URL',
  `phone` varchar(80) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `created_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间（秒级时间戳）',
  `updated_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间（秒级时间戳）',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `uk_openid` (`openid`),
  KEY `idx_phone` (`phone`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
```