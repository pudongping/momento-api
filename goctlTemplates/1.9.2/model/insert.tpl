func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, data *{{.upperStartCamelObject}}) (sql.Result,error) {
	{{if .withCache}}{{.keys}}
    ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
    ret,err:=m.conn.ExecCtx(ctx, query, {{.expressionValues}}){{end}}
	return ret,err
}

// 自定义模版中扩展的方法 start ----->

func (m *default{{.upperStartCamelObject}}Model) GetTableName() string {
    return m.tableName()
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

func (m *default{{.upperStartCamelObject}}Model) ExecContext(ctx context.Context, session sqlx.Session, query string, args ...interface{}) (sql.Result, error) {
    if session != nil {
        return session.ExecCtx(ctx, query, args...)
    }
    return m.conn.ExecCtx(ctx, query, args...)
}

func (m *default{{.upperStartCamelObject}}Model) DeleteFilter(ctx context.Context, session sqlx.Session, where squirrel.Sqlizer) (sql.Result, error) {
    deleteBuilder := squirrel.Delete(m.table).Where(where)
    query, values, err := deleteBuilder.ToSql()
    if err != nil {
        return nil, err
    }

    return m.ExecContext(ctx, session, query, values...)
}

func (m *default{{.upperStartCamelObject}}Model) UpdateFilter(ctx context.Context, session sqlx.Session, updateData map[string]interface{}, where squirrel.Sqlizer) (sql.Result, error) {
    updateBuilder := squirrel.Update(m.table).SetMap(updateData).Where(where)
    query, values, err := updateBuilder.ToSql()
    if err != nil {
        return nil, err
    }

    return m.ExecContext(ctx, session, query, values...)
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
