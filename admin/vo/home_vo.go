package vo

import "time"

type DataBlock struct {
	JobTotalNum        int64 `json:"jobTotalNum"`        // 任务总数
	RunningJobNum      int64 `json:"runningJobNum"`      // 运行中的任务数
	ExecutorTotalNum   int64 `json:"executorTotalNum"`   // 执行器总数
	ExecutorOfflineNum int64 `json:"executorOfflineNum"` // 离线执行器数
}

type ScheduleTrend struct {
	StatisticsDate string `json:"statisticsDate"` // 统计日期（yyyy-MM-dd）
	TotalNum       int64  `json:"totalNum"`       // 调度总数
	SuccessNum     int64  `json:"successNum"`     // 调度成功数
	FailNum        int64  `json:"failNum"`        // 调度失败数
}

type OperateLog struct {
	Operator    string   `json:"operator"`    // 操作人
	Description string   `json:"description"` // 操作内容
	OperateTime string   `json:"operateTime"` // 操作时间（yyyy-MM-dd HH:mm:ss）
	Details     []string `json:"details"`     // 操作日志详情
}

type OperateLogMsg struct {
	OperatorId  int64     `json:"operatorId"`  // 操作人
	Description string    `json:"description"` // 操作内容
	OperateTime time.Time `json:"operateTime"` // 操作时间（yyyy-MM-dd HH:mm:ss）
	OldValue    any       `json:"oldValue"`    // 旧数据
	NewValue    any       `json:"newValue"`    // 新数据
}
