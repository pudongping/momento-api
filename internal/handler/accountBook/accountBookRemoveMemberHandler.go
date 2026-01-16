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

// 移除账本成员
func AccountBookRemoveMemberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AccountBookRemoveMemberReq
		if err := httpx.Parse(r, &req); err != nil {
			responses.ToParamValidateResponse(r, w, err)
			return
		}

		if msg, ok := validator.CallValidate(&req, requests.AccountBookRemoveMemberRequestCheck); !ok {
			responses.ToParamValidateResponse(r, w, nil, msg)
			return
		}

		l := accountBook.NewAccountBookRemoveMemberLogic(r.Context(), svcCtx)
		resp, err := l.AccountBookRemoveMember(&req)
		responses.ToResponse(r, w, resp, err)
	}
}
