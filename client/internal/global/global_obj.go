package global

import (
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"sync/atomic"
)

var defaultObj atomic.Pointer[Obj]

type Obj struct {
	// fixme logger
	RemotingUtil *sdk.RemotingUtil
}

func Store(obj *Obj) {
	defaultObj.Store(obj)
}

func Pool() *Obj {
	return defaultObj.Load()
}
