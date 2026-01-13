package user

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4/request"
	"github.com/pudongping/momento-api/coreKit/responses"

	"github.com/pudongping/momento-api/internal/logic/user"
	"github.com/pudongping/momento-api/internal/svc"
)

// 用户退出登录
func UserLogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取 token
		token, err := request.BearerExtractor{}.ExtractToken(r)
		if err != nil {
			responses.ToResponse(r, w, nil, err)
			return
		}

		l := user.NewUserLogoutLogic(r.Context(), svcCtx)
		err = l.UserLogout(token)
		responses.ToResponse(r, w, nil, err)
	}
}
