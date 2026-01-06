// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/pudongping/momento-api/internal/config"
	"github.com/pudongping/momento-api/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config      config.Config
	RedisClient *redis.Redis
	MysqlClient sqlx.SqlConn

	UserModel model.UsersModel
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

	return &ServiceContext{
		Config:      c,
		RedisClient: redisClient,
		MysqlClient: mysqlConn,

		UserModel: model.NewUsersModel(mysqlConn),
	}
}
