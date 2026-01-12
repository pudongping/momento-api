package ctxData

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	CtxKeyJwtUserId = "jwtUserId" // ctx 中存储用户ID的 key
)

func GetUIDFromCtx(ctx context.Context) int64 {
	var uid int64
	jsonUID, ok := ctx.Value(CtxKeyJwtUserId).(json.Number)
	if !ok {
		return 0
	}

	if intUid, err := jsonUID.Int64(); err == nil {
		uid = intUid
	} else {
		logx.WithContext(ctx).Errorf("GetUIDFromCtx err : %+v", err)
	}
	return uid
}
