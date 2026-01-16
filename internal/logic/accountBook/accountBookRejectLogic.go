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

type AccountBookRejectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 拒绝邀请
func NewAccountBookRejectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountBookRejectLogic {
	return &AccountBookRejectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccountBookRejectLogic) AccountBookReject(req *types.AccountBookRejectReq) (resp *types.AccountBookRejectResp, err error) {
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

	// 3. 更新邀请状态
	invitation.Status = model.AccountBookInvitationStatusRejected
	invitation.UpdatedAt = uint64(time.Now().Unix())
	if err := l.svcCtx.AccountBookInvitationsModel.Update(l.ctx, invitation); err != nil {
		l.Logger.Errorf("Update invitation status error: %v", err)
		return nil, errcode.DBError.Msgr("更新邀请状态失败")
	}

	return &types.AccountBookRejectResp{}, nil
}
