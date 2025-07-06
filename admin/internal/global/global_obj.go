package global

import (
	"github.com/hongzhaomin/hzm-job/admin/internal/global/iface"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"sync/atomic"
)

var defaultObj atomic.Pointer[Obj]

type Obj struct {
	// fixme logger
	Mysql            *gorm.DB
	Cron             *cron.Cron
	RemotingUtil     *sdk.RemotingUtil
	CronFuncRegister iface.CronFuncRegister
	NodeSelectorMap  map[po.JobRouterStrategy]iface.ExecutorNodeSelector
}

func Store(obj *Obj) {
	defaultObj.Store(obj)
}

func SingletonPool() *Obj {
	return defaultObj.Load()
}
