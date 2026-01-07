package flake

import (
	"time"

	"github.com/sony/sonyflake"
)

var sonyFlake *sonyflake.Sonyflake

func init() {
	var (
		sonyMachineID uint16
		st            time.Time
		err           error
	)

	startTime := "2024-08-20" // 初始化一个开始的时间，表示从这个时间开始算起
	machineID := 1            // 机器 ID
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		panic(err)
	}

	sonyMachineID = uint16(machineID)
	settings := sonyflake.Settings{
		StartTime: st,
		MachineID: func() (uint16, error) { return sonyMachineID, nil },
	}
	sonyFlake = sonyflake.NewSonyflake(settings)
	if sonyFlake == nil {
		panic("sony_flake not created")
	}
}

// GenUniqueID 雪花算法生成唯一分布式ID
func GenUniqueID() int64 {
	id, err := sonyFlake.NextID()
	if err != nil {
		panic(err)
	}

	return int64(id)
}
