package dao

import (
	"errors"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"gorm.io/gorm"
	"time"
)

type HzmJobLogDao struct{}

func (my *HzmJobLogDao) Save(jobLog *po.HzmJobLog) error {
	return global.SingletonPool().Mysql.
		// fixme
		//Select("").
		Create(jobLog).
		Error
}

func (my *HzmJobLogDao) FindUnfinishLogForUpdate(jobId, executorId *int64) (*po.HzmJobLog, error) {
	if jobId == nil || executorId == nil {
		return nil, errors.New("jobId or executorId is nil")
	}
	var jobLog *po.HzmJobLog
	err := global.SingletonPool().Mysql.Transaction(func(tx *gorm.DB) error {
		var log po.HzmJobLog
		// 查询是否存在未完成任务日志
		err := tx.Raw("SELECT * FROM `hzm_job_log` WHERE valid = 1 and job_id = ? and status in ? ORDER BY `hzm_job_log`.`id` LIMIT 1 for update", jobId,
			[]po.LogStatus{po.LogToSchedule, po.LogJobRunning}).
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
	updates := map[string]any{
		"status":      po.LogJobFinished,
		"handle_code": jobLog.HandleCode,
		"finish_time": time.Now(),
	}
	if jobLog.HandleMsg != nil {
		updates["handle_msg"] = *jobLog.HandleMsg
	}
	return global.SingletonPool().Mysql.Model(jobLog).
		Where("valid = 1 and id = ? and status = ?", jobLog.Id, po.LogJobRunning).
		Updates(updates).
		Error
}

func (my *HzmJobLogDao) UpdateLog4JobRunningById(jobLog *po.HzmJobLog) error {
	if jobLog == nil || jobLog.Id == nil || jobLog.ExecutorNodeId == nil {
		return nil
	}
	updates := map[string]any{
		"status":           po.LogJobRunning,
		"executor_node_id": jobLog.ExecutorNodeId,
		"schedule_time":    time.Now(),
	}
	if jobLog.Parameters != nil {
		updates["parameters"] = jobLog.Parameters
	}
	return global.SingletonPool().Mysql.Model(&po.HzmJobLog{}).
		Where("valid = 1 and id = ? and status = ?", jobLog.Id, po.LogToSchedule).
		Updates(updates).
		Error
}

func (my *HzmJobLogDao) RollbackToScheduleById(id *int64) error {
	if id == nil {
		return nil
	}
	updates := map[string]any{
		"status":           po.LogToSchedule,
		"executor_node_id": "",
		"parameters":       "",
		"schedule_time":    nil,
	}
	return global.SingletonPool().Mysql.Model(&po.HzmJobLog{}).
		Where("valid = 1 and id = ? and status = ?", id, po.LogJobRunning).
		Updates(updates).
		Error
}

func (my *HzmJobLogDao) FindRunningLogIdsByNodeId(executorNodeId *int64) ([]int64, error) {
	if executorNodeId == nil {
		return nil, nil
	}
	var jobLogs []*po.HzmJobLog
	err := global.SingletonPool().Mysql.
		Select("id").
		Where("valid = 1 and executor_node_id = ? and status = ?", executorNodeId, po.LogJobRunning).
		Find(&jobLogs).
		Error

	var logIds []int64
	for _, jobLog := range jobLogs {
		logIds = append(logIds, *jobLog.Id)
	}
	return logIds, err
}
