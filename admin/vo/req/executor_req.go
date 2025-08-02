package req

type ExecutorPage struct {
	BasePage

	Name        string  `json:"name,omitempty" form:"name"`     // 执行器名称
	AppKey      string  `json:"appKey,omitempty" form:"appKey"` // 执行器标识
	ExecutorIds []int64 `json:"executorIds,omitempty"`          // 执行器列表，数据权限
}

type Executor struct {
	Id           *int64  `json:"id,omitempty" form:"id"`                                        // 执行器id
	Name         *string `json:"name,omitempty" form:"name" binding:"required"`                 // 执行器名称
	AppKey       *string `json:"appKey,omitempty" form:"appKey" binding:"required"`             // 执行器标识
	AppSecret    *string `json:"appSecret,omitempty" form:"appSecret"`                          // 执行器密钥，鉴权需要，空的表示不鉴权
	RegistryType *byte   `json:"registryType,omitempty" form:"registryType" binding:"required"` // 注册方式：0-自动；1-手动
	Addresses    *string `json:"addresses,omitempty" form:"addresses"`                          // 执行器节点地址列表，多个英文逗号连接
}
