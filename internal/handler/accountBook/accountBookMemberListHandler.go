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

// 获取账本成员列表
func AccountBookMemberListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AccountBookMemberListReq
		if err := httpx.Parse(r, &req); err != nil {
			responses.ToParamValidateResponse(r, w, err)
			return
		}

		if msg, ok := validator.CallValidate(&req, requests.AccountBookMemberListRequestCheck); !ok {
			responses.ToParamValidateResponse(r, w, nil, msg)
			return
		}

		l := accountBook.NewAccountBookMemberListLogic(r.Context(), svcCtx)
		resp, err := l.AccountBookMemberList(&req)
		responses.ToResponse(r, w, resp, err)
	}
}
