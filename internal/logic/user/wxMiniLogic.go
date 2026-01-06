package user

import (
	"context"

	"github.com/pudongping/momento-api/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"

	wechatPkg "github.com/silenceper/wechat/v2"
)

type WXMiniLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWXMiniLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WXMiniLogic {
	return &WXMiniLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WXMiniLogic) WXMiniAuth(code string) (openId string, err error) {
	wechatPkg.NewWechat()
	return "", err
}
