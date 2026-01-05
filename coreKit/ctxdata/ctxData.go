package ctxdata

import (
	"context"
	"encoding/json"

	"github.com/pudongping/momento-api/coreKit"
	"github.com/zeromicro/go-zero/core/logx"
)

func GetUIDFromCtx(ctx context.Context) int64 {
	var uid int64
	jsonUID, ok := ctx.Value(coreKit.CtxKeyJwtUserId).(json.Number)
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
