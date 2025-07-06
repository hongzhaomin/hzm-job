package dao

import (
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/po"
)

type HzmExecutorDao struct {
}

func (my *HzmExecutorDao) FindAll() ([]*po.HzmExecutor, error) {
	var executors []*po.HzmExecutor
	err := global.SingletonPool().Mysql.
		Where("valid = ?", true).
		Find(&executors).
		Error
	return executors, err
}

func (my *HzmExecutorDao) FindByIds(ids []int64) ([]*po.HzmExecutor, error) {
	if len(ids) <= 0 {
		return nil, nil
	}

	var executors []*po.HzmExecutor
	err := global.SingletonPool().Mysql.
		Where("valid = ? and id in (?)", true, ids).
		Find(&executors).
		Error
	return executors, err
}

func (my *HzmExecutorDao) FindById(id int64) (*po.HzmExecutor, error) {
	if id <= 0 {
		return nil, nil
	}

	var executor *po.HzmExecutor
	err := global.SingletonPool().Mysql.
		Where("valid = ? and id = ?", true, id).
		First(&executor).
		Error
	return executor, err
}

func (my *HzmExecutorDao) Save(executor *po.HzmExecutor) error {
	return global.SingletonPool().Mysql.Create(executor).Error
}
