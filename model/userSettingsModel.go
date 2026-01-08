package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserSettingsModel = (*customUserSettingsModel)(nil)

type (
	// UserSettingsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserSettingsModel.
	UserSettingsModel interface {
		userSettingsModel
		withSession(session sqlx.Session) UserSettingsModel
	}

	customUserSettingsModel struct {
		*defaultUserSettingsModel
	}
)

// NewUserSettingsModel returns a model for the database table.
func NewUserSettingsModel(conn sqlx.SqlConn) UserSettingsModel {
	return &customUserSettingsModel{
		defaultUserSettingsModel: newUserSettingsModel(conn),
	}
}

func (m *customUserSettingsModel) withSession(session sqlx.Session) UserSettingsModel {
	return NewUserSettingsModel(sqlx.NewSqlConnFromSession(session))
}
