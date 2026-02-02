package transaction

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
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

	if err = l.checkBookMember(req.BookId, userID); err != nil {
		return nil, err
	}

	// 1. 构建基础查询条件
	whereBuilder := l.buildBaseQuery(req)

	// 2. 初始化查询构建器
	// 显式指定表名以避免连接查询时的列名冲突
	countBuilder := l.svcCtx.TransactionsModel.CountBuilder("transactions.transaction_id")
	selectBuilder := l.svcCtx.TransactionsModel.SelectBuilder("transactions.*")

	// 如果有关键词搜索，需要关联 tags 表
	if req.Keyword != "" {
		joinClause := "tags ON transactions.tag_id = tags.tag_id"
		countBuilder = countBuilder.LeftJoin(joinClause)
		selectBuilder = selectBuilder.LeftJoin(joinClause)
	}

	// 3. 计算总数
	countBuilder = countBuilder.Where(whereBuilder)
	total, err := l.countTotal(countBuilder)
	if err != nil {
		return nil, errcode.DBError.WithError(errors.Wrap(err, "计算交易流水总数失败"))
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

	// 4. 构建列表查询并应用游标分页
	selectBuilder = selectBuilder.Where(whereBuilder)
	if req.LastTransactionId > 0 {
		selectBuilder = l.applyCursorPagination(selectBuilder, req.LastTransactionId)
	}

	// 5. 设置排序和分页限制
	// 默认只按 transaction_id 倒序排列
	selectBuilder = selectBuilder.OrderBy("transactions.transaction_id DESC").Limit(uint64(limit + 1))

	// 6. 执行查询
	list, err := l.svcCtx.TransactionsModel.FindAll(l.ctx, selectBuilder)
	if err != nil {
		l.Logger.Errorf("TransactionList FindAll error: %v", err)
		return nil, errcode.DBError.WithError(err)
	}

	// 7. 处理分页标志
	hasMore := false
	if len(list) > int(limit) {
		hasMore = true
		list = list[:limit]
	}

	// 8. 获取标签信息
	tagMap := l.getTagsMap(list)

	// 9. 组装响应数据
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
func (l *TransactionListLogic) buildBaseQuery(req *types.TransactionListReq) squirrel.Sqlizer {
	conditions := squirrel.And{
		squirrel.Eq{"transactions.book_id": req.BookId},
	}

	if req.Type != "" {
		conditions = append(conditions, squirrel.Eq{"transactions.type": req.Type})
	}
	if req.MinAmount > 0 {
		conditions = append(conditions, squirrel.GtOrEq{"transactions.amount": req.MinAmount})
	}
	if req.MaxAmount > 0 {
		conditions = append(conditions, squirrel.LtOrEq{"transactions.amount": req.MaxAmount})
	}
	if req.StartDate > 0 {
		conditions = append(conditions, squirrel.GtOrEq{"transactions.transaction_time": req.StartDate})
	}
	if req.EndDate > 0 {
		conditions = append(conditions, squirrel.LtOrEq{"transactions.transaction_time": req.EndDate})
	}
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		conditions = append(conditions, squirrel.Or{
			squirrel.Like{"transactions.remark": keyword},
			squirrel.Like{"tags.name": keyword},
		})
	}

	return conditions
}

// countTotal 计算总数
func (l *TransactionListLogic) countTotal(countBuilder squirrel.SelectBuilder) (int64, error) {
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
		return builder.Where(squirrel.Lt{"transactions.transaction_id": lastID})
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

// checkBookMember 检查用户是否为账本成员且状态正常
func (l *TransactionListLogic) checkBookMember(bookID int64, userID uint64) error {
	member, err := l.svcCtx.AccountBookMembersModel.FindOneByBookIdUserId(l.ctx, uint64(bookID), userID)
	if err != nil {
		if err == model.ErrNotFound {
			return errcode.Fail.Msgr("您不是该账本的成员，无法查看账单")
		}
		return errcode.DBError.WithError(errors.Wrap(err, "查询账本成员失败"))
	}

	// 检查成员状态是否为已加入
	if member.Status != model.AccountBookMemberStatusJoined {
		return errcode.Fail.Msgr("您尚未加入该账本，无法查看账单")
	}

	return nil
}
