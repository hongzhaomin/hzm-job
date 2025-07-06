package service

import (
	"github.com/hongzhaomin/hzm-job/admin/dao"
	"github.com/hongzhaomin/hzm-job/admin/po"
)

type HzmJobService struct {
	hzmJobDao dao.HzmJobDao
}

func (my *HzmJobService) FindRunningJobs() []*po.HzmJob {
	jobs, err := my.hzmJobDao.FindRunningJobs()
	if err != nil {
		// todo 打印日志
		return nil
	}
	return jobs
}
