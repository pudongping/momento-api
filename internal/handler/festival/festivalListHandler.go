// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package festival

import (
	"net/http"

	"github.com/pudongping/momento-api/coreKit/responses"

	"github.com/pudongping/momento-api/internal/logic/festival"
	"github.com/pudongping/momento-api/internal/svc"
)

// 获取节日列表
func FestivalListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := festival.NewFestivalListLogic(r.Context(), svcCtx)
		resp, err := l.FestivalList()
		responses.ToResponse(r, w, resp, err)
	}
}
