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

// 更新节日
func FestivalUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FestivalUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			responses.ToParamValidateResponse(r, w, err)
			return
		}

		if msg, ok := validator.CallValidate(&req, requests.FestivalUpdateRequestCheck); !ok {
			responses.ToParamValidateResponse(r, w, nil, msg)
			return
		}

		l := festival.NewFestivalUpdateLogic(r.Context(), svcCtx)
		resp, err := l.FestivalUpdate(&req)
		responses.ToResponse(r, w, resp, err)
	}
}
