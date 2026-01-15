// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package festival

import (
	"context"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/service"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type FestivalUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新节日
func NewFestivalUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FestivalUpdateLogic {
	return &FestivalUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FestivalUpdateLogic) FestivalUpdate(req *types.FestivalUpdateReq) (*types.FestivalUpdateResp, error) {
	userID := ctxData.GetUIDFromCtx(l.ctx)
	userIDUint := cast.ToUint64(userID)

	// 查询节日是否存在并检查权限
	svc := service.NewFestivalService(l.ctx, l.svcCtx)
	_, err := svc.CheckFestivalOwnership(userIDUint, req.FestivalId, "修改")
	if err != nil {
		return nil, err
	}

	// 仅更新非空字段
	updateMap := map[string]interface{}{}
	if s := strings.TrimSpace(req.FestivalName); s != "" {
		// 检查是否存在同名的节日
		nameCheckBuilder := l.svcCtx.FestivalsModel.CountBuilder().
			Where("user_id = ?", userIDUint).
			Where("festival_name = ?", s).
			Where("festival_id != ?", req.FestivalId)
		nameCount, err := l.svcCtx.FestivalsModel.FindCount(l.ctx, nameCheckBuilder)
		if err != nil {
			return nil, errcode.DBError.WithError(errors.Wrapf(err, "FestivalUpdate FindCount by name userID : %d, name: %s", userID, s)).Msgr("查询节日名称失败")
		}
		if nameCount > 0 {
			return nil, errcode.Fail.Msgr("节日名称已存在，请更换节日名称")
		}
		updateMap["festival_name"] = s
	}
	if req.FestivalDate > 0 {
		// 检查节日日期是否为未来日期
		if !svc.IsValidAndFutureFestivalDate(req.FestivalDate) {
			return nil, errcode.Fail.Msgr("节日日期必须是未来的日期")
		}
		updateMap["festival_date"] = cast.ToUint64(req.FestivalDate)
	}
	if req.IsShowHome > 0 {
		updateMap["is_show_home"] = cast.ToInt64(req.IsShowHome)
	}

	// 如果全部为空，则不做更新
	if len(updateMap) == 0 {
		return nil, errcode.Fail.Msgr("没有可更新的节日信息")
	}

	// 补充更新时间
	now := time.Now().Unix()
	updateMap["updated_at"] = cast.ToUint64(now)

	where := squirrel.Eq{"festival_id": cast.ToUint64(req.FestivalId)}
	if _, err := l.svcCtx.FestivalsModel.UpdateFilter(l.ctx, nil, updateMap, where); err != nil {
		l.Logger.Errorf("更新节日失败 userID : %d, festivalID : %d, err : %v", userID, req.FestivalId, err)
		return nil, errcode.DBError.WithError(errors.Wrapf(err, "FestivalUpdate UpdateFilter festivalID : %d", req.FestivalId)).Msgr("更新节日失败")
	}

	return &types.FestivalUpdateResp{}, nil
}
