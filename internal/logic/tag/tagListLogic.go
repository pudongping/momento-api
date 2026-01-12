// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tag

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"

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
	userID := ctxData.GetUIDFromCtx(l.ctx)
	userIDUint := cast.ToUint64(userID)

	// 构建查询条件
	queryBuilder := l.svcCtx.TagsModel.SelectBuilder()

	// 构建type过滤条件：如果type不为空，则按type过滤，否则查询所有type
	if req.Type != "" {
		queryBuilder = queryBuilder.Where("type = ?", req.Type)
	}

	// 查询系统标签和用户自定义标签
	// 系统标签（is_system=1）和用户自定义标签（is_system=2且user_id=当前用户ID）
	queryBuilder = queryBuilder.Where(
		"(is_system = ? OR (is_system = ? AND user_id = ?))",
		model.TagsIsSystemYes,
		model.TagsIsSystemNo,
		userIDUint,
	)

	// 按is_system升序（系统标签优先，is_system=1排在前），然后按sort_num倒序，最后按updated_at倒序
	queryBuilder = queryBuilder.OrderBy("is_system ASC, sort_num DESC, updated_at DESC")

	tags, err := l.svcCtx.TagsModel.FindAll(l.ctx, queryBuilder)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errcode.DBError.WithError(errors.Wrapf(err, "TagList FindAll")).Msgr("获取标签列表失败")
	}

	if tags == nil {
		return nil, nil
	}

	// 转换为响应格式
	resp = make([]types.TagListResp, 0, len(tags))
	for _, tag := range tags {
		var item types.TagListResp
		if err = copier.Copy(&item, tag); err != nil {
			return nil, errors.Wrap(err, "TagList copier")
		}
		item.UserId = cast.ToString(tag.UserId)
		item.TagId = cast.ToInt64(tag.TagId)
		item.SortNum = cast.ToInt64(tag.SortNum)
		resp = append(resp, item)
	}

	return resp, nil
}
