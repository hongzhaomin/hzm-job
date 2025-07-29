package vo

type Job struct {
	Id             *int64  `json:"id,omitempty"`             // 任务id
	ExecutorId     *int64  `json:"executorId,omitempty"`     // 执行器id
	Name           *string `json:"name,omitempty"`           // 任务名称
	ScheduleType   *byte   `json:"scheduleType,omitempty"`   // 调度类型：1-cron表达式；2-极简表达式
	ScheduleValue  *string `json:"scheduleValue,omitempty"`  // 调度值：如果scheduleType是1，则为cron表达式；如果scheduleType是2，则为极简表达式值
	Parameters     *string `json:"parameters,omitempty"`     // 任务参数
	Description    *string `json:"description,omitempty"`    // 任务描述
	Head           *string `json:"head,omitempty"`           // 负责人
	Status         *byte   `json:"status,omitempty"`         // 任务状态：0-未启动；1-已启动
	RouterStrategy *byte   `json:"routerStrategy,omitempty"` // 路由策略：0-轮询；1-随机；2-故障转移
}

type SimpleCronSelectBox struct {
	Name  string `json:"name,omitempty"`  // 表达式描述
	Value string `json:"value,omitempty"` // 表达式名称
}

type JobSelectBox struct {
	Name  *string `json:"name,omitempty"`  // 任务描述
	Value *int64  `json:"value,omitempty"` // 任务id
}
