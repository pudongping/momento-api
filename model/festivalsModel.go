package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ FestivalsModel = (*customFestivalsModel)(nil)

type (
	// FestivalsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFestivalsModel.
	FestivalsModel interface {
		festivalsModel
		withSession(session sqlx.Session) FestivalsModel
	}

	customFestivalsModel struct {
		*defaultFestivalsModel
	}
)

// NewFestivalsModel returns a model for the database table.
func NewFestivalsModel(conn sqlx.SqlConn) FestivalsModel {
	return &customFestivalsModel{
		defaultFestivalsModel: newFestivalsModel(conn),
	}
}

func (m *customFestivalsModel) withSession(session sqlx.Session) FestivalsModel {
	return NewFestivalsModel(sqlx.NewSqlConnFromSession(session))
}
