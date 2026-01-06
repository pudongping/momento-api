package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound

const (
	UserIsDisableYes = 1 // 启用
	UserIsDisableNo  = 2 // 禁用
)
