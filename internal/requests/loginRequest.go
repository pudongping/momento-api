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
			"min_cn:code长度必须大于等于2个字符",
			"max_cn:code长度不能超过150个字符",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}
