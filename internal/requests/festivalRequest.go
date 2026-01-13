package requests

import (
	"github.com/pudongping/momento-api/coreKit/validator"
	"github.com/thedevsaddam/govalidator"
)

func FestivalAddRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"festival_name": []string{"required", "min_cn:1", "max_cn:10"},
		"festival_date": []string{"required", "numeric_between:10000000,99999999"},
		"is_show_home":  []string{"in:1,2"},
	}

	messages := govalidator.MapData{
		"festival_name": []string{
			"required:节日名称为必填项",
			"min_cn:节日名称长度必须至少 1 个字符",
			"max_cn:节日名称长度不能超过 10 个字符",
		},
		"festival_date": []string{
			"required:节日日期为必填项",
			"numeric_between:节日日期格式应为 YYYYMMDD",
		},
		"is_show_home": []string{
			"in:是否显示仅支持 1 或 2",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}

func FestivalUpdateRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"festival_id": []string{"required", "numeric_between:1,"},
		"festival_name": []string{"min_cn:1", "max_cn:10"},
		"festival_date": []string{"numeric_between:10000000,99999999"},
		"is_show_home":  []string{"in:1,2"},
	}

	messages := govalidator.MapData{
		"festival_id": []string{
			"required:节日ID为必填项",
			"numeric_between:节日ID必须为数字并且大于等于 1",
		},
		"festival_name": []string{
			"min_cn:节日名称长度必须至少 1 个字符",
			"max_cn:节日名称长度不能超过 10 个字符",
		},
		"festival_date": []string{
			"numeric_between:节日日期格式应为 YYYYMMDD",
		},
		"is_show_home": []string{
			"in:是否显示仅支持 1 或 2",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}

func FestivalDeleteRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"festival_id": []string{"required", "numeric_between:1,"},
	}

	messages := govalidator.MapData{
		"festival_id": []string{
			"required:节日ID为必填项",
			"numeric_between:节日ID必须为数字并且大于等于 1",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}

func FestivalToggleRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"festival_id": []string{"required", "numeric_between:1,"},
		"is_show_home": []string{"required", "in:1,2"},
	}

	messages := govalidator.MapData{
		"festival_id": []string{
			"required:节日ID为必填项",
			"numeric_between:节日ID必须为数字并且大于等于 1",
		},
		"is_show_home": []string{
			"required:是否显示为必填项",
			"in:是否显示仅支持 1 或 2",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}
