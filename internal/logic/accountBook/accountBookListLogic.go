// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package accountBook

import (
	"context"

	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccountBookListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取账本列表
func NewAccountBookListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountBookListLogic {
	return &AccountBookListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccountBookListLogic) AccountBookList(req *types.AccountBookListReq) (resp []types.AccountBookListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
