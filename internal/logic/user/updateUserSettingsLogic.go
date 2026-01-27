package user

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

type UpdateUserSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新用户设置
func NewUpdateUserSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserSettingsLogic {
	return &UpdateUserSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserSettingsLogic) UpdateUserSettings(req *types.UpdateUserSettingsReq) error {
	userID := ctxData.GetUIDFromCtx(l.ctx)
	if userID <= 0 {
		return errcode.Fail.Msgr("用户未登录")
	}

	// 查找用户设置
	userSettings, err := l.svcCtx.UserSettingModel.FindOneByUserId(l.ctx, cast.ToUint64(userID))
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return errcode.DBError.WithError(errors.Wrapf(err, "UpdateUserSettings FindOneByUserId")).Msgr("获取用户设置失败")
	}

	now := uint64(time.Now().Unix())

	if userSettings == nil {
		// 不存在则创建
		userSettings = &model.UserSettings{
			UserId:    cast.ToUint64(userID),
			CreatedAt: now,
			UpdatedAt: now,
		}

		if req.BackgroundUrl != nil {
			userSettings.BackgroundUrl = *req.BackgroundUrl
		}
		if req.Budget != nil {
			userSettings.Budget = *req.Budget
		}

		_, err = l.svcCtx.UserSettingModel.Insert(l.ctx, userSettings)
		if err != nil {
			return errcode.DBError.WithError(errors.Wrapf(err, "UpdateUserSettings Insert")).Msgr("创建用户设置失败")
		}
	} else {
		// 存在则更新
		updateData := make(map[string]interface{})

		if req.BackgroundUrl != nil && *req.BackgroundUrl != userSettings.BackgroundUrl {
			updateData["background_url"] = *req.BackgroundUrl
		}
		if req.Budget != nil && *req.Budget != userSettings.Budget {
			updateData["budget"] = *req.Budget
		}

		if len(updateData) > 0 {
			updateData["updated_at"] = now
			where := squirrel.Eq{"user_setting_id": userSettings.UserSettingId}
			_, err = l.svcCtx.UserSettingModel.UpdateFilter(l.ctx, nil, updateData, where)
			if err != nil {
				return errcode.DBError.WithError(errors.Wrapf(err, "UpdateUserSettings UpdateFilter")).Msgr("更新用户设置失败")
			}
		}
	}

	return nil
}
