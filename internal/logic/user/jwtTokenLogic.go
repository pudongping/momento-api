package user

import (
	"context"
	"time"

	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/jwtToken"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type JwtTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJwtTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JwtTokenLogic {
	return &JwtTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GenerateToken 生成 JWT token
// args:
//
//	userID int64 用户 ID
//
// returns:
//
//	string 生成的 JWT token
//	int64 访问令牌过期时间（秒）
//	int64 刷新时间（秒）
//	error 错误信息
func (l *JwtTokenLogic) GenerateToken(userID int64) (string, int64, int64, error) {
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JWTAuth.AccessExpire
	accessSecret := l.svcCtx.Config.JWTAuth.AccessSecret
	refreshTime := now + accessExpire/2

	token, err := jwtToken.GenJwtToken(accessSecret, ctxData.CtxKeyJwtUserId, now, accessExpire, userID)
	return token, accessExpire, refreshTime, err
}
