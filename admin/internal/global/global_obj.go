package global

import (
	"github.com/hongzhaomin/hzm-job/admin/internal/global/iface"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log/slog"
	"sync/atomic"
)

var defaultObj atomic.Pointer[Obj]

type Obj struct {
	Log              *slog.Logger                                        // 日志对象
	LogLevelVar      *slog.LevelVar                                      // 日志级别配置
	Mysql            *gorm.DB                                            // 数据库对象
	Cron             *cron.Cron                                          // cron定时任务对象
	RemotingUtil     *sdk.RemotingUtil                                   // 远程工具对象
	CronFuncRegister iface.CronFuncRegister                              // 定时任务注册器
	NodeSelectorMap  map[po.JobRouterStrategy]iface.ExecutorNodeSelector // 执行器节点选择器对象列表
	ExeSecretCache   iface.ExecutorSecretCacheIface                      // 执行器密钥缓存对象
}

func Store(obj *Obj) {
	defaultObj.Store(obj)
}

func SingletonPool() *Obj {
	return defaultObj.Load()
}
