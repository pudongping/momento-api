// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package festival

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type FestivalToggleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 切换节日显示状态
func NewFestivalToggleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FestivalToggleLogic {
	return &FestivalToggleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FestivalToggleLogic) FestivalToggle(req *types.FestivalToggleReq) (*types.FestivalToggleResp, error) {
	userID := ctxData.GetUIDFromCtx(l.ctx)
	userIDUint := cast.ToUint64(userID)

	// 查询节日是否存在并检查权限
	svc := NewFestivalService(l.ctx, l.svcCtx)
	_, err := svc.CheckFestivalOwnership(userIDUint, req.FestivalId, "修改")
	if err != nil {
		return nil, err
	}

	// 更新显示状态
	updateMap := map[string]interface{}{
		"is_show_home": cast.ToInt64(req.IsShowHome),
		"updated_at":   cast.ToUint64(time.Now().Unix()),
	}

	where := squirrel.Eq{"festival_id": cast.ToUint64(req.FestivalId)}
	if _, err := l.svcCtx.FestivalsModel.UpdateFilter(l.ctx, nil, updateMap, where); err != nil {
		l.Logger.Errorf("切换节日显示状态失败 userID : %d, festivalID : %d, err : %v", userID, req.FestivalId, err)
		return nil, errcode.DBError.WithError(errors.Wrapf(err, "FestivalToggle UpdateFilter festivalID : %d", req.FestivalId)).Msgr("切换节日显示状态失败")
	}

	return &types.FestivalToggleResp{}, nil
}
