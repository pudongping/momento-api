// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	AppService struct {
		StaticFSRelativePath string
	}

	Mysql struct {
		DataSource string
	}

	Redis redis.RedisConf

	CacheRedis cache.CacheConf

	JWTAuth struct {
		AccessSecret string
		AccessExpire int64
	}

	WXMiniProgram struct {
		AppID     string `json:"AppID"`     // 微信小程序的 app_id
		AppSecret string `json:"AppSecret"` // 微信小程序的 app_secret
	}
}
