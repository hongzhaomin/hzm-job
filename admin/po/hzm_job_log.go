package po

import "time"

// HzmJobLog 任务日志表
type HzmJobLog struct {
	BasePo
	JobId          *int64     // 任务id
	ExecutorId     *int64     // 执行器id
	ExecutorNodeId *int64     // 执行器节点id
	Parameters     *string    // 任务参数
	ScheduleTime   *time.Time // 任务调度日志时间
	Status         *byte      // 任务调度日志状态：0-待调度；1-任务执行中；2-任务结束
	HandleCode     *int       // 任务结果编码
	HandleMsg      *string    // 任务结果消息
	FinishTime     *time.Time // 任务完成时间
}

// LogStatus 任务日志状态枚举
type LogStatus byte

func (my LogStatus) Is(status *byte) bool {
	return my == LogStatus(*status)
}

func (my LogStatus) ToPtr() *LogStatus {
	p := new(LogStatus)
	*p = my
	return p
}

const (
	LogToSchedule  LogStatus = iota // 待调度
	LogJobRunning                   // 任务执行中
	LogJobFinished                  // 任务结束
)
