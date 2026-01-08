// 自定义中间件
// https://go-zero.dev/docs/tutorials/api/middleware
package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4/request"
	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/internal/constant"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"
)

type AuthCheckMiddleware struct {
	svc *svc.ServiceContext
}

func NewAuthCheckMiddleware(svc *svc.ServiceContext) *AuthCheckMiddleware {
	return &AuthCheckMiddleware{
		svc: svc,
	}
}

func (m *AuthCheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 这里就不需要验证 jwt token 了，因为在 go-zero 框架中已经做了 jwt token 的验证
		// 可以详见：go-zero -> handler.Authorize() 中间件的操作
		// 这里只需要做额外的授权检查即可，比如检查 token 是否在黑名单中，用户是否被禁用等

		token, err := request.BearerExtractor{}.ExtractToken(r)
		if err != nil {
			if errors.Is(err, request.ErrNoTokenInRequest) {
				// 表示没有获取到 token
				next(w, r)
				return
			}
			// 有其他的错误
			panic(err)
		}

		if "" == token {
			// 没有 token，则不需要做授权处理，直接放行
			next(w, r)
			return
		}

		// 判断 token 是否有效
		cacheKey := constant.RedisKeyPrefixToken + token
		userID, err := m.svc.RedisClient.GetCtx(r.Context(), cacheKey)
		if err != nil {
			panic(err)
		}
		if userID == "" {
			// token 已经失效
			unauthorized(w, r, "token 已经过期，请重新登录")
			return
		}

		// 检查数据库中该用户是否存在
		queryBuilder := m.svc.UserModel.SelectBuilder().
			Where("user_id = ?", cast.ToUint64(userID))
		user, err := m.svc.UserModel.FindOneByQuery(r.Context(), queryBuilder)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				unauthorized(w, r, "用户不存在")
				return
			}
			panic(errors.Wrap(err, "AuthCheckMiddleware.Handle 查询用户失败"))
		}
		if model.UserIsDisableYes == user.IsDisable {
			unauthorized(w, r, "用户已被禁用")
			return
		}

		next(w, r)

	}
}

func unauthorized(w http.ResponseWriter, r *http.Request, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	resp := map[string]interface{}{
		"code": http.StatusUnauthorized,
		"msg":  msg,
	}
	bs, _ := json.Marshal(resp)
	if _, err := w.Write(bs); err != nil {
		panic(fmt.Sprintf("authCheckMiddleware.unauthorized 写响应失败: %v", err))
	}
}
