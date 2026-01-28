package user

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/constant"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"
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

func (l *UserLogoutLogic) UserLogout(req *types.UserLogoutReq, token string) error {
	userID := ctxData.GetUIDFromCtx(l.ctx)

	if token == "" {
		return errcode.Fail.Msgr("Token不能为空")
	}

	// 如果是注销账号
	if req.Type == "delete" {
		if err := l.deleteAccount(userID); err != nil {
			return err
		}
	}

	// 从 Redis 中删除 token
	cacheKey := constant.RedisKeyPrefixToken + token
	_, err := l.svcCtx.RedisClient.DelCtx(l.ctx, cacheKey)
	if err != nil {
		l.Logger.Errorf("用户退出登录失败，删除 Redis token 失败 userID : %d, err : %v", userID, err)
		return errcode.InternalServerError.WithError(errors.Wrapf(err, "UserLogout DelCtx userID : %d", userID)).Msgr("退出登录失败")
	}

	return nil
}

// deleteAccount 注销账号逻辑
func (l *UserLogoutLogic) deleteAccount(userID int64) error {
	// 1. 获取当前用户信息
	currentUser, err := l.svcCtx.UserModel.FindOne(l.ctx, cast.ToUint64(userID))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return errcode.Fail.Msgr("用户不存在")
		}
		return errcode.DBError.WithError(errors.Wrapf(err, "UserLogout FindOne")).Msgr("获取用户信息失败")
	}

	currentOpenid := currentUser.Openid

	// 2. 查找所有已注销的同源账号（openid like 'currentOpenid_%'）
	// 使用 SelectBuilder 构造查询
	// SELECT openid FROM users WHERE openid LIKE 'currentOpenid_%'
	likePattern := currentOpenid + "_%"
	queryBuilder := l.svcCtx.UserModel.SelectBuilder("openid").
		Where("openid LIKE ?", likePattern)

	var deletedUsers []*model.Users
	deletedUsers, err = l.svcCtx.UserModel.FindAll(l.ctx, queryBuilder)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		// 如果出错但不是 NotFound，则返回错误
		return errcode.DBError.WithError(errors.Wrapf(err, "UserLogout FindAll")).Msgr("查询历史注销账号失败")
	}

	// 3. 计算新的后缀
	maxSuffix := 0
	for _, u := range deletedUsers {
		// 提取后缀
		parts := strings.Split(u.Openid, "_")
		if len(parts) > 1 {
			suffixStr := parts[len(parts)-1]
			suffix, err := strconv.Atoi(suffixStr)
			if err == nil {
				if suffix > maxSuffix {
					maxSuffix = suffix
				}
			}
		}
	}

	newSuffix := maxSuffix + 1
	newOpenid := fmt.Sprintf("%s_%d", currentOpenid, newSuffix)

	// 4. 更新当前用户的 openid
	updateData := map[string]interface{}{
		"openid":     newOpenid,
		"updated_at": time.Now().Unix(),
	}
	where := squirrel.Eq{"user_id": currentUser.UserId}
	_, err = l.svcCtx.UserModel.UpdateFilter(l.ctx, nil, updateData, where)
	if err != nil {
		return errcode.DBError.WithError(errors.Wrapf(err, "UserLogout UpdateFilter")).Msgr("注销账号失败")
	}

	return nil
}
