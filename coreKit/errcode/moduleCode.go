// Package errcode 业务相关的错误码
package errcode

var (
	ErrorUploadFileFail = NewError(200000, "上传文件失败")
	NoData              = NewError(200001, "没有数据")
)
