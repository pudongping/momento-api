// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package transaction

import (
	"net/http"

	"github.com/pudongping/momento-api/coreKit/responses"

	"github.com/pudongping/momento-api/internal/logic/transaction"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取交易统计
func TransactionStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TransactionStatsReq
		if err := httpx.Parse(r, &req); err != nil {
			responses.ToParamValidateResponse(r, w, err)
			return
		}

		l := transaction.NewTransactionStatsLogic(r.Context(), svcCtx)
		resp, err := l.TransactionStats(&req)
		responses.ToResponse(r, w, resp, err)
	}
}
