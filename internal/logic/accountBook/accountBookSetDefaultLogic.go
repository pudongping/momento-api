package accountBook

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
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
		return nil, errcode.DBError.Msgr("查询成员信息失败")
	}

	// 2. 如果已经是默认账本，直接返回
	if member.IsDefault == model.AccountBookMemberIsDefaultYes {
		return &types.AccountBookSetDefaultResp{}, nil
	}

	// 3. 更新操作
	// 将该用户的所有账本 is_default 设为 No
	// 这里使用 squirrel 构建 update 语句，或者手动执行 SQL
	// 由于 go-zero model 没有提供批量更新方法，这里我们可以分步或者自定义 Query

	// Step 3.1: Reset all user's books to non-default
	// 最好使用事务，这里使用顺序操作

	// 使用自定义 SQL 更新
	// UPDATE account_book_members SET is_default = 2 WHERE user_id = ?
	query, args, err := squirrel.Update("account_book_members").
		Set("is_default", model.AccountBookMemberIsDefaultNo).
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, errcode.InternalServerError.Msgr("构建SQL失败")
	}

	// Execute SQL
	_, err = l.svcCtx.MysqlClient.ExecCtx(l.ctx, query, args...)
	if err != nil {
		l.Logger.Errorf("Reset default book error: %v", err)
		return nil, errcode.DBError.Msgr("重置默认账本失败")
	}

	// Step 3.2: Set current book as default
	member.IsDefault = model.AccountBookMemberIsDefaultYes
	member.UpdatedAt = uint64(time.Now().Unix())
	if err := l.svcCtx.AccountBookMembersModel.Update(l.ctx, member); err != nil {
		l.Logger.Errorf("Set default book error: %v", err)
		return nil, errcode.DBError.Msgr("设置默认账本失败")
	}

	return &types.AccountBookSetDefaultResp{}, nil
}
