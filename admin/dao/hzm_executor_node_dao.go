package dao

import (
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/po"
)

type HzmExecutorNodeDao struct {
}

func (my *HzmExecutorNodeDao) FindOnlineByExecutorIds(status *po.NodeStatus, executorIds ...int64) ([]*po.HzmExecutorNode, error) {
	if len(executorIds) <= 0 {
		return nil, nil
	}

	var nodes []*po.HzmExecutorNode
	db := global.SingletonPool().Mysql.
		Where("valid = ? and executor_id IN (?)", true, executorIds)
	if status != nil {
		db = db.Where("status = ?", status)
	}
	err := db.Order("id asc").Find(&nodes).Error
	return nodes, err
}

func (my *HzmExecutorNodeDao) Save(executorNode *po.HzmExecutorNode) error {
	return global.SingletonPool().Mysql.Create(executorNode).Error
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
