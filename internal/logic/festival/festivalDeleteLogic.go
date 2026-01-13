// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package festival

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type FestivalDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除节日
func NewFestivalDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FestivalDeleteLogic {
	return &FestivalDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FestivalDeleteLogic) FestivalDelete(req *types.FestivalDeleteReq) (*types.FestivalDeleteResp, error) {
	userID := ctxData.GetUIDFromCtx(l.ctx)
	userIDUint := cast.ToUint64(userID)

	// 查询节日是否存在并检查权限
	svc := NewFestivalService(l.ctx, l.svcCtx)
	_, err := svc.CheckFestivalOwnership(userIDUint, req.FestivalId, "删除")
	if err != nil {
		return nil, err
	}

	// 删除节日
	where := squirrel.Eq{"festival_id": cast.ToUint64(req.FestivalId)}
	if _, err := l.svcCtx.FestivalsModel.DeleteFilter(l.ctx, nil, where); err != nil {
		l.Logger.Errorf("删除节日失败 userID : %d, festivalID : %d, err : %v", userID, req.FestivalId, err)
		return nil, errcode.DBError.WithError(errors.Wrapf(err, "FestivalDelete DeleteFilter festivalID : %d", req.FestivalId)).Msgr("删除节日失败")
	}

	return &types.FestivalDeleteResp{}, nil
}
