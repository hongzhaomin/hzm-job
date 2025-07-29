package vo

type Executor struct {
	Id              *int64  `json:"id,omitempty"`           // 执行器id
	Name            *string `json:"name,omitempty"`         // 执行器名称
	AppKey          *string `json:"appKey,omitempty"`       // 执行器标识
	RegistryType    *byte   `json:"registryType,omitempty"` // 注册方式：0-自动；1-手动
	OnlineNodeCount int     `json:"onlineNodeCount"`        // 在线节点数量
}

type ExecutorNode struct {
	Id      *int64  `json:"id,omitempty"`      // 执行器节点id
	Address *string `json:"address,omitempty"` // 节点地址
	Status  *byte   `json:"status,omitempty"`  // 节点状态：0-离线；1-在线
}

type ExecutorSelectBox struct {
	Name  *string `json:"name,omitempty"`  // 执行器名称
	Value *int64  `json:"value,omitempty"` // 执行器id
}
