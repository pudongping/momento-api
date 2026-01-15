package transaction

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

type TransactionStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取交易统计
func NewTransactionStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransactionStatsLogic {
	return &TransactionStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransactionStatsLogic) TransactionStats(req *types.TransactionStatsReq) (resp *types.TransactionStatsResp, err error) {
	userID := cast.ToUint64(ctxData.GetUIDFromCtx(l.ctx))

	// 1. 构建基础查询条件
	baseWhere := l.buildBaseQuery(userID, req)

	// 2. 计算总收入
	totalIncome, err := l.calculateTotalIncome(req, baseWhere)
	if err != nil {
		return nil, err
	}

	// 3. 计算总支出
	totalExpense, err := l.calculateTotalExpense(req, baseWhere)
	if err != nil {
		return nil, err
	}

	// 4. 计算结余
	balance := totalIncome - totalExpense

	// 5. 计算各标签占比
	tagStats, err := l.calculateTagStats(baseWhere)
	if err != nil {
		l.Logger.Errorf("TransactionStats calculateTagStats error: %v", err)
		return nil, errcode.DBError.Msgr("统计标签数据失败")
	}

	return &types.TransactionStatsResp{
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		Balance:      balance,
		TagStats:     tagStats,
	}, nil
}

// calculateTotalIncome 计算总收入
func (l *TransactionStatsLogic) calculateTotalIncome(req *types.TransactionStatsReq, baseWhere squirrel.And) (float64, error) {
	// 如果筛选条件指定了 type=expense，则收入为 0
	if req.Type == "expense" {
		return 0, nil
	}

	// 复制 baseWhere 防止修改原切片
	incomeWhere := make(squirrel.And, len(baseWhere))
	copy(incomeWhere, baseWhere)

	// 如果没有指定 type，则限定为 income
	if req.Type == "" {
		incomeWhere = append(incomeWhere, squirrel.Eq{"type": "income"})
	}

	sumBuilder := l.svcCtx.TransactionsModel.SumBuilder("amount").Where(incomeWhere)
	totalIncome, err := l.svcCtx.TransactionsModel.FindSum(l.ctx, sumBuilder)
	if err != nil {
		l.Logger.Errorf("TransactionStats FindSum Income error: %v", err)
		return 0, errcode.DBError.Msgr("统计收入失败")
	}
	return totalIncome, nil
}

// calculateTotalExpense 计算总支出
func (l *TransactionStatsLogic) calculateTotalExpense(req *types.TransactionStatsReq, baseWhere squirrel.And) (float64, error) {
	// 如果筛选条件指定了 type=income，则支出为 0
	if req.Type == "income" {
		return 0, nil
	}

	// 复制 baseWhere 防止修改原切片
	expenseWhere := make(squirrel.And, len(baseWhere))
	copy(expenseWhere, baseWhere)

	// 如果没有指定 type，则限定为 expense
	if req.Type == "" {
		expenseWhere = append(expenseWhere, squirrel.Eq{"type": "expense"})
	}

	sumBuilder := l.svcCtx.TransactionsModel.SumBuilder("amount").Where(expenseWhere)
	totalExpense, err := l.svcCtx.TransactionsModel.FindSum(l.ctx, sumBuilder)
	if err != nil {
		l.Logger.Errorf("TransactionStats FindSum Expense error: %v", err)
		return 0, errcode.DBError.Msgr("统计支出失败")
	}
	return totalExpense, nil
}

// buildBaseQuery 构建通用的过滤条件
func (l *TransactionStatsLogic) buildBaseQuery(userID uint64, req *types.TransactionStatsReq) squirrel.And {
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

type TagStatDbResult struct {
	TagId  uint64  `db:"tag_id"`
	Count  int64   `db:"count"`
	Amount float64 `db:"amount"`
}

// TagInfo 内部使用的标签信息结构
type TagInfo struct {
	Name  string
	Color string
	Icon  string
}

// calculateTagStats 计算各标签占比
func (l *TransactionStatsLogic) calculateTagStats(where squirrel.Sqlizer) ([]types.TagStatItem, error) {
	// 1. 获取原始统计数据
	rawStats, err := l.getRawTagStats(where)
	if err != nil {
		return nil, err
	}

	if len(rawStats) == 0 {
		return []types.TagStatItem{}, nil
	}

	// 2. 获取标签详情
	tagIDs := l.extractTagIDs(rawStats)
	tagMap := l.getTagsMap(tagIDs)

	// 3. 组装结果
	return l.buildTagStatItems(rawStats, tagMap), nil
}

// getRawTagStats 从数据库获取原始统计数据
func (l *TransactionStatsLogic) getRawTagStats(where squirrel.Sqlizer) ([]TagStatDbResult, error) {
	// SELECT tag_id, COUNT(*) as count, SUM(amount) as amount FROM transactions WHERE ... GROUP BY tag_id ORDER BY amount DESC
	builder := l.svcCtx.TransactionsModel.SelectBuilder("tag_id", "COUNT(*) as count", "SUM(amount) as amount").
		Where(where).
		GroupBy("tag_id").
		OrderBy("amount DESC")

	var results []TagStatDbResult
	err := l.svcCtx.TransactionsModel.FindAny(l.ctx, builder, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// extractTagIDs 提取标签ID
func (l *TransactionStatsLogic) extractTagIDs(rawStats []TagStatDbResult) []uint64 {
	tagIDs := make([]uint64, 0, len(rawStats))
	for _, r := range rawStats {
		tagIDs = append(tagIDs, r.TagId)
	}
	return tagIDs
}

// getTagsMap 获取标签映射信息
func (l *TransactionStatsLogic) getTagsMap(tagIDs []uint64) map[uint64]TagInfo {
	tagMap := make(map[uint64]TagInfo)
	if len(tagIDs) == 0 {
		return tagMap
	}

	tagBuilder := l.svcCtx.TagsModel.SelectBuilder().Where(squirrel.Eq{"tag_id": tagIDs})
	tags, err := l.svcCtx.TagsModel.FindAll(l.ctx, tagBuilder)
	if err != nil {
		// 记录日志但不中断，允许返回无标签详情的数据
		l.Logger.Errorf("TransactionStats FindAll Tags error: %v", err)
		return tagMap
	}

	for _, t := range tags {
		tagMap[t.TagId] = TagInfo{
			Name:  t.Name,
			Color: t.Color,
			Icon:  t.Icon,
		}
	}
	return tagMap
}

// buildTagStatItems 组装最终结果
func (l *TransactionStatsLogic) buildTagStatItems(rawStats []TagStatDbResult, tagMap map[uint64]TagInfo) []types.TagStatItem {
	respStats := make([]types.TagStatItem, 0, len(rawStats))
	for _, r := range rawStats {
		item := types.TagStatItem{
			TagId:  int64(r.TagId),
			Count:  r.Count,
			Amount: r.Amount,
		}

		if info, ok := tagMap[r.TagId]; ok {
			item.TagName = info.Name
			item.TagColor = info.Color
			item.TagIcon = info.Icon
		} else {
			// 默认值或处理已删除标签
			item.TagName = "未知标签"
		}

		respStats = append(respStats, item)
	}
	return respStats
}
