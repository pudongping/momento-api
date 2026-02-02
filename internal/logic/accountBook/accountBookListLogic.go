package accountBook

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

type AccountBookListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取账本列表
func NewAccountBookListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountBookListLogic {
	return &AccountBookListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type accountBookListRow struct {
	BookId        int64  `db:"book_id"`
	Name          string `db:"name"`
	CreatorUserId int64  `db:"creator_user_id"`
	IsCreator     int64  `db:"is_creator"`
	IsDefault     int64  `db:"is_default"`
	CreatedAt     int64  `db:"created_at"`
}

func (l *AccountBookListLogic) AccountBookList(req *types.AccountBookListReq) (resp []types.AccountBookListResp, err error) {
	userId := ctxData.GetUIDFromCtx(l.ctx)

	// 构建查询
	// SELECT ab.book_id, ab.name, ab.creator_user_id, ab.created_at, abm.is_creator, abm.is_default
	// FROM account_books ab
	// JOIN account_book_members abm ON ab.book_id = abm.book_id
	// WHERE abm.user_id = ?
	// ORDER BY abm.is_default ASC, ab.created_at DESC

	builder := squirrel.Select(
		"ab.book_id",
		"ab.name",
		"ab.creator_user_id",
		"ab.created_at",
		"abm.is_creator",
		"abm.is_default",
	).
		From("account_books ab").
		Join("account_book_members abm ON ab.book_id = abm.book_id").
		Where(squirrel.Eq{"abm.user_id": userId}).
		OrderBy("abm.is_default ASC", "ab.created_at DESC")

	var rows []accountBookListRow
	// 使用任意 Model 执行查询即可，因为是 Join 查询，不局限于特定 Table
	err = l.svcCtx.AccountBooksModel.FindAny(l.ctx, builder, &rows)
	if err != nil {
		return nil, errcode.DBError.WithError(err).Msgr("查询账本列表失败")
	}

	if len(rows) == 0 {
		return []types.AccountBookListResp{}, nil
	}

	// 统计成员数量（只统计已加入的成员）
	bookIDs := make([]int64, len(rows))
	for i, row := range rows {
		bookIDs[i] = row.BookId
	}

	countBuilder := l.svcCtx.AccountBookMembersModel.SelectBuilder("book_id", "COUNT(*) as count").
		Where(squirrel.Eq{
			"book_id": bookIDs,
			"status":  model.AccountBookMemberStatusJoined,
		}).
		GroupBy("book_id")

	var countRows []struct {
		BookId int64 `db:"book_id"`
		Count  int64 `db:"count"`
	}
	err = l.svcCtx.AccountBookMembersModel.FindAny(l.ctx, countBuilder, &countRows)
	if err != nil {
		return nil, errcode.DBError.WithError(err).Msgr("查询成员数量失败")
	}

	countMap := make(map[int64]int64)
	for _, r := range countRows {
		countMap[r.BookId] = r.Count
	}

	resp = make([]types.AccountBookListResp, 0, len(rows))
	for _, row := range rows {
		resp = append(resp, types.AccountBookListResp{
			BookId:        row.BookId,
			Name:          row.Name,
			CreatorUserId: cast.ToString(row.CreatorUserId),
			IsCreator:     row.IsCreator,
			IsDefault:     row.IsDefault,
			MemberCount:   countMap[row.BookId],
			CreatedAt:     row.CreatedAt,
		})
	}

	return resp, nil
}
