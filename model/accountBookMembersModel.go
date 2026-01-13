package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AccountBookMembersModel = (*customAccountBookMembersModel)(nil)

type (
	// AccountBookMembersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAccountBookMembersModel.
	AccountBookMembersModel interface {
		accountBookMembersModel
		withSession(session sqlx.Session) AccountBookMembersModel
	}

	customAccountBookMembersModel struct {
		*defaultAccountBookMembersModel
	}
)

// NewAccountBookMembersModel returns a model for the database table.
func NewAccountBookMembersModel(conn sqlx.SqlConn) AccountBookMembersModel {
	return &customAccountBookMembersModel{
		defaultAccountBookMembersModel: newAccountBookMembersModel(conn),
	}
}

func (m *customAccountBookMembersModel) withSession(session sqlx.Session) AccountBookMembersModel {
	return NewAccountBookMembersModel(sqlx.NewSqlConnFromSession(session))
}
