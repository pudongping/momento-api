package accountBook

import (
	"context"

	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

type AccountBookExitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 退出账本
func NewAccountBookExitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountBookExitLogic {
	return &AccountBookExitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccountBookExitLogic) AccountBookExit(req *types.AccountBookExitReq) (resp *types.AccountBookExitResp, err error) {
	userID := cast.ToUint64(ctxData.GetUIDFromCtx(l.ctx))

	// 1. 查询成员信息
	member, err := l.svcCtx.AccountBookMembersModel.FindOneByBookIdUserId(l.ctx, uint64(req.BookId), userID)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errcode.NotFound.Msgr("您不是该账本成员")
		}
		return nil, errcode.DBError.Msgr("查询成员信息失败")
	}

	// 2. 检查是否是创建者
	if member.IsCreator == model.AccountBookMemberIsCreatorYes {
		return nil, errcode.Forbidden.Msgr("创建者不能退出账本，请使用删除账本功能")
	}

	// 3. 删除成员关联
	if err := l.svcCtx.AccountBookMembersModel.Delete(l.ctx, member.BookMemberId); err != nil {
		l.Logger.Errorf("Delete member error: %v", err)
		return nil, errcode.DBError.Msgr("退出账本失败")
	}

	return &types.AccountBookExitResp{}, nil
}
