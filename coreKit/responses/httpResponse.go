package responses

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

func ToResponse(r *http.Request, w http.ResponseWriter, data interface{}, err error) {
	var resp interface{}
	if err != nil {
		resp = Err(r, w, err)
	} else {
		resp = Success(r, w, data)
	}

	httpx.OkJsonCtx(r.Context(), w, resp)
}

// ToParamValidateResponse 参数验证失败时的返回
func ToParamValidateResponse(r *http.Request, w http.ResponseWriter, err error, msg ...string) {
	e := errcode.UnprocessableEntity
	if err != nil {
		e = e.Msgf("：" + err.Error())
	}
	if len(msg) > 0 && "" != msg[0] {
		e = e.Msgr(msg[0])
	}
	resp := NewErrorResp(e.Code(), e.Msg())

	httpx.OkJsonCtx(r.Context(), w, resp)
}

func mustNullJson(data interface{}) interface{} {
	// isNil() 方法只能用于指向指针类型的指针
	if data == nil || reflect.ValueOf(data).IsNil() {
		return NullJson{}
	}
	return data
}

func Success(r *http.Request, w http.ResponseWriter, data interface{}) *SuccessBean {
	code := errcode.Success
	return NewSuccessResp(code.Code(), code.Msg(), mustNullJson(data))
}

func Err(r *http.Request, w http.ResponseWriter, err error, messages ...string) *ErrorBean {
	var e = errcode.InternalServerError
	var errLog = fmt.Sprintf("【API-ERR】==> %+v \n", err)

	causeErr := errors.Cause(err) // err 类型
	if customErr, ok := causeErr.(*errcode.Error); ok {
		// 如果 err 为自定义错误时
		e = customErr
	} else if grpcStatus, ok := status.FromError(causeErr); ok {
		// 如果为 grpc err 时
		e = errcode.GRPCError
		errLog += fmt.Sprintf("【GRPC】异常，code 为：%+v 错误提示信息为：%s \n", grpcStatus.Code(), grpcStatus.Message())
	}

	resp := NewErrorResp(e.Code(), e.Msg())

	if len(messages) > 0 && messages[0] != "" {
		resp.Msg = messages[0]
	}

	// 如果包含了详细错误提示信息时
	if len(e.Details()) > 0 {
		errLog += fmt.Sprintf("【Err Details】==> %s \n", strings.Join(e.Details(), "|"))
	}
	// 如果包含了原始错误信息时
	if e.Err() != nil {
		errLog += fmt.Sprintf("【Err Stack Trace】==> %+v \n", e.Err())
	}

	logx.WithContext(r.Context()).Error(errLog)

	return resp
}
