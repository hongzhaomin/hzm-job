package iface

type ExecutorSecretCacheIface interface {

	// GetSecretByAppKey 根据执行器标识获取执行器密钥
	GetSecretByAppKey(appKey string) string

	// DeleteByAppKey 根据执行器标识删除缓存
	DeleteByAppKey(appKey string)
}
