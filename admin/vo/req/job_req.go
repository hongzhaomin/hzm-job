package req

type JobPage struct {
	BasePage

	ExecutorId  *int64  `json:"executorId,omitempty" form:"executorId"`   // 执行器id
	Status      *byte   `json:"status,omitempty" form:"status"`           // 任务状态
	Name        string  `json:"name,omitempty" form:"name"`               // 任务名称
	Description string  `json:"description,omitempty" form:"description"` // 任务描述
	Head        string  `json:"head,omitempty" form:"head"`               // 负责人
	ExecutorIds []int64 `json:"executorIds,omitempty"`                    // 执行器id，数据权限
}

type Job struct {
	Id             *int64  `json:"id,omitempty" form:"id"`                                            // 任务id
	ExecutorId     *int64  `json:"executorId,omitempty" form:"executorId" binding:"required"`         // 执行器id
	Name           *string `json:"name,omitempty" form:"name" binding:"required"`                     // 任务名称
	ScheduleType   *byte   `json:"scheduleType,omitempty" form:"scheduleType" binding:"required"`     // 调度类型：1-cron表达式；2-极简表达式
	ScheduleValue  *string `json:"scheduleValue,omitempty" form:"scheduleValue" binding:"required"`   // 调度值：如果scheduleType是1，则为cron表达式；如果scheduleType是2，则为极简表达式值
	Parameters     *string `json:"parameters,omitempty" form:"parameters"`                            // 任务参数
	Description    *string `json:"description,omitempty" form:"description" binding:"required"`       // 任务描述
	Head           *string `json:"head,omitempty" form:"head" binding:"required"`                     // 负责人
	RouterStrategy *byte   `json:"routerStrategy,omitempty" form:"routerStrategy" binding:"required"` // 路由策略：0-轮询；1-随机；2-故障转移
}

type JobRunOnce struct {
	Id             *int64 `json:"id,omitempty" form:"id" binding:"required"`      // 任务id
	Parameters     string `json:"parameters,omitempty" form:"parameters"`         // 任务参数
	ExecutorNodeId *int64 `json:"executorNodeId,omitempty" form:"executorNodeId"` // 执行节点id
}
