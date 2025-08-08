package dao

import (
	"errors"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/admin/vo/req"
	"gorm.io/gorm"
)

type HzmExecutorDao struct{}

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
		Select("Name", "AppKey", "AppSecret", "RegistryType").
		Create(executor).
		Error
}

func (my *HzmExecutorDao) Page(param req.ExecutorPage) (int64, []*po.HzmExecutor, error) {
	// 构造条件
	db := global.SingletonPool().Mysql
	db = db.Where("valid = ?", 1)
	if param.Name != "" {
		db = db.Where("name LIKE ?", "%"+param.Name+"%")
	}
	if param.AppKey != "" {
		db = db.Where("app_key = ?", param.AppKey)
	}
	if len(param.ExecutorIds) > 0 {
		db = db.Where("id in(?)", param.ExecutorIds)
	}

	var count int64
	db.Model(po.HzmExecutor{}).Count(&count)
	if count == 0 {
		return 0, nil, nil
	}

	var executors []*po.HzmExecutor
	err := db.Offset(param.Start()).Limit(param.Limit()).Find(&executors).Error
	return count, executors, err
}

func (my *HzmExecutorDao) Update(executor *po.HzmExecutor) error {
	return global.SingletonPool().Mysql.
		Save(executor).
		Error
}

func (my *HzmExecutorDao) LogicDeleteBatch(ids []int64) error {
	if len(ids) <= 0 {
		return nil
	}
	var executors []*po.HzmExecutor
	err := global.SingletonPool().Mysql.Select("app_key").
		Where("valid = 1 and id in (?)", ids).
		Find(&executors).
		Error
	if err != nil {
		return err
	}
	for _, executor := range executors {
		global.SingletonPool().ExeSecretCache.DeleteByAppKey(*executor.AppKey)
	}

	return global.SingletonPool().Mysql.
		Model(&po.HzmExecutor{}).
		Where("valid = 1 and id in (?)", ids).
		Update("valid", false).
		Error
}

func (my *HzmExecutorDao) CountStatistics() (int64, int64, error) {
	db := global.SingletonPool().Mysql
	var total int64
	err := db.Model(po.HzmExecutor{}).
		Where("valid = 1").
		Count(&total).
		Error
	if err != nil {
		return 0, 0, err
	}

	var online int64
	err = db.Raw("select count(distinct executor_id) from hzm_executor_node where valid = 1 and status = ?", po.NodeOnline).
		Scan(&online).
		Error
	return total, total - online, err
}
