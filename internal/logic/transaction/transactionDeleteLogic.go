package transaction

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

type TransactionDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除交易记录
func NewTransactionDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransactionDeleteLogic {
	return &TransactionDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransactionDeleteLogic) TransactionDelete(req *types.TransactionDeleteReq) (*types.TransactionDeleteResp, error) {
	// 获取当前登录用户 ID
	userID := ctxData.GetUIDFromCtx(l.ctx)

	// 查询交易记录是否存在
	transactionId := cast.ToUint64(req.TransactionId)
	transactionData, err := l.svcCtx.TransactionsModel.FindOne(l.ctx, transactionId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errcode.NotFound.Msgr("交易记录不存在")
		}
		return nil, errcode.DBError.Msgr("查询交易记录失败")
	}

	// 权限验证
	if transactionData.UserId != uint64(userID) {
		return nil, errcode.Forbidden.Msgr("无权操作此记录")
	}

	// 执行删除
	err = l.svcCtx.TransactionsModel.Delete(l.ctx, transactionId)
	if err != nil {
		l.Logger.Errorf("Delete transaction failed: %v", err)
		return nil, errcode.DBError.Msgr("删除交易记录失败")
	}

	return &types.TransactionDeleteResp{}, nil
}
