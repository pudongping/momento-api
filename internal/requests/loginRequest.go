package requests

import (
	"github.com/pudongping/momento-api/coreKit/validator"
	"github.com/thedevsaddam/govalidator"
)

func LoginRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"code": []string{"required", "min_cn:2", "max_cn:150"},
	}

	messages := govalidator.MapData{
		"code": []string{
			"required:code为必填项",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}
