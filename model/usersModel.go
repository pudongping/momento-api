package model

import (
	"context"
	"database/sql"
	"strings"

	"github.com/pudongping/momento-api/coreKit/paginator"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsersModel = (*customUsersModel)(nil)

type (
	// UsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersModel.
	UsersModel interface {
		usersModel
		withSession(session sqlx.Session) UsersModel

		GetTableName() string
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		DeleteFilter(ctx context.Context, where squirrel.Sqlizer) (sql.Result, error)
		UpdateFilter(ctx context.Context, updateData map[string]interface{}, where squirrel.Sqlizer) (sql.Result, error)

		Transaction(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		SelectBuilder(fields ...string) squirrel.SelectBuilder
		CountBuilder(field ...string) squirrel.SelectBuilder
		SumBuilder(field string) squirrel.SelectBuilder
		FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error)
		FindSum(ctx context.Context, sumBuilder squirrel.SelectBuilder) (float64, error)
		FindAny(ctx context.Context, rowBuilder squirrel.SelectBuilder, bindings interface{}) error
		FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*Users, error)
		FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*Users, error)
		FindListByPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, perPage int64, orderBy string) ([]*Users, *paginator.Pagination, error)
	}

	customUsersModel struct {
		*defaultUsersModel
	}
)

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn),
	}
}

func (m *customUsersModel) withSession(session sqlx.Session) UsersModel {
	return NewUsersModel(sqlx.NewSqlConnFromSession(session))
}

// 自定义模版中扩展的方法 start ----->

func (m *defaultUsersModel) GetTableName() string {
	return m.table
}

func (m *defaultUsersModel) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {

	return m.conn.ExecCtx(ctx, query, args...)

}

func (m *defaultUsersModel) DeleteFilter(ctx context.Context, where squirrel.Sqlizer) (sql.Result, error) {
	deleteBuilder := squirrel.Delete(m.table).Where(where)
	query, values, err := deleteBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	return m.ExecContext(ctx, query, values...)
}

func (m *defaultUsersModel) UpdateFilter(ctx context.Context, updateData map[string]interface{}, where squirrel.Sqlizer) (sql.Result, error) {
	updateBuilder := squirrel.Update(m.table).SetMap(updateData).Where(where)
	query, values, err := updateBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	return m.ExecContext(ctx, query, values...)
}

func (m *defaultUsersModel) Transaction(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {

	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})

}

func (m *defaultUsersModel) SelectBuilder(fields ...string) squirrel.SelectBuilder {
	f := usersRows
	if len(fields) > 0 {
		f = strings.Join(fields, ",")
	}
	return squirrel.Select(f).From(m.table)
}

func (m *defaultUsersModel) CountBuilder(field ...string) squirrel.SelectBuilder {
	f := "*"
	if len(field) > 0 && field[0] != "" {
		f = field[0]
	}
	return squirrel.Select("COUNT(" + f + ")").From(m.table)
}

func (m *defaultUsersModel) SumBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("IFNULL(SUM(" + field + "),0)").From(m.table)
}

func (m *defaultUsersModel) FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error) {
	query, values, err := countBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var resp int64

	err = m.conn.QueryRowCtx(ctx, &resp, query, values...)

	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultUsersModel) FindSum(ctx context.Context, sumBuilder squirrel.SelectBuilder) (float64, error) {
	query, values, err := sumBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var resp float64

	err = m.conn.QueryRowCtx(ctx, &resp, query, values...)

	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultUsersModel) FindAny(ctx context.Context, rowBuilder squirrel.SelectBuilder, bindings interface{}) error {
	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return err
	}

	return m.conn.QueryRowsCtx(ctx, bindings, query, values...)

}

func (m *defaultUsersModel) FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*Users, error) {
	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp Users

	err = m.conn.QueryRowCtx(ctx, &resp, query, values...)

	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}

}

func (m *defaultUsersModel) FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*Users, error) {
	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Users

	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultUsersModel) FindListByPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, perPage int64, orderBy string) ([]*Users, *paginator.Pagination, error) {
	// 构建查询数据总条数的 sql 语句
	countQuery, countValues, err := paginator.CountDataSqlBuilder(rowBuilder).ToSql()
	if err != nil {
		return nil, nil, err
	}

	var total int64

	err = m.conn.QueryRowCtx(ctx, &total, countQuery, countValues...)

	if err != nil {
		return nil, nil, err
	}

	currentPage, limit, offset := paginator.PrepareOffsetLimit(page, perPage)
	pagination := paginator.NewPagination(total, currentPage, limit)

	builder := rowBuilder.Offset(uint64(offset)).Limit(uint64(limit))
	builder = paginator.WithOrderBy(orderBy, builder)

	query, values, err := builder.ToSql()
	if err != nil {
		return nil, nil, err
	}

	var resp []*Users

	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...)

	switch err {
	case nil:
		return resp, pagination, nil
	default:
		return nil, nil, err
	}
}

// 自定义模版中扩展的方法 end ---->
