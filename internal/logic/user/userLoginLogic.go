// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/coreKit/helpers/flake"
	"github.com/pudongping/momento-api/internal/constant"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/pudongping/momento-api/service"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stringx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 小程序授权登录
func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.LoginReq, clientIP string) (*types.LoginResp, error) {
	var resp types.LoginResp
	var err error

	openID, err := l.FetchOpenID(req.Code)
	if err != nil {
		return nil, errors.Wrapf(err, "UserLogin FetchOpenID code : %s", req.Code)
	}

	// 先检查当前用户是否存在，存在则登录，否则注册新用户
	user, err := l.svcCtx.UserModel.FindOneByOpenid(l.ctx, openID)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errcode.DBError.WithError(errors.Wrap(err, "UserLogin FindOneByOpenid")).Msgr("查询用户失败")
	}
	if user == nil {
		// 用户不存在，注册新用户
		user, err = l.CreateUser(openID, clientIP)
		if err != nil {
			return nil, err
		}
	} else {
		// 用户存在，更新登录信息
		err = l.UpdateUser(cast.ToUint64(user.UserId), clientIP)
		if err != nil {
			return nil, err
		}
	}

	jwtTokenLogic := NewJwtTokenLogic(l.ctx, l.svcCtx)
	token, accessExpire, _, err := jwtTokenLogic.GenerateToken(cast.ToInt64(user.UserId))
	if err != nil {
		return nil, errcode.InternalServerError.WithError(errors.Wrapf(err, "UserLogin GenerateToken userID : %d", user.UserId)).Msgr("生成登录令牌失败")
	}

	// 将 token 存入 redis 中，并设置过期时间
	cacheKey := constant.RedisKeyPrefixToken + token
	_, err = l.svcCtx.RedisClient.SetnxExCtx(l.ctx, cacheKey, cast.ToString(user.UserId), int(accessExpire))
	if err != nil {
		return nil, errcode.InternalServerError.WithError(errors.Wrapf(err, "UserLogin Redis SetnxExCtx UID : %d", user.UserId)).Msgr("存储登录令牌失败")
	}

	if err = copier.Copy(&resp, user); err != nil {
		return nil, errors.Wrapf(err, "拷贝用户信息失败 userID : %d", user.UserId)
	}

	resp.UserID = cast.ToString(user.UserId)
	resp.Token = token

	return &resp, nil
}

func (l *UserLoginLogic) FetchOpenID(code string) (openID string, err error) {
	if "dev" == l.svcCtx.Config.Mode {
		return code, nil
	}

	// 调用微信接口获取 openID
	wechatService := service.NewWechatService(l.ctx)
	appID := l.svcCtx.Config.WXMiniProgram.AppID
	appSecret := l.svcCtx.Config.WXMiniProgram.AppSecret
	openID, err = wechatService.GetWXMiniOpenID(appID, appSecret, code)
	if err != nil {
		return "", errors.Wrapf(err, "FetchOpenID code : %s", code)
	}

	return openID, nil
}

func (l *UserLoginLogic) CreateUser(openID, clientIP string) (*model.Users, error) {
	now := time.Now().Unix()
	userID := flake.GenUniqueID()
	newUserNickname, err := stringx.Substr(cast.ToString(userID), 0, 6)
	if err != nil {
		return nil, errcode.Fail.WithError(errors.Wrapf(err, "生成新用户昵称失败 userID : %d", userID))
	}
	newUserNickname = "新用户_" + newUserNickname

	user := new(model.Users)
	user.UserId = cast.ToUint64(userID)
	user.Openid = openID
	user.Nickname = newUserNickname
	user.IsDisable = model.UserIsDisableNo
	user.RegisterIp = clientIP
	user.LoginIp = clientIP
	user.LastLoginTime = cast.ToUint64(now)
	user.CreatedAt = cast.ToUint64(now)
	user.UpdatedAt = cast.ToUint64(now)

	insertResult, err := l.svcCtx.UserModel.Insert(l.ctx, user)
	if err != nil {
		return nil, errcode.DBError.WithError(errors.Wrapf(err, "用户注册失败 userID : %d, openID : %s", userID, openID))
	}
	lastInsertID, err := insertResult.LastInsertId()
	if err != nil {
		return nil, errcode.DBError.WithError(errors.Wrapf(err, "获取用户注册最后插入ID失败 userID : %d, openID : %s lastInsertID : %d", userID, openID, lastInsertID))
	}

	l.Logger.Infof("新用户注册成功 userID : %d, openID : %s", userID, openID)

	return user, nil
}

func (l *UserLoginLogic) UpdateUser(userID uint64, clientIP string) error {
	now := time.Now().Unix()
	updateMap := map[string]interface{}{
		"last_login_time": cast.ToUint64(now),
		"login_ip":        clientIP,
		"updated_at":      cast.ToUint64(now),
	}

	where := squirrel.Eq{"user_id": userID}

	_, err := l.svcCtx.UserModel.UpdateFilter(l.ctx, nil, updateMap, where)
	if err != nil {
		l.Logger.Errorf("更新用户登录信息失败 userID : %d, err : %v", userID, err)
		return errors.Wrapf(err, "UpdateUser UpdateFilter userID : %d", userID)
	}

	return nil
}
