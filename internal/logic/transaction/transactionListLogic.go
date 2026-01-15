package transaction

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/coreKit/paginator"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransactionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取交易流水列表
func NewTransactionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransactionListLogic {
	return &TransactionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransactionListLogic) TransactionList(req *types.TransactionListReq) (resp *types.TransactionListResp, err error) {
	userID := cast.ToUint64(ctxData.GetUIDFromCtx(l.ctx))

	// 1. 构建基础查询条件
	whereBuilder := l.buildBaseQuery(userID, req)

	// 2. 计算总数
	total, err := l.countTotal(whereBuilder)
	if err != nil {
		return nil, err
	}

	_, limit, _ := paginator.PrepareOffsetLimit(req.Page, req.PerPage)

	// 如果总数为 0，直接返回空结果，避免多余的查询
	if total == 0 {
		return &types.TransactionListResp{
			List:    []types.TransactionItem{},
			HasMore: false,
			Total:   0,
			Page:    req.Page,
			PerPage: limit,
		}, nil
	}

	// 3. 构建列表查询并应用游标分页
	builder := l.svcCtx.TransactionsModel.SelectBuilder().Where(whereBuilder)
	builder = l.applyCursorPagination(builder, req.LastTransactionId)

	// 4. 设置排序和分页限制
	// 默认只按 transaction_id 倒序排列
	builder = builder.OrderBy("transaction_id DESC").Limit(uint64(limit + 1))

	// 5. 执行查询
	list, err := l.svcCtx.TransactionsModel.FindAll(l.ctx, builder)
	if err != nil {
		l.Logger.Errorf("TransactionList FindAll error: %v", err)
		return nil, errcode.DBError.WithError(err)
	}

	// 6. 处理分页标志
	hasMore := false
	if len(list) > int(limit) {
		hasMore = true
		list = list[:limit]
	}

	// 7. 获取标签信息
	tagMap := l.getTagsMap(list)

	// 8. 组装响应数据
	respList := l.assembleTransactionItems(list, tagMap)

	return &types.TransactionListResp{
		List:    respList,
		HasMore: hasMore,
		Total:   total,
		Page:    req.Page,
		PerPage: limit,
	}, nil
}

// buildBaseQuery 构建通用的过滤条件
func (l *TransactionListLogic) buildBaseQuery(userID uint64, req *types.TransactionListReq) squirrel.Sqlizer {
	conditions := squirrel.And{
		squirrel.Eq{"user_id": userID},
		squirrel.Eq{"book_id": req.BookId},
	}

	if req.Type != "" {
		conditions = append(conditions, squirrel.Eq{"type": req.Type})
	}
	if req.TagId > 0 {
		conditions = append(conditions, squirrel.Eq{"tag_id": req.TagId})
	}
	if req.StartDate > 0 {
		conditions = append(conditions, squirrel.GtOrEq{"transaction_time": req.StartDate})
	}
	if req.EndDate > 0 {
		conditions = append(conditions, squirrel.LtOrEq{"transaction_time": req.EndDate})
	}

	return conditions
}

// countTotal 计算总数
func (l *TransactionListLogic) countTotal(where squirrel.Sqlizer) (int64, error) {
	countBuilder := l.svcCtx.TransactionsModel.CountBuilder().Where(where)
	total, err := l.svcCtx.TransactionsModel.FindCount(l.ctx, countBuilder)
	if err != nil {
		l.Logger.Errorf("TransactionList FindCount error: %v", err)
		return 0, errcode.DBError.WithError(err)
	}
	return total, nil
}

// applyCursorPagination 应用游标分页
func (l *TransactionListLogic) applyCursorPagination(builder squirrel.SelectBuilder, lastID int64) squirrel.SelectBuilder {
	// 如果 lastID > 0，说明需要加载更早的数据（transaction_id 更小）
	if lastID > 0 {
		return builder.Where(squirrel.Lt{"transaction_id": lastID})
	}
	// 如果 lastID == 0，说明是第一页，不需要额外条件
	return builder
}

// getTagsMap 获取标签 Map
func (l *TransactionListLogic) getTagsMap(list []*model.Transactions) map[uint64]*model.Tags {
	tagMap := make(map[uint64]*model.Tags)
	tagIDSet := make(map[uint64]struct{})

	for _, item := range list {
		if item.TagId > 0 {
			tagIDSet[item.TagId] = struct{}{}
		}
	}

	if len(tagIDSet) == 0 {
		return tagMap
	}

	tagIDs := make([]uint64, 0, len(tagIDSet))
	for id := range tagIDSet {
		tagIDs = append(tagIDs, id)
	}

	tagBuilder := l.svcCtx.TagsModel.SelectBuilder().Where(squirrel.Eq{"tag_id": tagIDs})
	tags, err := l.svcCtx.TagsModel.FindAll(l.ctx, tagBuilder)
	if err != nil {
		l.Logger.Errorf("TransactionList FindAll Tags error: %v", err)
		// 标签查询失败不影响主流程，只是没有标签信息
		return tagMap
	}

	for _, t := range tags {
		tagMap[t.TagId] = t
	}

	return tagMap
}

// assembleTransactionItems 组装交易记录数据
func (l *TransactionListLogic) assembleTransactionItems(list []*model.Transactions, tagMap map[uint64]*model.Tags) []types.TransactionItem {
	respList := make([]types.TransactionItem, 0, len(list))

	for _, item := range list {
		tItem := types.TransactionItem{}

		// 手动处理 ID 和 类型不匹配字段
		tItem.TransactionId = cast.ToString(item.TransactionId)
		tItem.UserId = cast.ToString(item.UserId)
		tItem.BookId = cast.ToInt64(item.BookId)
		tItem.Type = item.Type
		tItem.Amount = item.Amount
		tItem.TagId = cast.ToInt64(item.TagId)
		tItem.Remark = item.Remark
		tItem.TransactionTime = cast.ToInt64(item.TransactionTime)
		tItem.CreatedAt = cast.ToInt64(item.CreatedAt)
		tItem.UpdatedAt = cast.ToInt64(item.UpdatedAt)
		tItem.IsAutoGenerated = cast.ToInt64(item.IsAutoGenerated)

		// 填充标签信息
		if tag, ok := tagMap[item.TagId]; ok {
			tItem.TagName = tag.Name
			tItem.TagColor = tag.Color
			tItem.TagIcon = tag.Icon
		}

		respList = append(respList, tItem)
	}

	return respList
}
