// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"net/http"

	"github.com/pudongping/momento-api/coreKit/responses"

	"github.com/pudongping/momento-api/internal/logic/common"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 文件上传
func UploadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadFileReq
		if err := httpx.Parse(r, &req); err != nil {
			responses.ToParamValidateResponse(r, w, err)
			return
		}

		l := common.NewUploadFileLogic(r.Context(), svcCtx)
		resp, err := l.UploadFile(&req, r)
		responses.ToResponse(r, w, resp, err)
	}
}
