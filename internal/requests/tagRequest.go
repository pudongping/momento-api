package requests

import (
	"github.com/pudongping/momento-api/coreKit/validator"
	"github.com/thedevsaddam/govalidator"
)

func TagListRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"type": []string{"in:expense,income"},
	}

	messages := govalidator.MapData{
		"type": []string{
			"in:类型参数仅支持 expense 或 income",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}
