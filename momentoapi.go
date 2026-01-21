// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/coreKit/httpRest"
	"github.com/pudongping/momento-api/coreKit/responses"
	"github.com/pudongping/momento-api/cron"
	"github.com/pudongping/momento-api/internal/config"
	"github.com/pudongping/momento-api/internal/handler"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/momentoapi.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	opts := []rest.RunOption{
		rest.WithCustomCors(func(header http.Header) {
			header.Set("Access-Control-Allow-Headers", "Content-Type, Origin, X-CSRF-Token, Authorization, AccessToken, Token, Range, X-Token, X-Request-ID, X-Device-ID, X-User-ID")
		}, nil),
		rest.WithUnauthorizedCallback(httpRest.UnauthorizedHandler),
		rest.WithNotFoundHandler(httpRest.NotFoundHandler()),
		rest.WithNotAllowedHandler(httpRest.NotAllowedHandler()),
	}

	if c.Mode != "dev" || c.Mode != "test" {
		logx.DisableStat()
	}
	if c.Mode == "dev" {
		// 提供文件服务：https://go-zero.dev/faq/http/fileserver
		opts = append(opts, rest.WithFileServer("/static/public", http.Dir("public")))
	}

	server := rest.MustNewServer(c.RestConf, opts...)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 启动定时任务
	cron.Start(ctx)

	// 自定义错误
	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, interface{}) {
		logx.WithContext(ctx).Errorf("有错误，错误提示为：%+v 错误详情为：%+v", err.Error(), err)
		var customErr *errcode.Error
		switch e := err.(type) {
		case *errcode.Error:
			customErr = e
		default:
			customErr = errcode.InternalServerError
		}
		return http.StatusOK, responses.NewErrorResp(customErr.Code(), customErr.Msg())
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	if c.Mode == "dev" {
		server.PrintRoutes() // 打印路由表
	}

	server.Start()
}
