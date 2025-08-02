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
	AppKey          *string `json:"appKey,omitempty"`          // 执行器服务名称标识
	ExecutorAddress *string `json:"executorAddress,omitempty"` // 执行器地址（ip+端口）
}

type JobResultReq struct {
	LogId       *int64  `json:"logId,omitempty"`       // 任务日志id
	AppKey      *string `json:"appKey,omitempty"`      // 执行器服务名称标识
	HandlerCode *int    `json:"handlerCode,omitempty"` // 任务处理编码，200标识成功，其他失败
	HandlerMsg  *string `json:"handlerMsg,omitempty"`  // 任务处理结果消息
}
