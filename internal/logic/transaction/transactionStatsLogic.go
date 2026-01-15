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
	var totalIncome float64
	// 如果筛选条件指定了 type=expense，则收入为 0；否则计算收入
	if req.Type == "" || req.Type == "income" {
		incomeWhere := baseWhere
		if req.Type == "" {
			incomeWhere = append(incomeWhere, squirrel.Eq{"type": "income"})
		}
		// 如果 req.Type == "income"，baseWhere 已经包含了 type=income 条件
		sumBuilder := l.svcCtx.TransactionsModel.SumBuilder("amount").Where(incomeWhere)
		totalIncome, err = l.svcCtx.TransactionsModel.FindSum(l.ctx, sumBuilder)
		if err != nil {
			l.Logger.Errorf("TransactionStats FindSum Income error: %v", err)
			return nil, errcode.DBError.Msgr("统计收入失败")
		}
	}

	// 3. 计算总支出
	var totalExpense float64
	// 如果筛选条件指定了 type=income，则支出为 0；否则计算支出
	if req.Type == "" || req.Type == "expense" {
		expenseWhere := baseWhere
		if req.Type == "" {
			expenseWhere = append(expenseWhere, squirrel.Eq{"type": "expense"})
		}
		// 如果 req.Type == "expense"，baseWhere 已经包含了 type=expense 条件
		sumBuilder := l.svcCtx.TransactionsModel.SumBuilder("amount").Where(expenseWhere)
		totalExpense, err = l.svcCtx.TransactionsModel.FindSum(l.ctx, sumBuilder)
		if err != nil {
			l.Logger.Errorf("TransactionStats FindSum Expense error: %v", err)
			return nil, errcode.DBError.Msgr("统计支出失败")
		}
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

func (l *TransactionStatsLogic) calculateTagStats(where squirrel.Sqlizer) ([]types.TagStatItem, error) {
	// 构建分组查询
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

	if len(results) == 0 {
		return []types.TagStatItem{}, nil
	}

	// 获取标签详情
	tagIDs := make([]uint64, 0, len(results))
	for _, r := range results {
		tagIDs = append(tagIDs, r.TagId)
	}

	tagMap := make(map[uint64]struct {
		Name  string
		Color string
		Icon  string
	})

	if len(tagIDs) > 0 {
		tagBuilder := l.svcCtx.TagsModel.SelectBuilder().Where(squirrel.Eq{"tag_id": tagIDs})
		tags, err := l.svcCtx.TagsModel.FindAll(l.ctx, tagBuilder)
		if err == nil {
			for _, t := range tags {
				tagMap[t.TagId] = struct {
					Name  string
					Color string
					Icon  string
				}{
					Name:  t.Name,
					Color: t.Color,
					Icon:  t.Icon,
				}
			}
		} else {
			// 记录日志但不中断，允许返回无标签详情的数据
			l.Logger.Errorf("TransactionStats FindAll Tags error: %v", err)
		}
	}

	// 组装结果
	respStats := make([]types.TagStatItem, 0, len(results))
	for _, r := range results {
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

	return respStats, nil
}
