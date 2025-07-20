package api

import (
	"errors"
	"github.com/hongzhaomin/hzm-job/admin/dao"
	"github.com/hongzhaomin/hzm-job/admin/po"
)

var _ JobServerApi = (*JobServerApiImpl)(nil)

type JobServerApiImpl struct {
	hzmExecutorDao     dao.HzmExecutorDao
	hzmExecutorNodeDao dao.HzmExecutorNodeDao
	hzmJobLogDao       dao.HzmJobLogDao
}

func (my *JobServerApiImpl) Registry(req *RegistryReq) {
	executor, err := my.hzmExecutorDao.FindByAppKey(*req.AppKey)
	if err != nil {
		panic(err)
	}
	if executor == nil {
		panic(errors.New("executor not exist"))
	}
	if po.ManualRegistry.Is(executor.RegistryType) {
		// 手动注册则忽略
		return
	}

	err = my.hzmExecutorNodeDao.Save(&po.HzmExecutorNode{
		ExecutorId: executor.Id,
		Address:    req.ExecutorAddress,
		Status:     (*byte)(po.NodeOnline.ToPtr()),
	})
	if err != nil {
		panic(err)
	}
}

func (my *JobServerApiImpl) Offline(req *RegistryReq) {
	executor, err := my.hzmExecutorDao.FindByAppKey(*req.AppKey)
	if err != nil {
		panic(err)
	}
	if executor == nil {
		panic(errors.New("executor not exist"))
	}
	if po.ManualRegistry.Is(executor.RegistryType) {
		// 手动注册则忽略
		return
	}

	err = my.hzmExecutorNodeDao.DeleteByAddress(executor.Id, req.ExecutorAddress)
	if err != nil {
		panic(err)
	}
}

func (my *JobServerApiImpl) Callback(req *JobResultReq) {
	// 任务完成回调，更新调度日志记录
	jobLog := &po.HzmJobLog{
		BasePo: po.BasePo{
			Id: req.LogId,
		},
		HandleCode: req.HandlerCode,
		HandleMsg:  req.HandlerMsg,
	}
	if err := my.hzmJobLogDao.FinishJobLogById(jobLog); err != nil {
		panic(err)
	}
}
