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

// 接受邀请
func AccountBookAcceptHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AccountBookAcceptReq
		if err := httpx.Parse(r, &req); err != nil {
			responses.ToParamValidateResponse(r, w, err)
			return
		}

		if msg, ok := validator.CallValidate(&req, requests.AccountBookAcceptRequestCheck); !ok {
			responses.ToParamValidateResponse(r, w, nil, msg)
			return
		}

		l := accountBook.NewAccountBookAcceptLogic(r.Context(), svcCtx)
		resp, err := l.AccountBookAccept(&req)
		responses.ToResponse(r, w, resp, err)
	}
}
