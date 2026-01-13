package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RecurringTransactionsModel = (*customRecurringTransactionsModel)(nil)

type (
	// RecurringTransactionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRecurringTransactionsModel.
	RecurringTransactionsModel interface {
		recurringTransactionsModel
		withSession(session sqlx.Session) RecurringTransactionsModel
	}

	customRecurringTransactionsModel struct {
		*defaultRecurringTransactionsModel
	}
)

// NewRecurringTransactionsModel returns a model for the database table.
func NewRecurringTransactionsModel(conn sqlx.SqlConn) RecurringTransactionsModel {
	return &customRecurringTransactionsModel{
		defaultRecurringTransactionsModel: newRecurringTransactionsModel(conn),
	}
}

func (m *customRecurringTransactionsModel) withSession(session sqlx.Session) RecurringTransactionsModel {
	return NewRecurringTransactionsModel(sqlx.NewSqlConnFromSession(session))
}
