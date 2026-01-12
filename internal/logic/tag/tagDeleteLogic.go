// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tag

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除标签
func NewTagDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagDeleteLogic {
	return &TagDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagDeleteLogic) TagDelete(req *types.TagDeleteReq) (*types.TagDeleteResp, error) {
	userID := ctxData.GetUIDFromCtx(l.ctx)
	userIDUint := cast.ToUint64(userID)

	// 查询标签是否存在
	tag, err := l.svcCtx.TagsModel.FindOne(l.ctx, cast.ToUint64(req.TagId))
	if err != nil {
		l.Logger.Errorf("查询标签失败 userID : %d, tagID : %d, err : %v", userID, req.TagId, err)
		return nil, errcode.DBError.WithError(errors.Wrapf(err, "TagDelete FindOne tagID : %d", req.TagId)).Msgr("标签不存在")
	}

	// 检查标签是否为系统标签，不允许删除系统标签
	if tag.IsSystem == model.TagsIsSystemYes {
		return nil, errcode.Fail.Msgr("系统标签不允许删除")
	}

	// 检查标签是否属于当前用户，不允许删除其他用户的标签
	if tag.UserId != userIDUint {
		return nil, errcode.Fail.Msgr("没有权限删除该标签")
	}

	// 删除标签
	where := squirrel.Eq{"tag_id": cast.ToUint64(req.TagId)}
	if _, err := l.svcCtx.TagsModel.DeleteFilter(l.ctx, nil, where); err != nil {
		l.Logger.Errorf("删除标签失败 userID : %d, tagID : %d, err : %v", userID, req.TagId, err)
		return nil, errcode.DBError.WithError(errors.Wrapf(err, "TagDelete DeleteFilter tagID : %d", req.TagId)).Msgr("删除标签失败")
	}

	return &types.TagDeleteResp{}, nil
}
