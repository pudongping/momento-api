package user

import (
	"net/http"

	"github.com/pudongping/momento-api/coreKit/responses"

	"github.com/pudongping/momento-api/internal/logic/user"
	"github.com/pudongping/momento-api/internal/svc"
)

// 获取用户设置
func UserSettingsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserSettingsLogic(r.Context(), svcCtx)
		resp, err := l.UserSettings()
		responses.ToResponse(r, w, resp, err)
	}
}
