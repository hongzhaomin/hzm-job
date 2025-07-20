package dao

import (
	"errors"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"gorm.io/gorm"
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

	var executor po.HzmExecutor
	err := global.SingletonPool().Mysql.
		Where("valid = ? and id = ?", true, id).
		First(&executor).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &executor, err
}

func (my *HzmExecutorDao) FindByAppKey(appKey string) (*po.HzmExecutor, error) {
	if appKey == "" {
		return nil, nil
	}

	var executor po.HzmExecutor
	err := global.SingletonPool().Mysql.
		Where("valid = ? and app_key = ?", true, appKey).
		First(&executor).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &executor, err
}

func (my *HzmExecutorDao) Save(executor *po.HzmExecutor) error {
	return global.SingletonPool().Mysql.
		// fixme
		//Select("").
		Create(executor).
		Error
}
