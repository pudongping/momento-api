package responses

type NullJson struct {
}

type SuccessBean struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ErrorBean struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewSuccessResp(code int, msg string, data interface{}) *SuccessBean {
	return &SuccessBean{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func NewErrorResp(code int, msg string) *ErrorBean {
	return &ErrorBean{
		Code: code,
		Msg:  msg,
	}
}
