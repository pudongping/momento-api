// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package festival

import (
	"context"

	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type FestivalListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取节日列表
func NewFestivalListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FestivalListLogic {
	return &FestivalListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FestivalListLogic) FestivalList() (resp []types.FestivalListResp, err error) {
	userID := ctxData.GetUIDFromCtx(l.ctx)
	userIDUint := cast.ToUint64(userID)

	// 构建查询条件
	queryBuilder := l.svcCtx.FestivalsModel.SelectBuilder().
		Where("user_id = ?", userIDUint)

	// 按festival_date和updated_at倒序排列
	queryBuilder = queryBuilder.OrderBy("festival_date DESC, updated_at DESC")

	festivals, err := l.svcCtx.FestivalsModel.FindAll(l.ctx, queryBuilder)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		l.Logger.Errorf("查询节日列表失败 userID : %d, err : %v", userID, err)
		return nil, errcode.DBError.WithError(errors.Wrapf(err, "FestivalList FindAll userID : %d", userID)).Msgr("获取节日列表失败")
	}

	if festivals == nil {
		return nil, nil
	}

	// 转换为响应格式
	resp = make([]types.FestivalListResp, 0, len(festivals))
	for _, festival := range festivals {
		var item types.FestivalListResp
		item.FestivalId = cast.ToInt64(festival.FestivalId)
		item.UserId = cast.ToString(festival.UserId)
		item.FestivalName = festival.FestivalName
		item.FestivalDate = cast.ToInt64(festival.FestivalDate)
		item.IsShowHome = cast.ToInt32(festival.IsShowHome)
		item.CreatedAt = cast.ToInt64(festival.CreatedAt)
		item.UpdatedAt = cast.ToInt64(festival.UpdatedAt)
		resp = append(resp, item)
	}

	return resp, nil
}
