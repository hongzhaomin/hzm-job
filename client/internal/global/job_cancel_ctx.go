package global

import (
	"context"
	"sync"
)

func NewJobCancelCtx() *JobCancelCtx {
	return &JobCancelCtx{
		logId2CancelFuncMap: &sync.Map{},
	}
}

type JobCancelCtx struct {
	logId2CancelFuncMap *sync.Map // map[*int64]context.CancelFunc
}

func (my *JobCancelCtx) Put(logId *int64, cancel context.CancelFunc) {
	if logId != nil && cancel != nil {
		my.logId2CancelFuncMap.Store(*logId, cancel)
	}
}

func (my *JobCancelCtx) CancelAndRemove(logId *int64) {
	cancelFunc := my.getAndRemove(logId)
	if cancelFunc != nil {
		cancelFunc()
	}
}

func (my *JobCancelCtx) getAndRemove(logId *int64) context.CancelFunc {
	if logId == nil {
		return nil
	}
	value, ok := my.logId2CancelFuncMap.LoadAndDelete(logId)
	if !ok {
		return nil
	}
	return value.(context.CancelFunc)
}
