// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package accountBook

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type AccountBookDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除账本
func NewAccountBookDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountBookDeleteLogic {
	return &AccountBookDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccountBookDeleteLogic) AccountBookDelete(req *types.AccountBookDeleteReq) (resp *types.AccountBookDeleteResp, err error) {
	userId := ctxData.GetUIDFromCtx(l.ctx)
	userIdUint := cast.ToUint64(userId)
	bookIdUint := cast.ToUint64(req.BookId)

	// 1. Check if book exists and user is creator
	book, err := l.svcCtx.AccountBooksModel.FindOne(l.ctx, bookIdUint)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, errcode.NotFound.Msgr("账本不存在")
		}
		return nil, errcode.DBError.WithError(err).Msgr("查询账本失败")
	}

	if book.CreatorUserId != userIdUint {
		return nil, errcode.Forbidden.Msgr("只有创建者可以删除账本")
	}

	// 2. Cascade Delete in Transaction
	err = l.svcCtx.AccountBooksModel.Transaction(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// Delete Transactions
		_, err = l.svcCtx.TransactionsModel.DeleteFilter(ctx, session, squirrel.Eq{"book_id": bookIdUint})
		if err != nil {
			return errcode.DBError.WithError(errors.Wrap(err, "删除账单记录失败")).Msgr("删除账单记录失败")
		}

		// Delete Recurring Transactions
		_, err = l.svcCtx.RecurringTransactionsModel.DeleteFilter(ctx, session, squirrel.Eq{"book_id": bookIdUint})
		if err != nil {
			return errcode.DBError.WithError(errors.Wrap(err, "删除周期账单记录失败")).Msgr("删除周期账单记录失败")
		}

		// Delete Account Book Members
		_, err = l.svcCtx.AccountBookMembersModel.DeleteFilter(ctx, session, squirrel.Eq{"book_id": bookIdUint})
		if err != nil {
			return errcode.DBError.WithError(errors.Wrap(err, "删除账本成员失败")).Msgr("删除账本成员失败")
		}

		// Delete Account Book Invitations
		_, err = l.svcCtx.AccountBookInvitationsModel.DeleteFilter(ctx, session, squirrel.Eq{"book_id": bookIdUint})
		if err != nil {
			return errcode.DBError.WithError(errors.Wrap(err, "删除账本邀请失败")).Msgr("删除账本邀请失败")
		}

		// Delete Account Book
		_, err = l.svcCtx.AccountBooksModel.DeleteFilter(ctx, session, squirrel.Eq{"book_id": bookIdUint})
		if err != nil {
			return errcode.DBError.WithError(errors.Wrap(err, "删除账本失败")).Msgr("删除账本失败")
		}

		return nil
	})

	if err != nil {
		return nil, errcode.DBError.WithError(err).Msgr("删除账本失败")
	}

	return &types.AccountBookDeleteResp{}, nil
}
