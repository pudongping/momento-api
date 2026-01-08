// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户信息
func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoReq) (*types.UserInfoResp, error) {
	userID := ctxData.GetUIDFromCtx(l.ctx)
	if userID <= 0 {
		return nil, errcode.Fail.Msgr("用户未登录")
	}

	queryBuilder := l.svcCtx.UserModel.SelectBuilder().
		Where("user_id = ?", cast.ToUint64(userID)).
		Where("is_disable = ?", model.UserIsDisableNo)

	user, err := l.svcCtx.UserModel.FindOneByQuery(l.ctx, queryBuilder)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errcode.DBError.WithError(errors.Wrapf(err, "UserInfo FindOne")).Msgr("获取用户信息失败")
	}
	if user == nil {
		return nil, errcode.Fail.Msgr("当前用户不存在")
	}

	var resp types.UserInfoResp
	if err = copier.Copy(&resp, user); err != nil {
		return nil, errors.Wrap(err, "UserInfo copier")
	}

	return &resp, nil
}
