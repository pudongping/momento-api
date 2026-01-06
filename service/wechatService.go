package service

import (
	"context"

	"github.com/pkg/errors"

	wechatPkg "github.com/silenceper/wechat/v2"
	wecahtCache "github.com/silenceper/wechat/v2/cache"
	wechatMiniCfg "github.com/silenceper/wechat/v2/miniprogram/config"
)

type WechatService struct {
	ctx context.Context
}

func NewWechatService(ctx context.Context) *WechatService {
	return &WechatService{ctx: ctx}
}

func (svc *WechatService) GetWXMiniOpenID(appID, appSecret, code string) (openID string, err error) {
	miniprogram := wechatPkg.NewWechat().GetMiniProgram(&wechatMiniCfg.Config{
		AppID:     appID,
		AppSecret: appSecret,
		Cache:     wecahtCache.NewMemory(),
	})

	resp, err := miniprogram.GetAuth().Code2SessionContext(svc.ctx, code)
	if err != nil || resp.ErrCode != 0 || resp.OpenID == "" {
		return "", errors.Errorf("微信小程序授权请求失败 err -> %v, code -> %s, resp -> %+v", err, code, resp)
	}

	return resp.OpenID, nil
}
