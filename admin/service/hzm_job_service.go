package service

import (
	"github.com/hongzhaomin/hzm-job/admin/dao"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/po"
)

type HzmJobService struct {
	hzmJobDao dao.HzmJobDao
}

func (my *HzmJobService) FindRunningJobs() []*po.HzmJob {
	jobs, err := my.hzmJobDao.FindRunningJobs()
	if err != nil {
		global.SingletonPool().Log.Error("查询启动中的任务失败", "err", err)
		return nil
	}
	return jobs
}
