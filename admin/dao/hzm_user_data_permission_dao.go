package dao

import (
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/po"
)

type HzmUserDataPermissionDao struct{}

func (my *HzmUserDataPermissionDao) SaveBatch(dataPerms []*po.HzmUserDataPermission) error {
	return global.SingletonPool().Mysql.
		Select("UserId", "ExecutorId").
		Create(dataPerms).
		Error
}

func (my *HzmUserDataPermissionDao) DeleteByUserIds(userIds ...int64) error {
	if len(userIds) <= 0 {
		return nil
	}
	return global.SingletonPool().Mysql.
		Unscoped().
		Delete(&po.HzmUserDataPermission{}, "valid = 1 and user_id in (?)", userIds).
		Error
}

func (my *HzmUserDataPermissionDao) FindExecutorIdsByUserId(userId int64) ([]int64, error) {
	var dataPerms []*po.HzmUserDataPermission
	err := global.SingletonPool().Mysql.
		Select("executor_id").
		Where("valid = 1 and user_id = ?", userId).
		Find(&dataPerms).
		Error
	if err != nil {
		return nil, err
	}

	var executorIds []int64
	for _, dp := range dataPerms {
		executorIds = append(executorIds, *dp.ExecutorId)
	}
	return executorIds, nil
}

func (my *HzmUserDataPermissionDao) LogicDeleteBatchByExecutorIds(executorIds []int64) error {
	if len(executorIds) <= 0 {
		return nil
	}
	return global.SingletonPool().Mysql.
		Model(&po.HzmUserDataPermission{}).
		Where("valid = 1 and executor_id in (?)", executorIds).
		Update("valid", false).
		Error
}
