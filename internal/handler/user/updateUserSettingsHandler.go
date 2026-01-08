package user

import (
	"net/http"

	"github.com/pudongping/momento-api/coreKit/responses"
	"github.com/pudongping/momento-api/coreKit/validator"

	"github.com/pudongping/momento-api/internal/logic/user"
	"github.com/pudongping/momento-api/internal/requests"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更新用户设置
func UpdateUserSettingsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateUserSettingsReq
		if err := httpx.Parse(r, &req); err != nil {
			responses.ToParamValidateResponse(r, w, err)
			return
		}

		if msg, ok := validator.CallValidate(&req, requests.UpdateUserSettingsRequestCheck); !ok {
			responses.ToParamValidateResponse(r, w, nil, msg)
			return
		}

		l := user.NewUpdateUserSettingsLogic(r.Context(), svcCtx)
		resp, err := l.UpdateUserSettings(&req)
		responses.ToResponse(r, w, resp, err)
	}
}
