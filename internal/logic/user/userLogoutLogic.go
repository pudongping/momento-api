package user

import (
	"context"

	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/constant"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户退出登录
func NewUserLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogoutLogic {
	return &UserLogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLogoutLogic) UserLogout(req *types.UserLogoutReq, token string) (*types.UserLogoutResp, error) {
	userID := ctxData.GetUIDFromCtx(l.ctx)

	if token == "" {
		return nil, errcode.Fail.Msgr("Token不能为空")
	}

	// 从 Redis 中删除 token
	cacheKey := constant.RedisKeyPrefixToken + token
	_, err := l.svcCtx.RedisClient.DelCtx(l.ctx, cacheKey)
	if err != nil {
		l.Logger.Errorf("用户退出登录失败，删除 Redis token 失败 userID : %d, err : %v", userID, err)
		return nil, errcode.InternalServerError.WithError(errors.Wrapf(err, "UserLogout DelCtx userID : %d", userID)).Msgr("退出登录失败")
	}

	return &types.UserLogoutResp{}, nil
}
