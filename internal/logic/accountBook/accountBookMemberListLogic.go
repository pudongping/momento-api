package accountBook

import (
	"context"

	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

type AccountBookMemberListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取账本成员列表
func NewAccountBookMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountBookMemberListLogic {
	return &AccountBookMemberListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccountBookMemberListLogic) AccountBookMemberList(req *types.AccountBookMemberListReq) (resp []types.AccountBookMemberResp, err error) {
	userId := cast.ToUint64(ctxData.GetUIDFromCtx(l.ctx))

	// 2. 检查当前用户是否在账本中
	_, err = l.svcCtx.AccountBookMembersModel.FindOneByBookIdUserId(l.ctx, uint64(req.BookId), userId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errcode.Forbidden.Msgr("您不是该账本成员，无权查看")
		}
		return nil, errcode.DBError.Msgr("查询账本成员失败").WithError(errors.Wrap(err, "查询账本成员失败"))
	}

	// 3. 查询账本成员列表
	type MemberInfo struct {
		UserId    uint64 `db:"user_id"`
		IsCreator int64  `db:"is_creator"`
		Status    string `db:"status"`
		JoinedAt  uint64 `db:"joined_at"`
		Nickname  string `db:"nickname"`
		Avatar    string `db:"avatar"`
	}

	sb := squirrel.Select("m.user_id, m.is_creator, m.status, m.joined_at, u.nickname, u.avatar").
		From("account_book_members m").
		LeftJoin("users u ON m.user_id = u.user_id").
		Where(squirrel.Eq{"m.book_id": req.BookId})

	var members []MemberInfo
	err = l.svcCtx.AccountBookMembersModel.FindAny(l.ctx, sb, &members)
	if err != nil {
		return nil, errcode.DBError.Msgr("查询账本成员列表失败").WithError(errors.Wrap(err, "查询账本成员列表失败"))
	}

	resp = make([]types.AccountBookMemberResp, 0, len(members))
	for _, m := range members {
		resp = append(resp, types.AccountBookMemberResp{
			UserId:    cast.ToString(m.UserId),
			Nickname:  m.Nickname,
			Avatar:    m.Avatar,
			IsCreator: m.IsCreator,
			Status:    m.Status,
			JoinedAt:  int64(m.JoinedAt),
		})
	}

	return resp, nil
}
