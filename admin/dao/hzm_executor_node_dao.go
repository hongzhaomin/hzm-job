package dao

import (
	"errors"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"gorm.io/gorm"
)

type HzmExecutorNodeDao struct {
}

func (my *HzmExecutorNodeDao) FindOnlineByExecutorIds(status *po.NodeStatus, executorIds ...int64) ([]*po.HzmExecutorNode, error) {
	if len(executorIds) <= 0 {
		return nil, nil
	}

	var nodes []*po.HzmExecutorNode
	db := global.SingletonPool().Mysql.
		Where("valid = 1 and executor_id IN (?)", executorIds)
	if status != nil {
		db = db.Where("status = ?", status)
	}
	err := db.Order("id asc").Find(&nodes).Error
	return nodes, err
}

func (my *HzmExecutorNodeDao) Save(executorNode *po.HzmExecutorNode) error {
	return global.SingletonPool().Mysql.
		Select("ExecutorId", "Address", "Status").
		Create(executorNode).
		Error
}

func (my *HzmExecutorNodeDao) Delete(id *int64) error {
	if id == nil {
		return nil
	}
	return global.SingletonPool().Mysql.Unscoped().Delete(&po.HzmExecutorNode{}, id).Error
}

func (my *HzmExecutorNodeDao) UpdateStatus(id *int64, status po.NodeStatus) error {
	if id == nil {
		return nil
	}
	return global.SingletonPool().Mysql.Model(&po.HzmExecutorNode{}).
		Where("id = ? and valid = ?", id, true).
		Update("status", status).Error
}

func (my *HzmExecutorNodeDao) DeleteByAddress(executorId *int64, address *string) error {
	if executorId == nil || address == nil {
		return nil
	}
	return global.SingletonPool().Mysql.Unscoped().
		Where("valid = ? and executor_id = ? and address = ?", true, executorId, address).
		Delete(&po.HzmExecutorNode{}).
		Error
}

func (my *HzmExecutorNodeDao) IsOnline(id *int64) bool {
	if id == nil {
		return false
	}

	var node po.HzmExecutorNode
	err := global.SingletonPool().Mysql.
		Select("id").
		Where("valid = 1 and id = ? and status = ?", id, po.NodeOnline).
		Find(&node).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return node.Id != nil
}
