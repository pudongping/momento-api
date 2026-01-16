package accountBook

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type AccountBookSetDefaultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 设置默认账本
func NewAccountBookSetDefaultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountBookSetDefaultLogic {
	return &AccountBookSetDefaultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccountBookSetDefaultLogic) AccountBookSetDefault(req *types.AccountBookSetDefaultReq) (resp *types.AccountBookSetDefaultResp, err error) {
	userID := cast.ToUint64(ctxData.GetUIDFromCtx(l.ctx))

	// 1. 确认用户是该账本成员
	member, err := l.svcCtx.AccountBookMembersModel.FindOneByBookIdUserId(l.ctx, uint64(req.BookId), userID)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errcode.NotFound.Msgr("您不是该账本成员")
		}
		return nil, errcode.DBError.Msgr("查询成员信息失败").WithError(errors.Wrap(err, "查询成员信息失败"))
	}

	// 2. 如果已经是默认账本，直接返回
	if member.IsDefault == model.AccountBookMemberIsDefaultYes {
		return nil, nil
	}

	// 3. 开启事务更新
	err = l.svcCtx.AccountBookMembersModel.Transaction(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		now := time.Now().Unix()

		// Step 3.1: 将该用户的所有账本 is_default 设为 No
		updateAllData := map[string]interface{}{
			"is_default": model.AccountBookMemberIsDefaultNo,
			"updated_at": now,
		}
		whereAll := squirrel.Eq{"user_id": userID}

		if _, err := l.svcCtx.AccountBookMembersModel.UpdateFilter(ctx, session, updateAllData, whereAll); err != nil {
			return err
		}

		// Step 3.2: 将当前账本设为默认
		updateCurrentData := map[string]interface{}{
			"is_default": model.AccountBookMemberIsDefaultYes,
			"updated_at": now,
		}
		whereCurrent := squirrel.Eq{"book_member_id": member.BookMemberId}

		if _, err := l.svcCtx.AccountBookMembersModel.UpdateFilter(ctx, session, updateCurrentData, whereCurrent); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, errcode.DBError.Msgr("设置默认账本失败").WithError(errors.Wrap(err, "设置默认账本失败"))
	}

	return nil, nil
}
