package vo

type JobLog struct {
	Id                  *int64  `json:"id,omitempty"`                  // 调度日志id
	JobId               *int64  `json:"jobId,omitempty"`               // 任务id
	JobDescription      *string `json:"jobDescription,omitempty"`      // 任务描述
	JobName             *string `json:"jobName,omitempty"`             // 任务名称
	ExecutorId          *int64  `json:"executorId,omitempty"`          // 执行器id
	ExecutorName        *string `json:"executorName,omitempty"`        // 执行器名称
	ExecutorNodeAddress *string `json:"executorNodeAddress,omitempty"` // 执行器节点地址
	Parameters          *string `json:"parameters,omitempty"`          // 任务参数
	ScheduleTime        *string `json:"scheduleTime,omitempty"`        // 任务调度日志时间
	Status              *byte   `json:"status,omitempty"`              // 任务调度日志状态：0-待调度；1-任务执行中；2-任务结束
	HandleCode          *int    `json:"handleCode,omitempty"`          // 任务结果编码
	HandleMsg           *string `json:"handleMsg,omitempty"`           // 任务结果消息
	FinishTime          *string `json:"finishTime,omitempty"`          // 任务完成时间
}

type ClearTypeSelectBox struct {
	Name  string `json:"name,omitempty"`  // 清理策略名称
	Value int    `json:"value,omitempty"` // 清理策略编号
}
