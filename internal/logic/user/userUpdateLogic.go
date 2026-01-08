// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新用户信息
func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUpdateLogic) UserUpdate(req *types.UserUpdateReq) (*types.UserUpdateResp, error) {
	userID := ctxData.GetUIDFromCtx(l.ctx)

	// 仅更新非空字段
	updateMap := map[string]interface{}{}
	if s := strings.TrimSpace(req.Nickname); s != "" {
		updateMap["nickname"] = s
	}
	if s := strings.TrimSpace(req.Avatar); s != "" {
		updateMap["avatar"] = s
	}
	if s := strings.TrimSpace(req.Phone); s != "" {
		updateMap["phone"] = s
	}

	// 如果全部为空，则不做更新，直接返回当前信息
	if len(updateMap) == 0 {
		return nil, errcode.Fail.Msgr("没有可更新的用户信息")
	}

	// 补充更新时间
	now := time.Now().Unix()
	updateMap["updated_at"] = cast.ToUint64(now)

	where := squirrel.Eq{"user_id": cast.ToUint64(userID)}
	if _, err := l.svcCtx.UserModel.UpdateFilter(l.ctx, nil, updateMap, where); err != nil {
		l.Logger.Errorf("更新用户信息失败 userID : %d, err : %v", userID, err)
		return nil, errcode.DBError.WithError(errors.Wrapf(err, "UserUpdate UpdateFilter userID : %d", userID)).Msgr("更新用户信息失败")
	}

	// 返回最新数据
	userData, err := l.svcCtx.UserModel.FindOne(l.ctx, cast.ToUint64(userID))
	if err != nil {
		return nil, errcode.DBError.WithError(errors.Wrap(err, "UserUpdate FindOne")).Msgr("获取用户信息失败")
	}

	var resp types.UserUpdateResp
	if err = copier.Copy(&resp, userData); err != nil {
		return nil, errors.Wrap(err, "UserUpdate copier")
	}
	resp.UserId = cast.ToString(userData.UserId)

	return &resp, nil
}
