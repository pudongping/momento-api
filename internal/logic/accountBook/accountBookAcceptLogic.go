package accountBook

import (
	"context"
	"time"

	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
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

	// 1. 查询邀请记录
	invitation, err := l.svcCtx.AccountBookInvitationsModel.FindOne(l.ctx, uint64(req.InvitationId))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errcode.NotFound.Msgr("邀请记录不存在")
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
		invitation.Status = model.AccountBookInvitationStatusAccepted
		invitation.UpdatedAt = uint64(time.Now().Unix())
		if err := l.svcCtx.AccountBookInvitationsModel.Update(l.ctx, invitation); err != nil {
			l.Logger.Errorf("Update invitation status error: %v", err)
			return nil, errcode.DBError.Msgr("更新邀请状态失败")
		}
		return &types.AccountBookAcceptResp{}, nil
	}
	if err != model.ErrNotFound {
		return nil, errcode.DBError.Msgr("查询成员信息失败")
	}

	// 4. 更新邀请状态
	invitation.Status = model.AccountBookInvitationStatusAccepted
	invitation.UpdatedAt = uint64(time.Now().Unix())
	if err := l.svcCtx.AccountBookInvitationsModel.Update(l.ctx, invitation); err != nil {
		l.Logger.Errorf("Update invitation status error: %v", err)
		return nil, errcode.DBError.Msgr("更新邀请状态失败")
	}

	// 5. 加入账本
	member := &model.AccountBookMembers{
		BookId:    invitation.BookId,
		UserId:    userID,
		IsCreator: model.AccountBookMemberIsCreatorNo,
		IsDefault: model.AccountBookMemberIsDefaultNo,
		Status:    model.AccountBookMemberStatusJoined,
		CreatedAt: uint64(time.Now().Unix()),
		UpdatedAt: uint64(time.Now().Unix()),
	}
	if _, err := l.svcCtx.AccountBookMembersModel.Insert(l.ctx, member); err != nil {
		l.Logger.Errorf("Insert member error: %v", err)
		// 注意：这里没有事务，如果插入失败，邀请状态已经变成了 accepted，可能需要回滚或者人工介入
		// 考虑到是个人项目，暂不引入复杂事务处理
		return nil, errcode.DBError.Msgr("加入账本失败")
	}

	return &types.AccountBookAcceptResp{}, nil
}
