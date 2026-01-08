Insert(ctx context.Context, data *{{.upperStartCamelObject}}) (sql.Result,error)

// 自定义模版中扩展的方法 start ----->
GetTableName() string
Transaction(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
ExecContext(ctx context.Context, session sqlx.Session, query string, args ...interface{}) (sql.Result, error)
DeleteFilter(ctx context.Context, session sqlx.Session, where squirrel.Sqlizer) (sql.Result, error)
UpdateFilter(ctx context.Context, session sqlx.Session, updateData map[string]interface{}, where squirrel.Sqlizer) (sql.Result, error)
SelectBuilder(fields ...string) squirrel.SelectBuilder
CountBuilder(field ...string) squirrel.SelectBuilder
SumBuilder(field string) squirrel.SelectBuilder
FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error)
FindSum(ctx context.Context, sumBuilder squirrel.SelectBuilder) (float64, error)
FindAny(ctx context.Context, rowBuilder squirrel.SelectBuilder, bindings interface{}) error
FindOneByQuery(ctx context.Context, rowBuilder squirrel.SelectBuilder) (*{{.upperStartCamelObject}}, error)
FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*{{.upperStartCamelObject}}, error)
FindListByPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, perPage int64, orderBy string) ([]*{{.upperStartCamelObject}}, *paginator.Pagination, error)
// 自定义模版中扩展的方法 end ---->