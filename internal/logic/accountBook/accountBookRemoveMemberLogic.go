package accountBook

import (
	"context"
	"strconv"

	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type AccountBookRemoveMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 移除账本成员
func NewAccountBookRemoveMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountBookRemoveMemberLogic {
	return &AccountBookRemoveMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccountBookRemoveMemberLogic) AccountBookRemoveMember(req *types.AccountBookRemoveMemberReq) (resp *types.AccountBookRemoveMemberResp, err error) {
	operatorId := cast.ToUint64(ctxData.GetUIDFromCtx(l.ctx))
	targetUserId, err := strconv.ParseUint(req.UserId, 10, 64)
	if err != nil {
		return nil, errcode.Fail.Msgr("无效的用户ID").WithError(errors.Wrap(err, "解析用户ID失败"))
	}

	// 2. Get Book
	book, err := l.svcCtx.AccountBooksModel.FindOne(l.ctx, uint64(req.BookId))
	if err != nil {
		return nil, errcode.DBError.WithError(errors.Wrap(err, "查询账本失败"))
	}

	// 3. Check permission (Operator must be creator)
	if book.CreatorUserId != operatorId {
		return nil, errcode.Forbidden.Msgr("只有账本创建者可以移除成员")
	}

	// 4. Check target validity
	// Cannot remove self (use exit)
	if targetUserId == operatorId {
		return nil, errcode.Fail.Msgr("不能移除自己，请使用退出账本功能")
	}

	// Check if target is member
	_, err = l.svcCtx.AccountBookMembersModel.FindOneByBookIdUserId(l.ctx, uint64(req.BookId), targetUserId)
	if err != nil {
		return nil, errcode.Fail.Msgr("该用户不是账本成员")
	}

	// 5. Transaction
	err = l.svcCtx.AccountBooksModel.Transaction(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// Delete member
		_, err := l.svcCtx.AccountBookMembersModel.DeleteFilter(ctx, session, squirrel.Eq{
			"book_id": req.BookId,
			"user_id": targetUserId,
		})
		if err != nil {
			return errors.Wrap(err, "删除账本成员失败")
		}

		// Update member count
		_, err = l.svcCtx.AccountBooksModel.UpdateFilter(ctx, session,
			map[string]interface{}{
				"member_count": squirrel.Expr("member_count - 1"),
			},
			squirrel.Eq{"book_id": req.BookId},
		)
		return errors.Wrap(err, "更新账本成员数量失败")
	})

	if err != nil {
		return nil, errcode.DBError.Msgr("移除成员失败").WithError(err)
	}

	return nil, nil
}
