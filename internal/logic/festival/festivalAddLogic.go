// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package festival

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/service"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type FestivalAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 添加节日
func NewFestivalAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FestivalAddLogic {
	return &FestivalAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FestivalAddLogic) FestivalAdd(req *types.FestivalAddReq) (*types.FestivalAddResp, error) {
	userID := ctxData.GetUIDFromCtx(l.ctx)
	userIDUint := cast.ToUint64(userID)

	// 检查用户节日数量是否已达到上限（20个）
	countBuilder := l.svcCtx.FestivalsModel.CountBuilder().
		Where("user_id = ?", userIDUint)

	count, err := l.svcCtx.FestivalsModel.FindCount(l.ctx, countBuilder)
	if err != nil {
		l.Logger.Errorf("查询用户节日数量失败 userID : %d, err : %v", userID, err)
		return nil, errcode.DBError.WithError(errors.Wrapf(err, "FestivalAdd FindCount userID : %d", userID)).Msgr("查询节日数量失败")
	}

	if count >= 20 {
		return nil, errcode.Fail.Msgr("节日数量已达到上限(20个)")
	}

	// 检查是否存在同名的节日
	nameCheckBuilder := l.svcCtx.FestivalsModel.CountBuilder().
		Where("user_id = ?", userIDUint).
		Where("festival_name = ?", strings.TrimSpace(req.FestivalName))

	nameCount, err := l.svcCtx.FestivalsModel.FindCount(l.ctx, nameCheckBuilder)
	if err != nil {
		return nil, errcode.DBError.WithError(errors.Wrapf(err, "FestivalAdd FindCount by name userID : %d, name: %s", userID, req.FestivalName)).Msgr("查询节日名称失败")
	}
	if nameCount > 0 {
		return nil, errcode.Fail.Msgr("节日名称已存在，请勿重复添加")
	}

	// 检查节日日期是否为未来日期
	svc := service.NewFestivalService(l.ctx, l.svcCtx)
	if !svc.IsValidAndFutureFestivalDate(req.FestivalDate) {
		return nil, errcode.Fail.Msgr("节日日期必须是未来的日期")
	}

	isShowHome := req.IsShowHome
	if isShowHome == 0 {
		isShowHome = model.FestivalIsShowHomeYes // 默认为1-是
	}

	now := time.Now().Unix()

	// 新增节日
	festival := &model.Festivals{
		UserId:       userIDUint,
		FestivalName: strings.TrimSpace(req.FestivalName),
		FestivalDate: cast.ToUint64(req.FestivalDate),
		IsShowHome:   cast.ToInt64(isShowHome),
		CreatedAt:    cast.ToUint64(now),
		UpdatedAt:    cast.ToUint64(now),
	}

	_, err = l.svcCtx.FestivalsModel.Insert(l.ctx, festival)
	if err != nil {
		l.Logger.Errorf("添加节日失败 userID : %d, err : %v", userID, err)
		return nil, errcode.DBError.WithError(errors.Wrapf(err, "FestivalAdd Insert userID : %d", userID)).Msgr("添加节日失败")
	}

	return &types.FestivalAddResp{}, nil
}
