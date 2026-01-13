package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TransactionsModel = (*customTransactionsModel)(nil)

type (
	// TransactionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTransactionsModel.
	TransactionsModel interface {
		transactionsModel
		withSession(session sqlx.Session) TransactionsModel
	}

	customTransactionsModel struct {
		*defaultTransactionsModel
	}
)

// NewTransactionsModel returns a model for the database table.
func NewTransactionsModel(conn sqlx.SqlConn) TransactionsModel {
	return &customTransactionsModel{
		defaultTransactionsModel: newTransactionsModel(conn),
	}
}

func (m *customTransactionsModel) withSession(session sqlx.Session) TransactionsModel {
	return NewTransactionsModel(sqlx.NewSqlConnFromSession(session))
}
