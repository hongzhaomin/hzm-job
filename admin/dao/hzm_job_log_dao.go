package dao

import (
	"errors"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/admin/vo"
	"github.com/hongzhaomin/hzm-job/admin/vo/req"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type HzmJobLogDao struct{}

func (my *HzmJobLogDao) FindUnfinishLogForUpdate(jobId, executorId *int64) (*po.HzmJobLog, error) {
	if jobId == nil || executorId == nil {
		return nil, errors.New("jobId or executorId is nil")
	}
	var jobLog *po.HzmJobLog
	err := global.SingletonPool().Mysql.Transaction(func(tx *gorm.DB) error {
		var log po.HzmJobLog
		// 查询是否存在未完成任务日志
		//err := tx.Raw("SELECT * FROM `hzm_job_log` WHERE valid = 1 and job_id = ? and status in ? ORDER BY `hzm_job_log`.`id` LIMIT 1 for update", jobId,
		//	[]po.LogStatus{po.LogToSchedule, po.LogJobRunning}).
		//	First(&log).
		//	Error

		// SELECT * FROM `hzm_job_log` WHERE valid = 1 and job_id = 1 and status in (0,1) ORDER BY `hzm_job_log`.`id` LIMIT 1 FOR UPDATE
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("valid = 1 and job_id = ? and status in ?", jobId, []po.LogStatus{po.LogToSchedule, po.LogJobRunning}).
			First(&log).
			Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 数据不存在，则创建任务日志
				newJobLog := &po.HzmJobLog{
					JobId:      jobId,
					ExecutorId: executorId,
					Status:     (*byte)(po.LogToSchedule.ToPtr()),
				}
				if err = tx.Select("JobId", "ExecutorId", "Status").Create(newJobLog).Error; err != nil {
					return err
				}
				jobLog = newJobLog
				return nil
			}
			return err
		}
		jobLog = &log
		return nil
	})
	return jobLog, err
}

func (my *HzmJobLogDao) FinishJobLogById(jobLog *po.HzmJobLog) error {
	if jobLog == nil || jobLog.Id == nil || jobLog.HandleCode == nil {
		return nil
	}
	queryLog, e := my.FindById(*jobLog.Id)
	if e != nil {
		return e
	}

	updates := map[string]any{
		"status":      po.LogJobFinished,
		"handle_code": jobLog.HandleCode,
		"finish_time": time.Now(),
	}
	if jobLog.HandleMsg != nil {
		updates["handle_msg"] = *jobLog.HandleMsg
	}
	err := global.SingletonPool().Mysql.Model(jobLog).
		Where("valid = 1 and id = ?", jobLog.Id).
		Updates(updates).
		Error
	if err == nil {
		// 发消息，更新调度统计成功或失败数
		go func() {
			var total int64
			var success int64
			var fail int64
			if po.LogToSchedule.Is(queryLog.Status) {
				total = 1
			}
			if jobLog.HandleCode != nil && *jobLog.HandleCode == 200 {
				success = 1
			} else {
				fail = 1
			}
			global.SingletonPool().MessageBus.SendMsg(&vo.ScheduleStaMsg{
				StaDay:      queryLog.CreateTime.Format(time.DateOnly),
				TotalIncr:   total,
				SuccessIncr: success,
				FailIncr:    fail,
			})
		}()
	}

	return err
}

func (my *HzmJobLogDao) UpdateLog4JobRunningById(jobLog *po.HzmJobLog) error {
	if jobLog == nil || jobLog.Id == nil || jobLog.ExecutorNodeAddress == nil {
		return nil
	}
	updates := map[string]any{
		"status":                po.LogJobRunning,
		"executor_node_address": jobLog.ExecutorNodeAddress,
		"schedule_time":         time.Now(),
	}
	if jobLog.Parameters != nil {
		updates["parameters"] = jobLog.Parameters
	}
	err := global.SingletonPool().Mysql.Model(&po.HzmJobLog{}).
		Where("valid = 1 and id = ? and status = ?", jobLog.Id, po.LogToSchedule).
		Updates(updates).
		Error
	if err == nil {
		// 发消息，更新调度统计总数
		go func() {
			queryLog, e := my.FindById(*jobLog.Id)
			if e != nil {
				global.SingletonPool().Log.Error("Find job log by id error", e)
				return
			}
			global.SingletonPool().MessageBus.SendMsg(&vo.ScheduleStaMsg{
				StaDay:      queryLog.CreateTime.Format(time.DateOnly),
				TotalIncr:   1,
				SuccessIncr: 0,
				FailIncr:    0,
			})
		}()
	}

	return err
}

