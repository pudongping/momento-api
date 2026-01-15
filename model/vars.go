package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound

const (
	// users 表相关
	UserIsDisableYes = 1 // 启用
	UserIsDisableNo  = 2 // 禁用
)

const (
	// tags 表相关
	TagsIsSystemYes = 1 // 系统标签
	TagsIsSystemNo  = 2 // 用户自定义标签

	// type
	TagsTypeExpense = "expense" // 支出
	TagsTypeIncome  = "income"  // 收入
)

const (
	// festivals 表相关

	// is_show_home 是否显示在首页
	FestivalIsShowHomeYes = 1 // 显示
	FestivalIsShowHomeNo  = 2 // 不显示
)
