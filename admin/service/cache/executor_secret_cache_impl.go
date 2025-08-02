package cache

import (
	"github.com/hongzhaomin/hzm-job/admin/dao"
	"github.com/hongzhaomin/hzm-job/admin/internal/global/iface"
	"sync"
	"time"
)

type secretObj struct {
	appSecret string // 密钥
	cacheTime int64  // 缓存时间（毫秒时间戳）
}

var _ iface.ExecutorSecretCacheIface = (*ExecutorSecretCacheImpl)(nil)

func NewExecutorSecretCache() *ExecutorSecretCacheImpl {
	return &ExecutorSecretCacheImpl{
		appKey2SecretMap: &sync.Map{},
		lock:             &sync.Mutex{},
		survivalTime:     time.Minute,
	}
}

type ExecutorSecretCacheImpl struct {
	appKey2SecretMap *sync.Map
	lock             *sync.Mutex
	survivalTime     time.Duration
	hzmExecutorDao   dao.HzmExecutorDao
}

func (my *ExecutorSecretCacheImpl) GetSecretByAppKey(appKey string) string {
	if appKey == "" {
		return ""
	}
	if obj, ok := my.appKey2SecretMap.Load(appKey); ok {
		s := obj.(*secretObj)
		if !my.expired(s) {
			return s.appSecret
		}
	}

	my.lock.Lock()
	defer my.lock.Unlock()
	if obj, ok := my.appKey2SecretMap.Load(appKey); ok {
		s := obj.(*secretObj)
		if !my.expired(s) {
			return s.appSecret
		}
	}

	// 查询数据库
	executor, _ := my.hzmExecutorDao.FindByAppKey(appKey)
	if executor == nil {
		return ""
	}

	var appSecret string
	if executor.AppSecret != nil && *executor.AppSecret != "" {
		appSecret = *executor.AppSecret
	}
	my.put(appKey, appSecret)
	return appSecret
}

func (my *ExecutorSecretCacheImpl) DeleteByAppKey(appKey string) {
	my.lock.Lock()
	defer my.lock.Unlock()
	my.appKey2SecretMap.Delete(appKey)
}

func (my *ExecutorSecretCacheImpl) expired(attr *secretObj) bool {
	if my.survivalTime == 0 {
		// 存活时间为0，表示永不过期
		return false
	}
	return attr.cacheTime+my.survivalTime.Milliseconds() < time.Now().UnixMilli()
}

func (my *ExecutorSecretCacheImpl) put(appKey string, appSecret string) {
	if appKey != "" && appSecret != "" {
		my.appKey2SecretMap.Store(appKey, &secretObj{
			appSecret: appSecret,
			cacheTime: time.Now().UnixMilli(),
		})
	}
}
