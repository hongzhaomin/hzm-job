package req

import "time"

type JobLogPage struct {
	BasePage

	ExecutorId        *int64  `json:"executorId,omitempty" form:"executorId"`               // 执行器id
	JobId             *int64  `json:"jobId,omitempty" form:"jobId"`                         // 任务id
	Status            *byte   `json:"status,omitempty" form:"status"`                       // 日志状态：0-待调度；1-任务执行中；2-任务结束
	ScheduleStartTime string  `json:"scheduleStartTime,omitempty" form:"scheduleStartTime"` // 调度开始时间
	ScheduleEndTime   string  `json:"scheduleEndTime,omitempty" form:"scheduleEndTime"`     // 调度结束时间
	ExecutorIds       []int64 `json:"executorIds,omitempty"`                                // 执行器id，数据权限
}

type LogDelParam struct {
	ExecutorId *int64 `json:"executorId,omitempty"`                   // 执行器id
	JobId      *int64 `json:"jobId,omitempty"`                        // 任务id
	ClearType  *int   `json:"clearType,omitempty" binding:"required"` // 清理策略：@see cleartype.ClearType
}

type LogDelDaoParam struct {
	ExecutorId       *int64     // 执行器id
	JobId            *int64     // 任务id
	CreateTimeBefore *time.Time // 创建时间之前
	CountBefore      *int       // 数量之前
}

type StopJobParam struct {
	Id      *int64  `json:"id,omitempty" binding:"required"`      // 调度日志id
	Address *string `json:"address,omitempty" binding:"required"` // 执行器地址
}
