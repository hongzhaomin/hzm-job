package po

import "time"

// HzmScheduleStatistics 任务调度统计表
type HzmScheduleStatistics struct {
	BasePo
	Day        *time.Time // 统计日期（yyyy-MM-dd）
	TotalNum   *int64     // 调度总数
	SuccessNum *int64     // 调度成功数
	FailNum    *int64     // 调度失败数
}
