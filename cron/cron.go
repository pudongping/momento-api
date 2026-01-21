package cron

import (
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

func Start(svcCtx *svc.ServiceContext) {
	// 创建 cron 实例
	c := cron.New(cron.WithSeconds()) // 开启秒级控制

	// 添加任务：每分钟执行一次
	// 0 * * * * * (每分钟的第0秒执行)
	recurringJob := NewRecurringJob(svcCtx)
	_, err := c.AddJob("0 * * * * *", recurringJob)
	if err != nil {
		logx.Errorf("Failed to add recurring job: %v", err)
		return
	}

	c.Start()
	logx.Info("Cron scheduler started")
}
