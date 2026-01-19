package requests

import (
	"github.com/pudongping/momento-api/coreKit/validator"
	"github.com/thedevsaddam/govalidator"
)

func UserUpdateRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"nickname": []string{"min_cn:1", "max_cn:8"},
		"avatar":   []string{"url"},
		"phone":    []string{"digits:11"},
	}

	messages := govalidator.MapData{
		"nickname": []string{
			"min_cn:昵称必须大于等于1个字符",
			"max_cn:昵称不能超过8个字符",
		},
		"avatar": []string{
			"url:头像地址必须为 http 或 https 链接",
		},
		"phone": []string{
			"digits:手机号格式错误",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}
