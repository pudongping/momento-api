package requests

import (
	"github.com/pudongping/momento-api/coreKit/validator"
	"github.com/thedevsaddam/govalidator"
)

func RecurringListRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"book_id": []string{"required", "numeric_between:1,"},
	}

	messages := govalidator.MapData{
		"book_id": []string{
			"required:账本ID为必填项",
			"numeric_between:账本ID必须为数字并且大于等于 1",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}

func RecurringDeleteRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"recurring_id": []string{"required", "numeric_between:1,"},
	}

	messages := govalidator.MapData{
		"recurring_id": []string{
			"required:周期性记账规则ID为必填项",
			"numeric_between:周期性记账规则ID必须为数字并且大于等于 1",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}
