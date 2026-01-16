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
)

type AccountBookInviteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 邀请用户加入账本
func NewAccountBookInviteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountBookInviteLogic {
	return &AccountBookInviteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccountBookInviteLogic) AccountBookInvite(req *types.AccountBookInviteReq) (resp *types.AccountBookInviteResp, err error) {
	userID := cast.ToUint64(ctxData.GetUIDFromCtx(l.ctx))
	now := time.Now().Unix()

	targetUID := cast.ToUint64(req.TargetUid)
	if targetUID == 0 {
		return nil, errcode.UnprocessableEntity.Msgf("被邀请人ID无效")
	}

	if userID == targetUID {
		return nil, errcode.UnprocessableEntity.Msgf("不能邀请自己")
	}

	// 2. 检查被邀请人是否存在
	_, err = l.svcCtx.UserModel.FindOne(l.ctx, targetUID)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errcode.NotFound.Msgr("被邀请人不存在")
		}
		return nil, errcode.DBError.Msgr("查询被邀请人失败")
	}

	// 3. 检查账本是否存在
	_, err = l.svcCtx.AccountBooksModel.FindOne(l.ctx, cast.ToUint64(req.BookId))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errcode.NotFound.Msgr("账本不存在")
		}
		return nil, errcode.DBError.Msgr("查询账本失败")
	}

	// 4. 检查邀请人是否为账本成员
	_, err = l.svcCtx.AccountBookMembersModel.FindOneByBookIdUserId(l.ctx, uint64(req.BookId), userID)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errcode.Forbidden.Msgr("您不是该账本成员，无权邀请")
		}
		return nil, errcode.DBError.Msgr("查询账本成员失败")
	}

	// 5. 检查被邀请人是否已经是账本成员
	_, err = l.svcCtx.AccountBookMembersModel.FindOneByBookIdUserId(l.ctx, uint64(req.BookId), targetUID)
	if err == nil {
		return nil, errcode.BadRequest.Msgr("该用户已经是账本成员")
	}
	if err != model.ErrNotFound {
		return nil, errcode.DBError.Msgr("查询账本成员失败")
	}

	// 6. 检查是否已经邀请过（待处理状态）
	rowBuilder := l.svcCtx.AccountBookInvitationsModel.SelectBuilder().
		Where(squirrel.Eq{
			"book_id":    req.BookId,
			"target_uid": targetUID,
			"status":     model.AccountBookInvitationStatusPending,
		})
	invitation, err := l.svcCtx.AccountBookInvitationsModel.FindOneByQuery(l.ctx, rowBuilder)
	if err != nil && err != model.ErrNotFound {
		return nil, errcode.DBError.Msgr("查询邀请记录失败")
	}
	if invitation != nil {
		return nil, errcode.BadRequest.Msgr("已邀请过该用户，请等待对方处理")
	}

	// 7. 创建邀请记录
	newInvitation := &model.AccountBookInvitations{
		BookId:     uint64(req.BookId),
		InviterUid: userID,
		TargetUid:  targetUID,
		Status:     model.AccountBookInvitationStatusPending,
		CreatedAt:  uint64(now),
		UpdatedAt:  uint64(now),
	}

	_, err = l.svcCtx.AccountBookInvitationsModel.Insert(l.ctx, newInvitation)
	if err != nil {
		return nil, errcode.DBError.Msgr("邀请失败").WithError(errors.Wrap(err, "AccountBookInvite Insert error"))
	}

	return &types.AccountBookInviteResp{}, nil
}
