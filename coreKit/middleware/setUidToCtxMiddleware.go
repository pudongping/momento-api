package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/pudongping/momento-api/coreKit/ctxData"
)

const (
	HeaderUserID = "X-User-ID" // 从客户端请求头中需要带过来的用户ID字段
)

type SetUidToCtxMiddleware struct {
}

func NewSetUidToCtxMiddleware() *SetUidToCtxMiddleware {
	return &SetUidToCtxMiddleware{}
}

func (m *SetUidToCtxMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.ParseInt(r.Header.Get(HeaderUserID), 10, 64)
		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxData.CtxKeyJwtUserId, userId)

		next(w, r.WithContext(ctx))
	}
}
