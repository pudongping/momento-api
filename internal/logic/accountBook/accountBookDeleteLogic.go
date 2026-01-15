// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package accountBook

import (
	"context"

	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccountBookDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除账本
func NewAccountBookDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountBookDeleteLogic {
	return &AccountBookDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccountBookDeleteLogic) AccountBookDelete(req *types.AccountBookDeleteReq) (resp *types.AccountBookDeleteResp, err error) {
	// todo: add your logic here and delete this line

	return
}
