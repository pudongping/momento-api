// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"net/http"

	"github.com/pudongping/momento-api/coreKit/responses"

	"github.com/pudongping/momento-api/internal/logic/user"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取用户信息
func UserInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			responses.ToParamValidateResponse(r, w, err)
			return
		}

		l := user.NewUserInfoLogic(r.Context(), svcCtx)
		resp, err := l.UserInfo(&req)
		responses.ToResponse(r, w, resp, err)
	}
}
