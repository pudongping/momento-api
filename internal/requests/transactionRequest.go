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
