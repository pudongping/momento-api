package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AccountBookInvitationsModel = (*customAccountBookInvitationsModel)(nil)

type (
	// AccountBookInvitationsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAccountBookInvitationsModel.
	AccountBookInvitationsModel interface {
		accountBookInvitationsModel
		withSession(session sqlx.Session) AccountBookInvitationsModel
	}

	customAccountBookInvitationsModel struct {
		*defaultAccountBookInvitationsModel
	}
)

// NewAccountBookInvitationsModel returns a model for the database table.
func NewAccountBookInvitationsModel(conn sqlx.SqlConn) AccountBookInvitationsModel {
	return &customAccountBookInvitationsModel{
		defaultAccountBookInvitationsModel: newAccountBookInvitationsModel(conn),
	}
}

func (m *customAccountBookInvitationsModel) withSession(session sqlx.Session) AccountBookInvitationsModel {
	return NewAccountBookInvitationsModel(sqlx.NewSqlConnFromSession(session))
}
