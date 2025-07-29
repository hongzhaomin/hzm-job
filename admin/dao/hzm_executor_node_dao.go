package dao

import (
	"errors"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"gorm.io/gorm"
)

type HzmExecutorNodeDao struct {
}

func (my *HzmExecutorNodeDao) FindByExecutorIds(status *po.NodeStatus, executorIds ...int64) ([]*po.HzmExecutorNode, error) {
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

func (my *HzmExecutorNodeDao) SaveBatch(executorNodes ...*po.HzmExecutorNode) error {
	if len(executorNodes) <= 0 {
		return nil
	}
	return global.SingletonPool().Mysql.
		Select("ExecutorId", "Address", "Status").
		Create(executorNodes).
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
		Where("valid = 1 and id = ?", id).
		Update("status", status).Error
}

func (my *HzmExecutorNodeDao) DeleteByAddress(executorId *int64, address *string) error {
	if executorId == nil || address == nil {
		return nil
	}
	return global.SingletonPool().Mysql.Unscoped().
		Where("valid = 1 and executor_id = ? and address = ?", executorId, address).
		Delete(&po.HzmExecutorNode{}).
		Error
}

func (my *HzmExecutorNodeDao) IsOnline(address *string) bool {
	if address == nil {
		return false
	}

	var node po.HzmExecutorNode
	err := global.SingletonPool().Mysql.
		Select("id").
		Where("valid = 1 and address = ? and status = ?", address, po.NodeOnline).
		Find(&node).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return node.Id != nil
}

func (my *HzmExecutorNodeDao) DeleteByExecutorId(executorId *int64) error {
	if executorId == nil {
		return nil
	}
	return global.SingletonPool().Mysql.Unscoped().
		Where("valid = 1 and executor_id = ?", executorId).
		Delete(&po.HzmExecutorNode{}).
		Error
}

func (my *HzmExecutorNodeDao) LogicDeleteBatchByExecutorIds(executorIds []int64) error {
	if len(executorIds) <= 0 {
		return nil
	}
	return global.SingletonPool().Mysql.
		Model(&po.HzmExecutorNode{}).
		Where("valid = 1 and executor_id in (?)", executorIds).
		Update("valid", false).
		Error
}

func (my *HzmExecutorNodeDao) CountOnlineByExecutorIds(executorIds []int64) (map[int64]int, error) {
	if len(executorIds) <= 0 {
		return nil, nil
	}

	var resultMap []map[string]any
	err := global.SingletonPool().Mysql.
		Model(&po.HzmExecutorNode{}).
		Select("executor_id, Count(*) as count").
		Where("valid = 1 and status = ? and executor_id IN (?)", po.NodeOnline, executorIds).
		Group("executor_id").
		Find(&resultMap).
		Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	result := make(map[int64]int, len(resultMap))
	for _, m := range resultMap {
		executorId := m["executor_id"].(*int64)
		count := m["count"].(int64)
		result[*executorId] = int(count)
	}
	return result, nil
}

func (my *HzmExecutorNodeDao) FindById(id *int64) (*po.HzmExecutorNode, error) {
	if id == nil {
		return nil, nil
	}

	var node po.HzmExecutorNode
	err := global.SingletonPool().Mysql.
		Where("valid = 1 and id = ?", id).
		First(&node).
		Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &node, nil
}

func (my *HzmExecutorNodeDao) FindByIds(ids []int64) ([]*po.HzmExecutorNode, error) {
	if len(ids) <= 0 {
		return nil, nil
	}

	var nodes []*po.HzmExecutorNode
	err := global.SingletonPool().Mysql.
		Where("valid = 1 and id in(?)", ids).
		Find(&nodes).
		Error
	return nodes, err
}
