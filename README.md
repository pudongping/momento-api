# momento-api

「时光小账本」API —— 每一笔账单，都是生活的旁白。

微信：

- [小程序用户头像昵称获取规则调整公告](https://developers.weixin.qq.com/community/develop/doc/00022c683e8a80b29bed2142b56c01)
- [小程序新的方式获取头像和昵称](https://developers.weixin.qq.com/miniprogram/dev/framework/open-ability/userProfile.html)

- [prf16/go-zero-box](https://github.com/prf16/go-zero-box/blob/main/app/doc/doc/doc.api)

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