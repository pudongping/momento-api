package requests

import (
	"github.com/pudongping/momento-api/coreKit/validator"
	"github.com/thedevsaddam/govalidator"
)

func AccountBookAddRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"name": []string{"required", "min_cn:1", "max_cn:8"},
	}

	messages := govalidator.MapData{
		"name": []string{
			"required:账本名称为必填项",
			"min_cn:账本名称长度必须至少 1 个字符",
			"max_cn:账本名称长度不能超过 8 个字符",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}

func AccountBookDeleteRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"book_id": []string{"required", "numeric_between:1,"},
	}

	messages := govalidator.MapData{
		"book_id": []string{
			"required:账本ID为必填项",
			"numeric_between:账本ID必须是数字并且大于等于 1",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}
