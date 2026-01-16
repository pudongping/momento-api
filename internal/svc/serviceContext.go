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

	UserModel                  model.UsersModel
	UserSettingModel           model.UserSettingsModel
	TagsModel                  model.TagsModel
	FestivalsModel             model.FestivalsModel
	TransactionsModel          model.TransactionsModel
	RecurringTransactionsModel model.RecurringTransactionsModel
	AccountBooksModel          model.AccountBooksModel
	AccountBookMembersModel    model.AccountBookMembersModel
	AccountBookInvitationsModel model.AccountBookInvitationsModel
	UploadFilesModel           model.UploadFilesModel
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
	tagsModel := model.NewTagsModel(mysqlConn)
	festivalsModel := model.NewFestivalsModel(mysqlConn)
	transactionsModel := model.NewTransactionsModel(mysqlConn)
	recurringTransactionsModel := model.NewRecurringTransactionsModel(mysqlConn)
	accountBooksModel := model.NewAccountBooksModel(mysqlConn)
	accountBookMembersModel := model.NewAccountBookMembersModel(mysqlConn)
	accountBookInvitationsModel := model.NewAccountBookInvitationsModel(mysqlConn)
	uploadFilesModel := model.NewUploadFilesModel(mysqlConn)

	return &ServiceContext{
		Config:      c,
		RedisClient: redisClient,
		MysqlClient: mysqlConn,

		AuthCheckMiddleware: middleware.NewAuthCheckMiddleware(userModel, redisClient).Handle,

		UserModel:                  userModel,
		UserSettingModel:           userSettingModel,
		TagsModel:                  tagsModel,
		FestivalsModel:             festivalsModel,
		TransactionsModel:          transactionsModel,
		RecurringTransactionsModel: recurringTransactionsModel,
		AccountBooksModel:          accountBooksModel,
		AccountBookMembersModel:    accountBookMembersModel,
		AccountBookInvitationsModel: accountBookInvitationsModel,
		UploadFilesModel:           uploadFilesModel,
	}
}
