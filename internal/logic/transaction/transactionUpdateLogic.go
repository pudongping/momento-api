package transaction

import (
	"context"
	"strings"
	"time"

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

	// 更新字段
	if req.Type != "" {
		transactionData.Type = req.Type
	}
	if req.Amount > 0 {
		transactionData.Amount = req.Amount
	}
	if req.TagId > 0 {
		// 验证 TagId 是否存在
		tag, err := l.svcCtx.TagsModel.FindOne(l.ctx, uint64(req.TagId))
		if err != nil {
			if err == model.ErrNotFound {
				return nil, errcode.BadRequest.Msgr("标签不存在")
			}
			return nil, errcode.DBError.Msgr("查询标签失败")
		}
		transactionData.TagId = uint64(req.TagId)
		// 如果更新了 TagId，后续返回时需要用到新标签信息，这里可以先存下来或者重新查
		// 实际上 FindOne 已经查到了 tag
		_ = tag
	}
	// 备注允许为空，但如果没传则不更新？
	// 规范说：Remark String 否 备注
	// 通常如果传了空字符串可能意为清空，但这里 req.Remark 是 optional，如果前端没传，go-zero 解析出来可能是空字符串。
	// 但如果前端传了 ""，也是空字符串。
	// 我们假设如果字段在 JSON 中存在且为空串，则是清空。
	// 但 go-zero optional 字段如果没传，也是默认值。
	// 这里简单处理：只要 req.Remark 不为空，或者显式逻辑。
	// 由于 go-zero 无法区分 "没传" 和 "传了空值" (除非用指针)，这里我们假设如果传了就更新。
	// 实际上 Update 接口通常全量更新或部分更新。
	// 如果前端传了该字段，我们应当更新。但这里我们无法区分。
	// 鉴于 go-zero 的 struct tag `optional`，我们只能认为非零值即更新，或者是覆盖更新。
	// 考虑到备注可能是要清空的，如果 req.Remark 是空字符串，我们是否应该更新？
	// 这是一个常见问题。通常做法是：如果不传，保持原样；如果传了空串，清空。
	// 但 struct 无法区分。
	// 让我们看 TagId，如果是 0 (默认值)，我们不更新。
	// Remark 如果是 ""，我们不更新？ 那用户怎么清空备注？
	// 这是一个权衡。
	// 暂时策略：如果 req.Remark != "" 则更新。如果用户想清空，可能需要传一个特殊值或者我们无法支持清空（只能改成空格）。
	// 或者我们假设前端会把所有字段都传过来（全量更新），那样就直接赋值。
	// 文档说 "更新已存在的普通交易记录"，参数都是可选。这暗示是 PATCH 语义。
	// 让我们看看 transactionData.Remark。
	// 如果 req.Remark != ""，更新。
	if req.Remark != "" {
		transactionData.Remark = strings.TrimSpace(req.Remark)
	}
	// 如果 req.CreatedAt > 0，更新
	if req.CreatedAt > 0 {
		transactionData.TransactionTime = uint64(req.CreatedAt)
	}

	transactionData.UpdatedAt = uint64(time.Now().Unix())

	err = l.svcCtx.TransactionsModel.Update(l.ctx, transactionData)
	if err != nil {
		l.Logger.Errorf("Update transaction failed: %v", err)
		return nil, errcode.DBError.Msgr("更新交易记录失败")
	}

	return &types.TransactionUpdateResp{}, nil
}
