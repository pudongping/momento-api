// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tag

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 添加自定义标签
func NewTagAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagAddLogic {
	return &TagAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagAddLogic) TagAdd(req *types.TagAddReq) error {
	userID := ctxData.GetUIDFromCtx(l.ctx)
	userIDUint := cast.ToUint64(userID)

	// 检查用户自定义标签数量是否已达到上限（20个）
	countBuilder := l.svcCtx.TagsModel.CountBuilder().
		Where("user_id = ?", userIDUint).
		Where("is_system = ?", model.TagsIsSystemNo)

	count, err := l.svcCtx.TagsModel.FindCount(l.ctx, countBuilder)
	if err != nil {
		l.Logger.Errorf("查询用户标签数量失败 userID : %d, err : %v", userID, err)
		return errcode.DBError.WithError(errors.Wrapf(err, "TagAdd FindCount userID : %d", userID)).Msgr("查询标签数量失败")
	}

	if count >= 20 {
		return errcode.Fail.Msgr("自定义标签数量已达到上限(20个)")
	}

	now := time.Now().Unix()

	// 新增标签
	tag := &model.Tags{
		UserId:    userIDUint,
		Name:      req.Name,
		Color:     req.Color,
		Icon:      req.Icon,
		IsSystem:  model.TagsIsSystemNo,
		Type:      req.Type,
		SortNum:   0,
		CreatedAt: cast.ToUint64(now),
		UpdatedAt: cast.ToUint64(now),
	}

	_, err = l.svcCtx.TagsModel.Insert(l.ctx, tag)
	if err != nil {
		l.Logger.Errorf("添加标签失败 userID : %d, err : %v", userID, err)
		return errcode.DBError.WithError(errors.Wrapf(err, "TagAdd Insert userID : %d", userID)).Msgr("添加标签失败")
	}

	return nil
}
