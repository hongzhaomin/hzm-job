package iface

// ExecutorSecretCache 执行器密钥缓存接口
type ExecutorSecretCache interface {

	// GetSecretByAppKey 根据执行器标识获取执行器密钥
	GetSecretByAppKey(appKey string) string

	// DeleteByAppKey 根据执行器标识删除缓存
	DeleteByAppKey(appKey string)
}
