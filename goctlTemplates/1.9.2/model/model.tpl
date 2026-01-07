package {{.pkg}}
{{if .withCache}}
import (
	"context"
	"strings"

	"your-project-module-name/coreKit/paginator"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/Masterminds/squirrel"
)
{{else}}

import (
	"context"
	"strings"

	"your-project-module-name/coreKit/paginator"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/Masterminds/squirrel"
)

{{end}}
var _ {{.upperStartCamelObject}}Model = (*custom{{.upperStartCamelObject}}Model)(nil)

type (
	// {{.upperStartCamelObject}}Model is an interface to be customized, add more methods here,
	// and implement the added methods in custom{{.upperStartCamelObject}}Model.
	{{.upperStartCamelObject}}Model interface {
		{{.lowerStartCamelObject}}Model
		{{if not .withCache}}withSession(session sqlx.Session) {{.upperStartCamelObject}}Model{{end}}

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
		FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*{{.upperStartCamelObject}}, error)
		FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*{{.upperStartCamelObject}}, error)
		FindListByPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, perPage int64, orderBy string) ([]*{{.upperStartCamelObject}}, *paginator.Pagination, error)

	}

	custom{{.upperStartCamelObject}}Model struct {
		*default{{.upperStartCamelObject}}Model
	}
)

// New{{.upperStartCamelObject}}Model returns a model for the database table.
func New{{.upperStartCamelObject}}Model(conn sqlx.SqlConn{{if .withCache}}, c cache.CacheConf, opts ...cache.Option{{end}}) {{.upperStartCamelObject}}Model {
	return &custom{{.upperStartCamelObject}}Model{
		default{{.upperStartCamelObject}}Model: new{{.upperStartCamelObject}}Model(conn{{if .withCache}}, c, opts...{{end}}),
	}
}

{{if not .withCache}}
func (m *custom{{.upperStartCamelObject}}Model) withSession(session sqlx.Session) {{.upperStartCamelObject}}Model {
    return New{{.upperStartCamelObject}}Model(sqlx.NewSqlConnFromSession(session))
}
{{end}}


// 自定义模版中扩展的方法 start ----->

func (m *default{{.upperStartCamelObject}}Model) GetTableName() string {
    return m.table
}

func (m *default{{.upperStartCamelObject}}Model) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
    {{if .withCache}}return m.ExecNoCacheCtx(ctx, query, args...){{else}}
    return m.conn.ExecCtx(ctx, query, args...)
    {{end}}
}

func (m *default{{.upperStartCamelObject}}Model) DeleteFilter(ctx context.Context, where squirrel.Sqlizer) (sql.Result, error) {
    deleteBuilder := squirrel.Delete(m.table).Where(where)
    query, values, err := deleteBuilder.ToSql()
    if err != nil {
        return nil, err
    }

    return m.ExecContext(ctx, query, values...)
}

func (m *default{{.upperStartCamelObject}}Model) UpdateFilter(ctx context.Context, updateData map[string]interface{}, where squirrel.Sqlizer) (sql.Result, error) {
    updateBuilder := squirrel.Update(m.table).SetMap(updateData).Where(where)
    query, values, err := updateBuilder.ToSql()
    if err != nil {
        return nil, err
    }

    return m.ExecContext(ctx, query, values...)
}

func (m *default{{.upperStartCamelObject}}Model) Transaction(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	{{if .withCache}}
	return m.TransactCtx(ctx,func(ctx context.Context,session sqlx.Session) error {
		return  fn(ctx,session)
	})
	{{else}}
	return m.conn.TransactCtx(ctx,func(ctx context.Context,session sqlx.Session) error {
		return  fn(ctx,session)
	})
	{{end}}
}

func (m *default{{.upperStartCamelObject}}Model) SelectBuilder(fields ...string) squirrel.SelectBuilder {
	f := {{.lowerStartCamelObject}}Rows
	if len(fields) > 0 {
		f = strings.Join(fields, ",")
	}
	return squirrel.Select(f).From(m.table)
}

func (m *default{{.upperStartCamelObject}}Model) CountBuilder(field ...string) squirrel.SelectBuilder {
	f := "*"
	if len(field) > 0 && field[0] != "" {
		f = field[0]
	}
	return squirrel.Select("COUNT(" + f + ")").From(m.table)
}

func (m *default{{.upperStartCamelObject}}Model) SumBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("IFNULL(SUM(" + field + "),0)").From(m.table)
}

func (m *default{{.upperStartCamelObject}}Model) FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error) {
	query, values, err := countBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var resp int64
	{{if .withCache}}err = m.QueryRowNoCacheCtx(ctx,&resp, query, values...){{else}}
	err = m.conn.QueryRowCtx(ctx,&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *default{{.upperStartCamelObject}}Model) FindSum(ctx context.Context, sumBuilder squirrel.SelectBuilder) (float64, error) {
	query, values, err := sumBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var resp float64
	{{if .withCache}}err = m.QueryRowNoCacheCtx(ctx,&resp, query, values...){{else}}
	err = m.conn.QueryRowCtx(ctx,&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *default{{.upperStartCamelObject}}Model) FindAny(ctx context.Context, rowBuilder squirrel.SelectBuilder, bindings interface{}) error {
	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return err
	}

	{{if .withCache}}return m.QueryRowsNoCacheCtx(ctx, bindings, query, values...){{else}}
	return m.conn.QueryRowsCtx(ctx, bindings, query, values...)
	{{end}}
}

func (m *default{{.upperStartCamelObject}}Model) FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*{{.upperStartCamelObject}}, error) {
	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp {{.upperStartCamelObject}}
	{{if .withCache}}err = m.QueryRowNoCacheCtx(ctx,&resp, query, values...){{else}}
	err = m.conn.QueryRowCtx(ctx,&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}

}

func (m *default{{.upperStartCamelObject}}Model) FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*{{.upperStartCamelObject}}, error) {
	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*{{.upperStartCamelObject}}
	{{if .withCache}}err = m.QueryRowsNoCacheCtx(ctx,&resp, query, values...){{else}}
	err = m.conn.QueryRowsCtx(ctx,&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *default{{.upperStartCamelObject}}Model) FindListByPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, perPage int64, orderBy string) ([]*{{.upperStartCamelObject}}, *paginator.Pagination, error) {
	// 构建查询数据总条数的 sql 语句
	countQuery, countValues, err := paginator.CountDataSqlBuilder(rowBuilder).ToSql()
	if err != nil {
		return nil, nil, err
	}

	var total int64
	{{if .withCache}}err = m.QueryRowNoCacheCtx(ctx,&total, countQuery, countValues...){{else}}
	err = m.conn.QueryRowCtx(ctx,&total, countQuery, countValues...)
	{{end}}
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

	var resp []*{{.upperStartCamelObject}}
	{{if .withCache}}err = m.QueryRowsNoCacheCtx(ctx,&resp, query, values...){{else}}
	err = m.conn.QueryRowsCtx(ctx,&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return resp, pagination, nil
	default:
		return nil, nil, err
	}
}

// 自定义模版中扩展的方法 end ---->