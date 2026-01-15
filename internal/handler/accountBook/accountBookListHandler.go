// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package accountBook

import (
	"net/http"

	"github.com/pudongping/momento-api/coreKit/responses"

	"github.com/pudongping/momento-api/internal/logic/accountBook"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取账本列表
func AccountBookListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AccountBookListReq
		if err := httpx.Parse(r, &req); err != nil {
			responses.ToParamValidateResponse(r, w, err)
			return
		}

		l := accountBook.NewAccountBookListLogic(r.Context(), svcCtx)
		resp, err := l.AccountBookList(&req)
		responses.ToResponse(r, w, resp, err)
	}
}
