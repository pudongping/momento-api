package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/pudongping/momento-api/coreKit"
)

type SetUidToCtxMiddleware struct {
}

func NewSetUidToCtxMiddleware() *SetUidToCtxMiddleware {
	return &SetUidToCtxMiddleware{}
}

func (m *SetUidToCtxMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.ParseInt(r.Header.Get(coreKit.HeaderUserID), 10, 64)
		ctx := r.Context()
		ctx = context.WithValue(ctx, coreKit.CtxKeyJwtUserId, userId)

		next(w, r.WithContext(ctx))
	}
}
