package dao

import (
	"errors"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/admin/vo"
	"github.com/hongzhaomin/hzm-job/admin/vo/req"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type HzmJobDao struct{}

func (my *HzmJobDao) Save(job *po.HzmJob) error {
	err := global.SingletonPool().Mysql.
		Select("ExecutorId", "Name", "ScheduleType", "ScheduleValue",
			"Parameters", "Description", "Head", "RouterStrategy").
		Create(job).
		Error

	if err != nil {
		return err
	}

	global.SingletonPool().MessageBus.SendMsg(vo.SseDataBlock)

	return nil
}

func (my *HzmJobDao) FindRunningJobs() ([]*po.HzmJob, error) {
	var jobs []*po.HzmJob
	err := global.SingletonPool().Mysql.
		Where("valid = 1 and status = ?", po.JobRunning).
		Find(&jobs).
		Error
	return jobs, err
}

func (my *HzmJobDao) LogicDeleteBatchByExecutorIds(executorIds []int64) error {
	if len(executorIds) <= 0 {
		return nil
	}
	err := global.SingletonPool().Mysql.
		Model(&po.HzmJob{}).
		Where("valid = 1 and executor_id in (?)", executorIds).
		Update("valid", false).
		Error
	if err != nil {
		return err
	}

	var jobs []*po.HzmJob
	err = global.SingletonPool().Mysql.
		Select("status", "cron_entry_id").
		Where("valid = 0 and executor_id in(?)", executorIds).
		Find(&jobs).
		Error
	if err != nil {
		return err
	}
	if len(jobs) <= 0 {
		return nil
	}

	for _, job := range jobs {
		// 删除注册任务
		if job.CronEntryId != nil && po.JobRunning.Is(job.Status) {
			global.SingletonPool().Cron.Remove(cron.EntryID(*job.CronEntryId))
		}
	}

	global.SingletonPool().MessageBus.SendMsg(vo.SseDataBlock)
	return nil
}

func (my *HzmJobDao) Page(param req.JobPage) (int64, []*po.HzmJob, error) {
	// 构造条件
	db := global.SingletonPool().Mysql
	db = db.Where("valid = 1")
	if param.ExecutorId != nil {
		db = db.Where("executor_id = ?", param.ExecutorId)
	}
	if param.Status != nil {
		db = db.Where("status = ?", param.Status)
	}
	if param.Name != "" {
		db = db.Where("name = ?", param.Name)
	}
	if param.Description != "" {
		db = db.Where("description LIKE ?", "%"+param.Description+"%")
	}
	if param.Head != "" {
		db = db.Where("head = ?", param.Head)
	}
	if len(param.ExecutorIds) > 0 {
		db = db.Where("executor_id in(?)", param.ExecutorIds)
	}

	var count int64
	db.Model(po.HzmJob{}).Count(&count)
	if count == 0 {
		return 0, nil, nil
	}

	var jobs []*po.HzmJob
	err := db.Offset(param.Start()).Limit(param.Limit()).Find(&jobs).Error
	return count, jobs, err
}

func (my *HzmJobDao) FindById(id int64) (*po.HzmJob, error) {
	var job po.HzmJob
	err := global.SingletonPool().Mysql.
		Where("valid = 1 and id = ?", id).
		First(&job).
		Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &job, err
}

func (my *HzmJobDao) UpdateJobStatus(id int64, status po.JobStatus) error {
	err := global.SingletonPool().Mysql.
		Model(&po.HzmJob{}).
		Where("valid = 1 and id = ?", id).
		Update("status", status).
		Error

	if err != nil {
		return err
	}

	global.SingletonPool().MessageBus.SendMsg(vo.SseDataBlock)

	return nil
}

func (my *HzmJobDao) UpdateCronEntryId(id int64, cronEntryId int) error {
	return global.SingletonPool().Mysql.
		Model(&po.HzmJob{}).
		Where("valid = 1 and id = ?", id).
		Update("cron_entry_id", cronEntryId).
		Error
}

func (my *HzmJobDao) DeleteBatch(ids []int64) error {
	if len(ids) <= 0 {
		return nil
	}

	var jobs []*po.HzmJob
	err := global.SingletonPool().Mysql.
		Select("status", "cron_entry_id").
		Where("valid = 1 and id in(?)", ids).
		Find(&jobs).
		Error
	if err != nil {
		return err
	}
	if len(jobs) <= 0 {
		return nil
	}

	err = global.SingletonPool().Mysql.
		Unscoped().
		Where("valid = 1 and id in (?)", ids).
		Delete(&po.HzmJob{}).
		Error
	if err != nil {
		return err
	}

	for _, job := range jobs {
		// 删除注册任务
		if job.CronEntryId != nil && po.JobRunning.Is(job.Status) {
			global.SingletonPool().Cron.Remove(cron.EntryID(*job.CronEntryId))
		}
	}

	global.SingletonPool().MessageBus.SendMsg(vo.SseDataBlock)
	return nil
}

func (my *HzmJobDao) FindByExecutorIdAndName(executorId *int64, name *string) (*po.HzmJob, error) {
	if executorId == nil || name == nil || *name == "" {
		return nil, nil
	}
	var job po.HzmJob
	err := global.SingletonPool().Mysql.
		Where("valid = 1 and executor_id = ? and name = ?", *executorId, *name).
		First(&job).
		Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &job, err
}

func (my *HzmJobDao) Update(job *po.HzmJob) error {
	err := global.SingletonPool().Mysql.
		Save(job).
		Error

	if err != nil {
		return err
	}

	global.SingletonPool().MessageBus.SendMsg(vo.SseDataBlock)

	return nil
}

func (my *HzmJobDao) FindByExecutorId(executorId *int64) ([]*po.HzmJob, error) {
	if executorId == nil {
		return nil, nil
	}
	var jobs []*po.HzmJob
	err := global.SingletonPool().Mysql.
		// 只查询 id 和 description 两个字段
		Select("id", "description").
		Where("valid = 1 and executor_id = ?", *executorId).
		Find(&jobs).
		Error
	return jobs, err
}

func (my *HzmJobDao) FindByIds(ids []int64) ([]*po.HzmJob, error) {
	if ids == nil {
		return nil, nil
	}
	var jobs []*po.HzmJob
	err := global.SingletonPool().Mysql.
		Where("valid = 1 and id in(?)", ids).
		Find(&jobs).
		Error
	return jobs, err
}

func (my *HzmJobDao) CountStatistics() (int64, int64, error) {
	db := global.SingletonPool().Mysql
	var total int64
	err := db.Model(po.HzmJob{}).
		Where("valid = 1").
		Count(&total).
		Error

	var running int64
	err = db.Model(po.HzmJob{}).
		Where("valid = 1 and status = ?", po.JobRunning).
		Count(&running).
		Error
	return total, running, err
}
