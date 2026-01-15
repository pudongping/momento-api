package recurring

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecurringListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取周期性记账列表
func NewRecurringListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecurringListLogic {
	return &RecurringListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecurringListLogic) RecurringList(req *types.RecurringListReq) (resp *types.RecurringListResp, err error) {
	userID := cast.ToUint64(ctxData.GetUIDFromCtx(l.ctx))

	builder := l.svcCtx.RecurringTransactionsModel.SelectBuilder().Where(squirrel.Eq{
		"user_id": userID,
		"book_id": req.BookId,
	}).OrderBy("recurring_id DESC")

	list, err := l.svcCtx.RecurringTransactionsModel.FindAll(l.ctx, builder)
	if err != nil {
		l.Logger.Errorf("RecurringList FindAll error: %v", err)
		return nil, errcode.DBError.Msgr("查询周期性记账列表失败")
	}

	respList := make([]types.RecurringItem, 0, len(list))
	for _, item := range list {
		rItem := types.RecurringItem{}

		// 手动转换字段
		rItem.RecurringId = cast.ToString(item.RecurringId)
		rItem.Name = item.Name
		rItem.Type = item.Type
		rItem.Amount = item.Amount
		rItem.RecurringType = item.RecurringType
		rItem.RecurringDay = cast.ToInt64(item.RecurringDay)
		rItem.RecurringHour = cast.ToInt64(item.RecurringHour)
		rItem.RecurringMinute = cast.ToInt64(item.RecurringMinute)
		rItem.IsRecurringEnabled = item.IsRecurringEnabled

		respList = append(respList, rItem)
	}

	return &types.RecurringListResp{
		List: respList,
	}, nil
}
