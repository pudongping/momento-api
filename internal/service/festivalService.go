package service

import (
	"context"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type FestivalService struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFestivalService(ctx context.Context, svcCtx *svc.ServiceContext) *FestivalService {
	return &FestivalService{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 检查节日是否存在且属于当前用户
func (svc *FestivalService) CheckFestivalOwnership(userID uint64, festivalID int64, operationType string) (*model.Festivals, error) {
	festival, err := svc.svcCtx.FestivalsModel.FindOne(svc.ctx, cast.ToUint64(festivalID))
	if err != nil {
		svc.Logger.Errorf("查询节日失败 userID : %d, festivalID : %d, err : %v", userID, festivalID, err)
		return nil, errcode.DBError.WithError(errors.Wrapf(err, "CheckFestivalOwnership FindOne festivalID : %d", festivalID)).Msgr("节日不存在")
	}

	if festival.UserId != userID {
		return nil, errcode.Fail.Msgr("没有权限" + operationType + "该节日")
	}

	return festival, nil
}

// 验证节日日期是否为有效的YYYYMMDD格式且为未来日期
func (svc *FestivalService) IsValidAndFutureFestivalDate(festivalDate int64) bool {
	dateStr := strconv.FormatInt(festivalDate, 10)
	if len(dateStr) != 8 {
		return false
	}

	year := dateStr[0:4]
	month := dateStr[4:6]
	day := dateStr[6:8]

	dateLayout := "20060102"
	t, err := time.Parse(dateLayout, year+month+day)
	if err != nil {
		svc.Logger.Errorf("解析日期失败 dateStr : %s, err : %v", dateStr, err)
		return false
	}

	// 检查日期是否为未来日期（比今天晚）
	today := time.Now()
	return t.After(today)
}
