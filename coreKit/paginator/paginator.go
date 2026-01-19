package paginator

import (
	"fmt"
	"math"
	"strings"

	"github.com/Masterminds/squirrel"
)

// DefaultPerPage 默认每页显示数
const DefaultPerPage = 20

type Pagination struct {
	Total    int64 `json:"total"`    // 数据总条数
	Page     int64 `json:"page"`     // 当前请求页码
	PerPage  int64 `json:"perPage"`  // 每页显示数
	LastPage int64 `json:"lastPage"` // 最大页码（最后一页）
}

func NewPagination(total, page, perPage int64) *Pagination {
	return &Pagination{
		Total:    total,
		Page:     page,
		PerPage:  perPage,
		LastPage: GetTotalPage(total, perPage),
	}
}

// WithOrderBy 在 sql 语句中添加排序字段
func WithOrderBy(orderBy string, rowBuilder squirrel.SelectBuilder) squirrel.SelectBuilder {
	if "" == orderBy {
		return rowBuilder
	}

	sortColumn, sortRules := PrepareOrderBy(orderBy)
	if len(sortColumn) > 0 && len(sortRules) > 0 {
		for i, v := range sortColumn {
			rowBuilder = rowBuilder.OrderBy(fmt.Sprintf("%s %s", v, sortRules[i]))
		}
	}

	return rowBuilder
}

// PrepareOrderBy 处理排序字段
func PrepareOrderBy(orderBy string) (sortColumn []string, sortRules []string) {
	if "" == orderBy {
		return
	}

	// eg： [id,desc name,asc]
	orderBySlice := strings.Split(orderBy, "|")
	for _, orderByItem := range orderBySlice {

		rank := strings.Split(orderByItem, ",")
		if len(rank) != 2 || rank[0] == "" || rank[1] == "" {
			continue
		}

		sortColumn = append(sortColumn, rank[0]) // eg：[id name]
		sortRules = append(sortRules, rank[1])   // eg：[desc asc]
	}

	return
}

// PrepareOffsetLimit 计算偏移量
func PrepareOffsetLimit(page, perPage int64) (int64, int64, int64) {
	if page < 1 {
		page = 1
	}
	if perPage <= 0 {
		perPage = DefaultPerPage
	}

	offset := (page - 1) * perPage

	return page, perPage, offset
}

// GetTotalPage 得到最大页码
func GetTotalPage(totalCount, perPage int64) int64 {
	if 0 == totalCount {
		return 0
	}

	return int64(math.Ceil(float64(totalCount) / float64(perPage)))
}

// CountDataSqlBuilder 构建查询数据总条数的 sql 语句
func CountDataSqlBuilder(rowBuilder squirrel.SelectBuilder) squirrel.SelectBuilder {
	// 获取数据总条数的时候不应该出现 order by、limit、offset
	subQuery := rowBuilder.RemoveOffset().RemoveLimit()
	return squirrel.Select("count(*) as aggregate").FromSelect(subQuery, "aggregate_table")
}
