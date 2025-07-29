package po

// HzmJob 任务表
type HzmJob struct {
	BasePo
	ExecutorId     *int64  // 执行器id
	Name           *string // 任务名称
	ScheduleType   *byte   // 调度类型：1-cron表达式；2-极简表达式
	ScheduleValue  *string // 调度值：如果scheduleType是1，则为cron表达式；如果scheduleType是2，则为极简表达式值
	Parameters     *string // 任务参数
	Description    *string // 任务描述
	Head           *string // 负责人
	Status         *byte   // 任务状态：0-未启动；1-已启动
	RouterStrategy *byte   // 路由策略：0-轮询；1-随机；2-故障转移
	CronEntryId    *int    // 注册到cron中的id
}

// JobScheduleType 任务调度类型枚举
type JobScheduleType byte

const (
	JobCron    JobScheduleType = iota + 1 // cron表达式
	JobSipCron                            // 极简表达式
)

// JobStatus 任务状态枚举
type JobStatus byte

func (my JobStatus) Is(status *byte) bool {
	return my == JobStatus(*status)
}

func (my JobStatus) ToPtr() *JobStatus {
	p := new(JobStatus)
	*p = my
	return p
}

const (
	JobStop    JobStatus = iota // 未启动
	JobRunning                  // 已启动
)

// JobRouterStrategy 任务路由策略枚举
type JobRouterStrategy byte

const (
	JobPoll    JobRouterStrategy = iota // 轮询
	JobRandom                           // 随机
	JobErrNext                          // 故障转移
)
