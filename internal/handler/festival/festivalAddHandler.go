// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package festival

import (
	"net/http"

	"github.com/pudongping/momento-api/coreKit/responses"
	"github.com/pudongping/momento-api/coreKit/validator"

	"github.com/pudongping/momento-api/internal/logic/festival"
	"github.com/pudongping/momento-api/internal/requests"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 添加节日
func FestivalAddHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FestivalAddReq
		if err := httpx.Parse(r, &req); err != nil {
			responses.ToParamValidateResponse(r, w, err)
			return
		}

		if msg, ok := validator.CallValidate(&req, requests.FestivalAddRequestCheck); !ok {
			responses.ToParamValidateResponse(r, w, nil, msg)
			return
		}

		l := festival.NewFestivalAddLogic(r.Context(), svcCtx)
		resp, err := l.FestivalAdd(&req)
		responses.ToResponse(r, w, resp, err)
	}
}
