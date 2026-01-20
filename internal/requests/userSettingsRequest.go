package requests

import (
	"github.com/pudongping/momento-api/coreKit/validator"
	"github.com/thedevsaddam/govalidator"
)

func UpdateUserSettingsRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"background_url": []string{"url"},
		"budget":         []string{"number_min:0", "number_max:1000000"},
	}

	messages := govalidator.MapData{
		"background_url": []string{
			"url:背景图片地址必须为 http 或 https 链接",
		},
		"budget": []string{
			"number_min:预算最小为0",
			"number_max:预算最大为1000000",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}
