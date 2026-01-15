package recurring

import (
	"context"

	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecurringDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除周期性记账规则
func NewRecurringDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecurringDeleteLogic {
	return &RecurringDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecurringDeleteLogic) RecurringDelete(req *types.RecurringDeleteReq) error {
	userID := cast.ToUint64(ctxData.GetUIDFromCtx(l.ctx))
	recurringId := cast.ToUint64(req.RecurringId)

	// 查询是否存在
	recurringData, err := l.svcCtx.RecurringTransactionsModel.FindOne(l.ctx, recurringId)
	if err != nil {
		if err == model.ErrNotFound {
			return errcode.NotFound.Msgr("周期性记账规则不存在")
		}
		return errcode.DBError.Msgr("查询周期性记账规则失败")
	}

	// 权限验证
	if recurringData.UserId != userID {
		return errcode.Forbidden.Msgr("无权操作此记录")
	}

	// 执行删除
	err = l.svcCtx.RecurringTransactionsModel.Delete(l.ctx, recurringId)
	if err != nil {
		l.Logger.Errorf("Delete recurring transaction failed: %v", err)
		return errcode.DBError.Msgr("删除周期性记账规则失败")
	}

	return nil
}
