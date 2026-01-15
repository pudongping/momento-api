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

// 获取交易流水列表
func TransactionListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TransactionListReq
		if err := httpx.Parse(r, &req); err != nil {
			responses.ToParamValidateResponse(r, w, err)
			return
		}

		l := transaction.NewTransactionListLogic(r.Context(), svcCtx)
		resp, err := l.TransactionList(&req)
		responses.ToResponse(r, w, resp, err)
	}
}
