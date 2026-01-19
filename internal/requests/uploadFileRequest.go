package requests

import (
	"github.com/pudongping/momento-api/coreKit/validator"
	"github.com/thedevsaddam/govalidator"
)

func UploadFileRequestCheck(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"file_type": []string{"required", "min_cn:1", "max_cn:10"},
	}

	messages := govalidator.MapData{
		"file_type": []string{
			"required:文件类型不能为空",
			"min_cn:文件类型长度不能少于1个字符",
			"max_cn:文件类型长度不能超过10个字符",
		},
	}

	return validator.ValidateStruct(data, rules, messages)
}
