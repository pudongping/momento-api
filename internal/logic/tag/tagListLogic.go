// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tag

import (
	"context"

	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取标签列表
func NewTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagListLogic {
	return &TagListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagListLogic) TagList(req *types.TagListReq) (resp []types.TagListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
