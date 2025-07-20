package api

// JobClientApi 分布式任务客户端api接口
// 服务端与客户端交互遵循此接口规范
type JobClientApi interface {

	// HeatBeat 服务心跳检测
	HeatBeat()

	// JobHandle 任务处理接口
	JobHandle(req *JobHandleReq)

	// CancelJob 任务终止接口
	CancelJob(req *JobCancelReq)
}

type JobHandleReq struct {
	LogId     *int64  `json:"logId,omitempty"`     // 本次调度日志id
	JobId     *int64  `json:"jobId,omitempty"`     // 任务id
	JobName   *string `json:"jobName,omitempty"`   // 任务名称
	JobParams *string `json:"jobParams,omitempty"` // 任务参数
}

type JobCancelReq struct {
	LogId *int64 `json:"logId,omitempty"` // 需要取消的调度日志id
}
