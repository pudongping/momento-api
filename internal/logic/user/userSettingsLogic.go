package user

import (
	"context"

	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户设置
func NewUserSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserSettingsLogic {
	return &UserSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserSettingsLogic) UserSettings(req *types.UserSettingsReq) (*types.UserSettingsResp, error) {
	userID := ctxData.GetUIDFromCtx(l.ctx)

	userSettings, err := l.svcCtx.UserSettingModel.FindOneByUserId(l.ctx, cast.ToUint64(userID))
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errcode.DBError.WithError(errors.Wrapf(err, "UserSettings FindOneByUserId")).Msgr("获取用户设置失败")
	}

	var resp types.UserSettingsResp

	if userSettings != nil {
		if err = copier.Copy(&resp, userSettings); err != nil {
			return nil, errors.Wrapf(err, "拷贝用户信息失败 userID : %d", userSettings.UserId)
		}
		// 返回到客户端的用户 ID 是一个字符串格式
		resp.UserId = cast.ToString(userSettings.UserId)
	}

	return &resp, nil
}
