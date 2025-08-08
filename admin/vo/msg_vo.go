package vo

type ScheduleStaMsg struct {
	StaDay      string // 统计日期（yyyy-MM-dd）
	TotalIncr   int64  // 总数增量
	SuccessIncr int64  // 成功数增量
	FailIncr    int64  // 失败数增量
}
