package validator

import (
	"net/http"

	"github.com/thedevsaddam/govalidator"
)

type ValidateFunc func(interface{}) map[string][]string

// 结构体中标签标识符
const tagIdentifier = "valid"

// CallValidate 调用验证方法进行验证
// obj 需要被验证的数据值
// handler 验证方法
func CallValidate(obj interface{}, handler ValidateFunc) (string, bool) {
	errs := handler(obj)

	if len(errs) == 0 {
		// 参数验证通过
		return "", true
	}

	var getFirstMsg = func(errs map[string][]string) string {
		for k := range errs {
			return errs[k][0]
		}
		return ""
	}

	return getFirstMsg(errs), false
}

// Validate 验证 form-data, x-www-form-urlencoded 和 query 传参类型的参数
// 如果要验证文件时 ref： https://github.com/thedevsaddam/govalidator/blob/master/doc/FILE_VALIDATION.md
func Validate(r *http.Request, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	opts := govalidator.Options{
		Request:  r,
		Rules:    rules,    // 验证规则
		Messages: messages, // 自定义错误消息
		// RequiredDefault: true,      // 所有的字段都要通过验证规则
		TagIdentifier: tagIdentifier, // 结构体中标签标识符
	}

	return govalidator.New(opts).Validate()
}

// ValidateStruct 验证已有的结构体数据（不支持验证含有指针变量的结构体，但是支持验证 map[string]interface{}）
// document link： https://github.com/thedevsaddam/govalidator/blob/master/doc/STRUCT_VALIDATION.md
func ValidateStruct(data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	opts := govalidator.Options{
		Data:          data,
		Rules:         rules,
		Messages:      messages,
		TagIdentifier: tagIdentifier, // 结构体中标签标识符
	}

	return govalidator.New(opts).ValidateStruct()
}
