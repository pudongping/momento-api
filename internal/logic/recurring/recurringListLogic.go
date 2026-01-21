package recurring

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
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

	// 1. 收集所有 TagId
	tagIDs := make([]uint64, 0)
	for _, item := range list {
		if item.TagId > 0 {
			tagIDs = append(tagIDs, item.TagId)
		}
	}

	// 2. 批量查询标签信息
	tagMap := make(map[uint64]*model.Tags)
	if len(tagIDs) > 0 {
		tagBuilder := l.svcCtx.TagsModel.SelectBuilder().Where(squirrel.Eq{"tag_id": tagIDs})
		tags, err := l.svcCtx.TagsModel.FindAll(l.ctx, tagBuilder)
		if err != nil {
			l.Logger.Errorf("RecurringList FindAll Tags error: %v", err)
			// 即使标签查询失败，也不影响主流程，只是缺少标签信息
		} else {
			for _, tag := range tags {
				tagMap[tag.TagId] = tag
			}
		}
	}

	respList := make([]types.RecurringItem, 0, len(list))
	for _, item := range list {
		rItem := types.RecurringItem{}

		// 手动转换字段
		rItem.RecurringId = cast.ToString(item.RecurringId)
		rItem.BookId = cast.ToInt64(item.BookId)
		rItem.UserId = cast.ToString(item.UserId)
		rItem.Name = item.Name
		rItem.Type = item.Type
		rItem.Amount = item.Amount
		rItem.TagId = cast.ToInt64(item.TagId)
		rItem.Remark = item.Remark
		rItem.RecurringType = item.RecurringType
		rItem.RecurringHour = cast.ToInt64(item.RecurringHour)
		rItem.RecurringMinute = cast.ToInt64(item.RecurringMinute)
		rItem.RecurringWeekday = cast.ToInt64(item.RecurringWeekday)
		rItem.RecurringMonth = cast.ToInt64(item.RecurringMonth)
		rItem.RecurringDay = cast.ToInt64(item.RecurringDay)
		rItem.IsRecurringEnabled = item.IsRecurringEnabled
		rItem.LastExecutedAt = cast.ToInt64(item.LastExecutedAt)
		rItem.CreatedAt = cast.ToInt64(item.CreatedAt)
		rItem.UpdatedAt = cast.ToInt64(item.UpdatedAt)

		// 填充标签信息
		if tag, ok := tagMap[item.TagId]; ok {
			rItem.TagName = tag.Name
			rItem.TagColor = tag.Color
			rItem.TagIcon = tag.Icon
		}

		respList = append(respList, rItem)
	}

	return &types.RecurringListResp{
		List: respList,
	}, nil
}
