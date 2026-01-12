// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tag

import (
	"net/http"

	"github.com/pudongping/momento-api/coreKit/responses"
	"github.com/pudongping/momento-api/coreKit/validator"

	"github.com/pudongping/momento-api/internal/logic/tag"
	"github.com/pudongping/momento-api/internal/requests"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更新自定义标签
func TagUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TagUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			responses.ToParamValidateResponse(r, w, err)
			return
		}

		if msg, ok := validator.CallValidate(&req, requests.TagUpdateRequestCheck); !ok {
			responses.ToParamValidateResponse(r, w, nil, msg)
			return
		}

		l := tag.NewTagUpdateLogic(r.Context(), svcCtx)
		resp, err := l.TagUpdate(&req)
		responses.ToResponse(r, w, resp, err)
	}
}
