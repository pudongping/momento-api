// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tag

import (
	"context"
	"strings"
	"time"

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

type TagUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新自定义标签
func NewTagUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagUpdateLogic {
	return &TagUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagUpdateLogic) TagUpdate(req *types.TagUpdateReq) error {
	userID := ctxData.GetUIDFromCtx(l.ctx)
	userIDUint := cast.ToUint64(userID)

	// 查询标签是否存在
	tag, err := l.svcCtx.TagsModel.FindOne(l.ctx, cast.ToUint64(req.TagId))
	if err != nil {
		l.Logger.Errorf("查询标签失败 userID : %d, tagID : %d, err : %v", userID, req.TagId, err)
		return errcode.DBError.WithError(errors.Wrapf(err, "TagUpdate FindOne tagID : %d", req.TagId)).Msgr("标签不存在")
	}

	// 检查标签是否为系统标签，不允许修改系统标签
	if tag.IsSystem == model.TagsIsSystemYes {
		return errcode.Fail.Msgr("系统标签不允许修改")
	}

	// 检查标签是否属于当前用户，不允许修改其他用户的标签
	if tag.UserId != userIDUint {
		return errcode.Fail.Msgr("没有权限修改该标签")
	}

	// 仅更新非空字段
	updateMap := map[string]interface{}{}
	if s := strings.TrimSpace(req.Name); s != "" {
		// 检查是否存在同名的自定义标签
		nameCheckBuilder := l.svcCtx.TagsModel.CountBuilder().
			Where("user_id = ?", userIDUint).
			Where("name = ?", s).
			Where("is_system = ?", model.TagsIsSystemNo).
			Where("tag_id != ?", req.TagId)
		nameCount, err := l.svcCtx.TagsModel.FindCount(l.ctx, nameCheckBuilder)
		if err != nil {
			return errcode.DBError.WithError(errors.Wrapf(err, "TagUpdate FindCount by name userID : %d, name: %s", userID, s)).Msgr("查询标签名称失败")
		}
		if nameCount > 0 {
			return errcode.Fail.Msgr("标签名称已存在，请更换标签名称")
		}
		updateMap["name"] = s
	}
	if s := strings.TrimSpace(req.Color); s != "" {
		updateMap["color"] = s
	}
	if s := strings.TrimSpace(req.Icon); s != "" {
		updateMap["icon"] = s
	}
	if s := strings.TrimSpace(req.Type); s != "" {
		updateMap["type"] = s
	}

	// 如果全部为空，则不做更新
	if len(updateMap) == 0 {
		return errcode.Fail.Msgr("没有可更新的标签信息")
	}

	// 补充更新时间
	now := time.Now().Unix()
	updateMap["updated_at"] = cast.ToUint64(now)

	where := squirrel.Eq{"tag_id": cast.ToUint64(req.TagId)}
	if _, err := l.svcCtx.TagsModel.UpdateFilter(l.ctx, nil, updateMap, where); err != nil {
		l.Logger.Errorf("更新标签失败 userID : %d, tagID : %d, err : %v", userID, req.TagId, err)
		return errcode.DBError.WithError(errors.Wrapf(err, "TagUpdate UpdateFilter tagID : %d", req.TagId)).Msgr("更新标签失败")
	}

	return nil
}
