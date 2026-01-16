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

type AccountBookAcceptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 接受邀请
func NewAccountBookAcceptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountBookAcceptLogic {
	return &AccountBookAcceptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccountBookAcceptLogic) AccountBookAccept(req *types.AccountBookAcceptReq) (resp *types.AccountBookAcceptResp, err error) {
	userID := cast.ToUint64(ctxData.GetUIDFromCtx(l.ctx))
	now := time.Now().Unix()

	// 1. 查询邀请记录
	invitation, err := l.svcCtx.AccountBookInvitationsModel.FindOne(l.ctx, uint64(req.InvitationId))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errcode.NotFound.Msgf("邀请记录")
		}
		return nil, errcode.DBError.Msgr("查询邀请记录失败")
	}

	// 2. 权限检查
	if invitation.TargetUid != userID {
		return nil, errcode.Forbidden.Msgr("无权处理该邀请")
	}

	if invitation.Status != model.AccountBookInvitationStatusPending {
		return nil, errcode.BadRequest.Msgr("该邀请已处理")
	}

	// 3. 检查是否已经是成员
	_, err = l.svcCtx.AccountBookMembersModel.FindOneByBookIdUserId(l.ctx, invitation.BookId, userID)
	if err == nil {
		// 已经是成员，直接更新邀请状态为已接受
		updateData := map[string]interface{}{
			"status":     model.AccountBookInvitationStatusAccepted,
			"updated_at": now,
		}
		where := squirrel.Eq{"invitation_id": invitation.InvitationId}

		// 显式传 nil session，表示不使用事务
		if _, err := l.svcCtx.AccountBookInvitationsModel.UpdateFilter(l.ctx, nil, updateData, where); err != nil {
			return nil, errcode.DBError.Msgr("更新邀请状态失败").WithError(errors.Wrap(err, "更新邀请状态失败"))
		}
		return &types.AccountBookAcceptResp{}, nil
	}
	if err != model.ErrNotFound {
		return nil, errcode.DBError.Msgr("查询成员信息失败")
	}

	// 4. 开启事务：更新邀请状态 + 加入账本
	// 使用 AccountBookInvitationsModel 开启事务（也可以用其他 Model，只要共享连接）
	err = l.svcCtx.AccountBookInvitationsModel.Transaction(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 4.1 更新邀请状态
		updateData := map[string]interface{}{
			"status":     model.AccountBookInvitationStatusAccepted,
			"updated_at": now,
		}
		where := squirrel.Eq{"invitation_id": invitation.InvitationId}

		if _, err := l.svcCtx.AccountBookInvitationsModel.UpdateFilter(ctx, session, updateData, where); err != nil {
			return err
		}

		// 4.2 加入账本
		memberData := map[string]interface{}{
			"book_id":    invitation.BookId,
			"user_id":    userID,
			"is_creator": model.AccountBookMemberIsCreatorNo,
			"is_default": model.AccountBookMemberIsDefaultNo,
			"status":     model.AccountBookMemberStatusJoined,
			"joined_at":  now,
			"created_at": now,
			"updated_at": now,
		}

		if _, err := l.svcCtx.AccountBookMembersModel.InsertWithSession(ctx, session, memberData); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, errcode.DBError.Msgr("加入账本失败").WithError(errors.Wrap(err, "加入账本失败"))
	}

	return &types.AccountBookAcceptResp{}, nil
}