func (my *HzmJobLogDao) FindRunningLogIdsByAddress(address *string) ([]int64, error) {
	if address == nil {
		return nil, nil
	}
	var jobLogs []*po.HzmJobLog
	err := global.SingletonPool().Mysql.
		Select("id").
		Where("valid = 1 and executor_node_address = ? and status = ?", address, po.LogJobRunning).
		Find(&jobLogs).
		Error

	var logIds []int64
	for _, jobLog := range jobLogs {
		logIds = append(logIds, *jobLog.Id)
	}
	return logIds, err
}

func (my *HzmJobLogDao) LogicDeleteBatchByExecutorIds(executorIds []int64) error {
	if len(executorIds) <= 0 {
		return nil
	}
	return global.SingletonPool().Mysql.
		Model(&po.HzmJobLog{}).
		Where("valid = 1 and executor_id in (?)", executorIds).
		Update("valid", false).
		Error
}

func (my *HzmJobLogDao) Page(param req.JobLogPage) (int64, []*po.HzmJobLog, error) {
	// 构造条件
	db := global.SingletonPool().Mysql
	db = db.Where("valid = ?", 1)
	if param.ExecutorId != nil {
		db = db.Where("executor_id = ?", param.ExecutorId)
	}
	if param.JobId != nil {
		db = db.Where("job_id = ?", param.JobId)
	}
	if param.Status != nil {
		db = db.Where("status = ?", param.Status)
	} else {
		db = db.Where("status in (?)", []po.LogStatus{po.LogJobRunning, po.LogJobFinished})
	}
	if param.ScheduleStartTime != "" {
		db = db.Where("schedule_time >= ?", param.ScheduleStartTime)
	}
	if param.ScheduleEndTime != "" {
		db = db.Where("schedule_time <= ?", param.ScheduleEndTime)
	}
	if len(param.ExecutorIds) > 0 {
		db = db.Where("executor_id in(?)", param.ExecutorIds)
	}

	var count int64
	db.Model(po.HzmJobLog{}).Count(&count)
	if count == 0 {
		return 0, nil, nil
	}

	var logs []*po.HzmJobLog
	err := db.Order("id desc").Offset(param.Start()).Limit(param.Limit()).Find(&logs).Error
	return count, logs, err
}

func (my *HzmJobLogDao) DeleteByJobIds(jobIds []int64) error {
	if len(jobIds) <= 0 {
		return nil
	}
	return global.SingletonPool().Mysql.
		Unscoped().
		Where("valid = 1 and job_id in (?)", jobIds).
		Delete(&po.HzmJobLog{}).
		Error
}

func (my *HzmJobLogDao) DeleteByQuery(param req.LogDelDaoParam) error {
	executorId := param.ExecutorId
	jobId := param.JobId
	createTimeBefore := param.CreateTimeBefore
	countBefore := param.CountBefore

	dbBase := global.SingletonPool().Mysql.Where("valid = 1")
	if executorId != nil {
		dbBase = dbBase.Where("executor_id = ?", executorId)
	}
	if jobId != nil {
		dbBase = dbBase.Where("job_id = ?", jobId)
	}
	if createTimeBefore != nil {
		dbBase = dbBase.Where("create_time < ?", createTimeBefore)
	}

	// 该条件放在最后，需要共用db，根据上面的查询条件查询 countBefore 条数据之前的id
	// 创建隔离副本，复用上面构造的条件
	if countBefore != nil && *countBefore > 0 {
		var log po.HzmJobLog
		err := dbBase.Session(&gorm.Session{}).
			Select("id").
			Offset(*countBefore - 1).
			Last(&log).
			Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if log.Id == nil {
			// 未查询到数据，表示数据不需要删除
			return nil
		}
		dbBase = dbBase.Where("id < ?", log.Id)
	}
	return dbBase.Unscoped().Delete(&po.HzmJobLog{}).Error
}

func (my *HzmJobLogDao) FindById(id int64) (*po.HzmJobLog, error) {
	var log po.HzmJobLog
	err := global.SingletonPool().Mysql.
		Where("valid = 1 and id = ?", id).
		First(&log).
		Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &log, err
}

func (my *HzmJobLogDao) CountStatistics(day string) (int64, int64, int64, error) {
	if day == "" {
		return 0, 0, 0, nil
	}
	startTime := day + " 00:00:00"
	endTime := day + " 23:59:59"

	db := global.SingletonPool().Mysql
	var total int64
	err := db.Model(po.HzmJobLog{}).
		Where("valid = 1 and status in(?) and create_time >= ? and create_time <= ?",
			[]po.LogStatus{po.LogJobRunning, po.LogJobFinished}, startTime, endTime).
		Count(&total).
		Error

	var success int64
	err = db.Model(po.HzmJobLog{}).
		Where("valid = 1 and handle_code = 200 and create_time >= ? and create_time <= ?", startTime, endTime).
		Count(&success).
		Error

	var fail int64
	err = db.Model(po.HzmJobLog{}).
		Where("valid = 1 and handle_code != 200 and create_time >= ? and create_time <= ?", startTime, endTime).
		Count(&fail).
		Error
	return total, success, fail, err
}
