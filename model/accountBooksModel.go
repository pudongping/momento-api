package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AccountBooksModel = (*customAccountBooksModel)(nil)

type (
	// AccountBooksModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAccountBooksModel.
	AccountBooksModel interface {
		accountBooksModel
		withSession(session sqlx.Session) AccountBooksModel
	}

	customAccountBooksModel struct {
		*defaultAccountBooksModel
	}
)

// NewAccountBooksModel returns a model for the database table.
func NewAccountBooksModel(conn sqlx.SqlConn) AccountBooksModel {
	return &customAccountBooksModel{
		defaultAccountBooksModel: newAccountBooksModel(conn),
	}
}

func (m *customAccountBooksModel) withSession(session sqlx.Session) AccountBooksModel {
	return NewAccountBooksModel(sqlx.NewSqlConnFromSession(session))
}
