package dao

import (
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/po"
)

type HzmJobDao struct{}

func (my *HzmJobDao) Save(job *po.HzmJob) error {
	return global.SingletonPool().Mysql.Create(job).Error
}

func (my *HzmJobDao) FindRunningJobs() ([]*po.HzmJob, error) {
	var jobs []*po.HzmJob
	err := global.SingletonPool().Mysql.
		Where("valid = ? and status = ?", true, po.JobRunning).
		Find(&jobs).
		Error
	return jobs, err
}
