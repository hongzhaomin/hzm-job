package api

// JobServerApi 分布式任务服务端api接口
// 客户端与服务端交互遵循此接口规范
type JobServerApi interface {

	// Registry 客户端注册接口
	Registry(req *RegistryReq)

	// Offline 客户端下线接口
	Offline(req *RegistryReq)

	// Callback 回调接口
	Callback(req *JobResultReq)
}

type RegistryReq struct {
	ExecutorAddress *string `json:"executorAddress,omitempty"` // 执行器地址（ip+端口）
	ExecutorName    *string `json:"executorName,omitempty"`    // 执行器服务名称
}

type JobResultReq struct {
	JobId       *int64  `json:"jobId,omitempty"`       // 任务id
	HandlerCode *int    `json:"handlerCode,omitempty"` // 任务处理编码，200标识成功，其他失败
	HandlerMsg  *string `json:"handlerMsg,omitempty"`  // 任务处理结果消息
}
