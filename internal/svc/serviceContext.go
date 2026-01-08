// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/pudongping/momento-api/internal/config"
	"github.com/pudongping/momento-api/internal/middleware"
	"github.com/pudongping/momento-api/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config      config.Config
	RedisClient *redis.Redis
	MysqlClient sqlx.SqlConn

	AuthCheckMiddleware rest.Middleware

	UserModel        model.UsersModel
	UserSettingModel model.UserSettingsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	redisClient := redis.MustNewRedis(redis.RedisConf{
		Host:        c.Redis.Host,
		Type:        c.Redis.Type,
		User:        c.Redis.User,
		Pass:        c.Redis.Pass,
		Tls:         false,
		NonBlock:    false,
		PingTimeout: 0,
	})

	userModel := model.NewUsersModel(mysqlConn)
	userSettingModel := model.NewUserSettingsModel(mysqlConn)

	return &ServiceContext{
		Config:      c,
		RedisClient: redisClient,
		MysqlClient: mysqlConn,

		AuthCheckMiddleware: middleware.NewAuthCheckMiddleware(userModel, redisClient).Handle,

		UserModel:        userModel,
		UserSettingModel: userSettingModel,
	}
}
