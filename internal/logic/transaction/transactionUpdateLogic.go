package transaction

import (
	"context"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransactionUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新交易记录
func NewTransactionUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransactionUpdateLogic {
	return &TransactionUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransactionUpdateLogic) TransactionUpdate(req *types.TransactionUpdateReq) (resp *types.TransactionUpdateResp, err error) {
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

	// 验证 TagId 是否存在（如果 TagId > 0）
	if req.TagId > 0 {
		_, err := l.svcCtx.TagsModel.FindOne(l.ctx, uint64(req.TagId))
		if err != nil {
			if err == model.ErrNotFound {
				return nil, errcode.BadRequest.Msgr("标签不存在")
			}
			return nil, errcode.DBError.Msgr("查询标签失败")
		}
	}

	updateMap := map[string]interface{}{
		"type":             req.Type,
		"amount":           req.Amount,
		"tag_id":           req.TagId,
		"remark":           strings.TrimSpace(req.Remark),
		"transaction_time": req.CreatedAt,
		"updated_at":       time.Now().Unix(),
	}

	// 使用 UpdateFilter 更新
	where := squirrel.Eq{"transaction_id": transactionId}
	_, err = l.svcCtx.TransactionsModel.UpdateFilter(l.ctx, nil, updateMap, where)
	if err != nil {
		l.Logger.Errorf("Update transaction failed: %v", err)
		return nil, errcode.DBError.Msgr("更新交易记录失败")
	}

	return &types.TransactionUpdateResp{}, nil
}
