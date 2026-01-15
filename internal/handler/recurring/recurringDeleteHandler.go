// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package recurring

import (
	"net/http"

	"github.com/pudongping/momento-api/coreKit/responses"
	"github.com/pudongping/momento-api/coreKit/validator"
	"github.com/pudongping/momento-api/internal/requests"

	"github.com/pudongping/momento-api/internal/logic/recurring"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除周期性记账规则
func RecurringDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RecurringDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			responses.ToParamValidateResponse(r, w, err)
			return
		}

		// 参数校验
		if msg, ok := validator.CallValidate(&req, requests.RecurringDeleteRequestCheck); !ok {
			responses.ToParamValidateResponse(r, w, nil, msg)
			return
		}

		l := recurring.NewRecurringDeleteLogic(r.Context(), svcCtx)
		err := l.RecurringDelete(&req)
		responses.ToResponse(r, w, nil, err)
	}
}
