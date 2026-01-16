// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package accountBook

import (
	"net/http"

	"github.com/pudongping/momento-api/coreKit/responses"
	"github.com/pudongping/momento-api/coreKit/validator"
	"github.com/pudongping/momento-api/internal/requests"

	"github.com/pudongping/momento-api/internal/logic/accountBook"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 设置默认账本
func AccountBookSetDefaultHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AccountBookSetDefaultReq
		if err := httpx.Parse(r, &req); err != nil {
			responses.ToParamValidateResponse(r, w, err)
			return
		}

		if msg, ok := validator.CallValidate(&req, requests.AccountBookSetDefaultRequestCheck); !ok {
			responses.ToParamValidateResponse(r, w, nil, msg)
			return
		}

		l := accountBook.NewAccountBookSetDefaultLogic(r.Context(), svcCtx)
		resp, err := l.AccountBookSetDefault(&req)
		responses.ToResponse(r, w, resp, err)
	}
}
