package global

import (
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"log/slog"
	"sync/atomic"
)

var defaultObj atomic.Pointer[Obj]

type Obj struct {
	Log          *slog.Logger      // 日志对象
	LogLevelVar  *slog.LevelVar    // 日志级别配置
	RemotingUtil *sdk.RemotingUtil // 远程工具对象
	JobCancelCtx *JobCancelCtx     // 任务取消容器
}

func Store(obj *Obj) {
	defaultObj.Store(obj)
}

func SingletonPool() *Obj {
	return defaultObj.Load()
}
