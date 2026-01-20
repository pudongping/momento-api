package requests

import (
	"github.com/pudongping/momento-api/coreKit/validator"
	"github.com/thedevsaddam/govalidator"
)

func TransactionListRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"book_id":             []string{"required", "numeric_between:1,"},
		"type":                []string{"in:expense,income"},
		"tag_id":              []string{"numeric_between:1,"},
		"start_date":          []string{"numeric"},
		"end_date":            []string{"numeric"},
		"page":                []string{"numeric_between:1,"},
		"per_page":            []string{"numeric_between:1,"},
		"last_transaction_id": []string{"required", "numeric_between:0,"},
	}

	messages := govalidator.MapData{
		"book_id": []string{
			"required:账本ID为必填项",
			"numeric_between:账本ID必须为数字并且大于等于 1",
		},
		"type": []string{
			"in:交易类型仅支持 expense 或 income",
		},
		"tag_id": []string{
			"numeric_between:标签ID必须为数字并且大于等于 1",
		},
		"start_date": []string{
			"numeric:开始时间必须为时间戳格式",
		},
		"end_date": []string{
			"numeric:结束时间必须为时间戳格式",
		},
		"page": []string{
			"numeric_between:页码必须为数字并且大于等于 1",
		},
		"per_page": []string{
			"numeric_between:每页数量必须为数字并且大于等于 1",
		},
		"last_transaction_id": []string{
			"required:LastTransactionId 为必填项",
			"numeric_between:LastTransactionId 必须为数字并且大于等于 0",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}

func TransactionStatsRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"book_id":    []string{"required", "numeric_between:1,"},
		"type":       []string{"in:expense,income"},
		"tag_id":     []string{"numeric_between:1,"},
		"start_date": []string{"numeric"},
		"end_date":   []string{"numeric"},
	}

	messages := govalidator.MapData{
		"book_id": []string{
			"required:账本ID为必填项",
			"numeric_between:账本ID必须为数字并且大于等于 1",
		},
		"type": []string{
			"in:交易类型仅支持 expense 或 income",
		},
		"tag_id": []string{
			"numeric_between:标签ID必须为数字并且大于等于 1",
		},
		"start_date": []string{
			"numeric:开始时间必须为时间戳格式",
		},
		"end_date": []string{
			"numeric:结束时间必须为时间戳格式",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}

func TransactionAddRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"book_id":           []string{"required", "numeric_between:1,"},
		"type":              []string{"required", "in:expense,income"},
		"amount":            []string{"required", "float", "number_min:0.01"},
		"tag_id":            []string{"required", "numeric_between:1,"},
		"remark":            []string{"max_cn:200"},
		"created_at":        []string{"required", "numeric"},
		"name":              []string{"max_cn:50"},
		"recurring_type":    []string{"in:daily,weekly,monthly,quarterly,yearly"},
		"recurring_hour":    []string{"numeric_between:0,23"},
		"recurring_minute":  []string{"numeric_between:0,59"},
		"recurring_weekday": []string{"numeric_between:0,6"},
		"recurring_month":   []string{"numeric_between:1,12"},
		"recurring_day":     []string{"numeric_between:1,31"},
	}

	messages := govalidator.MapData{
		"book_id": []string{
			"required:账本ID为必填项",
			"numeric_between:账本ID必须为数字并且大于等于 1",
		},
		"type": []string{
			"required:交易类型为必填项",
			"in:交易类型仅支持 expense 或 income",
		},
		"amount": []string{
			"required:金额为必填项",
			"float:金额必须为数字",
			"number_min:金额必须大于等于 0.01",
		},
		"tag_id": []string{
			"required:标签ID为必填项",
			"numeric_between:标签ID必须为数字并且大于等于 1",
		},
		"remark": []string{
			"max_cn:备注不能超过 200 个字符",
		},
		"created_at": []string{
			"required:交易时间为必填项",
			"numeric:交易时间必须为时间戳格式",
		},
		"name": []string{
			"max_cn:规则名称不能超过 50 个字符",
		},
		"recurring_type": []string{
			"in:周期类型无效",
		},
		"recurring_hour": []string{
			"numeric_between:执行小时必须在 0-23 之间",
		},
		"recurring_minute": []string{
			"numeric_between:执行分钟必须在 0-59 之间",
		},
		"recurring_weekday": []string{
			"numeric_between:执行周必须在 0-6 之间",
		},
		"recurring_month": []string{
			"numeric_between:执行月必须在 1-12 之间",
		},
		"recurring_day": []string{
			"numeric_between:执行日必须在 1-31 之间",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}

func TransactionUpdateRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"transaction_id": []string{"required", "numeric_between:1,"},
		"type":           []string{"in:expense,income"},
		"amount":         []string{"float", "number_min:0.01"},
		"tag_id":         []string{"numeric_between:1,"},
		"remark":         []string{"max_cn:200"},
		"created_at":     []string{"numeric"},
	}

	messages := govalidator.MapData{
		"transaction_id": []string{
			"required:交易记录ID为必填项",
			"numeric_between:交易记录ID必须为数字并且大于等于 1",
		},
		"type": []string{
			"in:交易类型仅支持 expense 或 income",
		},
		"amount": []string{
			"float:金额必须为数字",
			"number_min:金额必须大于等于 0.01",
		},
		"tag_id": []string{
			"numeric_between:标签ID必须为数字并且大于等于 1",
		},
		"remark": []string{
			"max_cn:备注不能超过 200 个字符",
		},
		"created_at": []string{
			"numeric:交易时间必须为时间戳格式",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}

func TransactionDeleteRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"transaction_id": []string{"required", "numeric_between:1,"},
	}

	messages := govalidator.MapData{
		"transaction_id": []string{
			"required:交易记录ID为必填项",
			"numeric_between:交易记录ID必须为数字并且大于等于 1",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}
