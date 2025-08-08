package dao

import (
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"time"
)

type HzmOperateLogDao struct{}

func (my *HzmOperateLogDao) Insert(opeLog *po.HzmOperateLog) error {
	return global.SingletonPool().Mysql.
		Select("OperatorId", "Description", "Detail", "OperateTime").
		Create(opeLog).
		Error
}

func (my *HzmOperateLogDao) FindList() ([]*po.HzmOperateLog, error) {
	var result []*po.HzmOperateLog
	err := global.SingletonPool().Mysql.
		Where("valid = 1 and operate_time >= ?", time.Now().AddDate(0, 0, -3)).
		Order("id desc").
		Find(&result).
		Error
	return result, err
}
