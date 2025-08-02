package req

type Login struct {
	UserName string `json:"userName" form:"userName" binding:"required"` // 用户名
	Password string `json:"password" form:"password" binding:"required"` // 密码
}

type ApiAuthParam struct {
	AppKey *string `json:"appKey,omitempty"` // 执行器服务名称标识
}
