// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"net/http"

	"github.com/pudongping/momento-api/coreKit/helpers"
	"github.com/pudongping/momento-api/coreKit/responses"
	"github.com/pudongping/momento-api/coreKit/validator"
	"github.com/pudongping/momento-api/internal/requests"

	"github.com/pudongping/momento-api/internal/logic/user"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 小程序授权登录
func UserLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			responses.ToParamValidateResponse(r, w, err)
			return
		}

		if msg, ok := validator.CallValidate(&req, requests.LoginRequestCheck); !ok {
			responses.ToParamValidateResponse(r, w, nil, msg)
			return
		}

		// 客户端的 ip 地址
		clientIP := helpers.GetClientIP(r)

		l := user.NewUserLoginLogic(r.Context(), svcCtx)
		resp, err := l.UserLogin(&req, clientIP)
		responses.ToResponse(r, w, resp, err)
	}
}
