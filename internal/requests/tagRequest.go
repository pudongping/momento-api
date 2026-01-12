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

func TagAddRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"name":  []string{"required", "min_cn:1", "max_cn:6"},
		"color": []string{"css_color"},
		"icon":  []string{"max:10"},
		"type":  []string{"required", "in:expense,income"},
	}

	messages := govalidator.MapData{
		"name": []string{
			"required:标签名称为必填项",
			"min_cn:标签名称长度必须至少 1 个字符",
			"max_cn:标签名称长度不能超过 6 个字符",
		},
		"color": []string{
			"css_color:标签颜色必须是有效的颜色值(如 #E91E63)",
		},
		"icon": []string{
			"max:标签图标长度不能超过 10 个字符",
		},
		"type": []string{
			"required:标签类型为必填项",
			"in:标签类型仅支持 expense 或 income",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}
