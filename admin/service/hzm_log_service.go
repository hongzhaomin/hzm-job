package service

import (
	"errors"
	"fmt"
	"github.com/hongzhaomin/hzm-job/admin/dao"
	"github.com/hongzhaomin/hzm-job/admin/internal"
	"github.com/hongzhaomin/hzm-job/admin/internal/consts"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/internal/tool"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/admin/po/cleartype"
	"github.com/hongzhaomin/hzm-job/admin/vo"
	"github.com/hongzhaomin/hzm-job/admin/vo/req"
	"github.com/hongzhaomin/hzm-job/core/tools"
	"time"
)

type HzmLogService struct {
	hzmJobLogDao             dao.HzmJobLogDao
	hzmExecutorDao           dao.HzmExecutorDao
	hzmExecutorNodeDao       dao.HzmExecutorNodeDao
	hzmJobDao                dao.HzmJobDao
	hzmUserDataPermissionDao dao.HzmUserDataPermissionDao
}

func (my *HzmLogService) PageLogs(loginUser *vo.User, param req.JobLogPage) (int64, []*vo.JobLog) {
	if po.CommonUser == *loginUser.Role {
		// 页面的执行器下拉框已经做了数据权限了，理论上选择了执行器，就无需做数据权限了
		if param.ExecutorId == nil {
			// 数据权限
			executorIds, err := my.hzmUserDataPermissionDao.FindExecutorIdsByUserId(*loginUser.Id)
			if len(executorIds) <= 0 {
				if err != nil {
					global.SingletonPool().Log.Error(err.Error())
				}
				return 0, nil
			}
			param.ExecutorIds = executorIds
		}
	}

	count, jobLogs, err := my.hzmJobLogDao.Page(param)
	if count == 0 {
		if err != nil {
			global.SingletonPool().Log.Error(err.Error())
		}
		return 0, nil
	}

	jobMap := my.findJobMap(tools.GetIds4DistinctSlice(jobLogs, func(log *po.HzmJobLog) *int64 {
		return log.JobId
	}))
	executorMap := my.findExecutorMap(tools.GetIds4DistinctSlice(jobLogs, func(log *po.HzmJobLog) *int64 {
		return log.ExecutorId
	}))

	voJobLogs := tool.BeanConv[po.HzmJobLog, vo.JobLog](jobLogs, func(jobLog *po.HzmJobLog) (*vo.JobLog, bool) {
		job := jobMap[*jobLog.JobId]
		executor := executorMap[*jobLog.ExecutorId]

		jobDescription := tool.GetOrDefault[*string](job, nil, func() *string {
			return job.Description
		})
		jobName := tool.GetOrDefault[*string](job, nil, func() *string {
			return job.Name
		})
		executorName := tool.GetOrDefault[*string](executor, nil, func() *string {
			return executor.Name
		})
		scheduleTime := tool.GetOrDefault[string](jobLog.ScheduleTime, "", func() string {
			return jobLog.ScheduleTime.Format(time.DateTime)
		})
		finishTime := tool.GetOrDefault[string](jobLog.FinishTime, "", func() string {
			return jobLog.FinishTime.Format(time.DateTime)
		})
		return &vo.JobLog{
			Id:                  jobLog.Id,
			JobId:               jobLog.JobId,
			JobDescription:      jobDescription,
			JobName:             jobName,
			ExecutorId:          jobLog.ExecutorId,
			ExecutorName:        executorName,
			ExecutorNodeAddress: jobLog.ExecutorNodeAddress,
			Parameters:          jobLog.Parameters,
			ScheduleTime:        &scheduleTime,
			Status:              jobLog.Status,
			HandleCode:          jobLog.HandleCode,
			HandleMsg:           jobLog.HandleMsg,
			FinishTime:          &finishTime,
		}, true
	})
	if len(voJobLogs) <= 0 {
		return 0, nil
	}

	return count, voJobLogs
}

func (my *HzmLogService) findJobMap(jobIds []int64) map[int64]*po.HzmJob {
	if len(jobIds) <= 0 {
		return nil
	}

	jobs, err := my.hzmJobDao.FindByIds(jobIds)
	if err != nil {
		global.SingletonPool().Log.Error("查询任务失败", "jobIds", jobIds, "err", err)
		return nil
	}

	id2JobMap := make(map[int64]*po.HzmJob, len(jobs))
	for _, job := range jobs {
		if job != nil {
			id2JobMap[*job.Id] = job
		}
	}
	return id2JobMap
}

func (my *HzmLogService) findExecutorMap(executorIds []int64) map[int64]*po.HzmExecutor {
	if len(executorIds) <= 0 {
		return nil
	}

	executors, err := my.hzmExecutorDao.FindByIds(executorIds)
	if err != nil {
		global.SingletonPool().Log.Error("查询执行器失败", "executorIds", executorIds, "err", err)
		return nil
	}

	id2ExecutorMap := make(map[int64]*po.HzmExecutor, len(executors))
	for _, executor := range executors {
		if executor != nil {
			id2ExecutorMap[*executor.Id] = executor
		}
	}
	return id2ExecutorMap
}

func (my *HzmLogService) DeleteByQuery(param req.LogDelParam) error {
	clearTyp := cleartype.ConvClearType(*param.ClearType)
	if clearTyp == nil {
		return errors.New("清理策略不合法")
	}
	timeBefore, countBefore := clearTyp.GetQueryParams()

	if err := my.hzmJobLogDao.DeleteByQuery(req.LogDelDaoParam{
		ExecutorId:       param.ExecutorId,
		JobId:            param.JobId,
		CreateTimeBefore: timeBefore,
		CountBefore:      countBefore,
	}); err != nil {
		global.SingletonPool().Log.Error("清理日志失败", "err", err)
		return consts.ServerError
	}
	return nil
}

func (my *HzmLogService) StopJob(param req.StopJobParam, userId int64) error {
	jobLog, err := my.hzmJobLogDao.FindById(*param.Id)
	if jobLog == nil {
		if err != nil {
			global.SingletonPool().Log.Error("查询任务调度日志失败", "err", err)
			return consts.ServerError
		}
		return errors.New("任务调度日志不存在")
	}

	executorId := jobLog.ExecutorId
	if executorId == nil {
		return errors.New("任务调度日志执行器不存在")
	}

	executor, err := my.hzmExecutorDao.FindById(*executorId)
	if executor == nil {
		if err != nil {
			global.SingletonPool().Log.Error("查询执行器失败", "err", err)
			return consts.ServerError
		}
		return errors.New("执行器不存在或已被删除")
	}

	err = internal.JobCancel2Client(*param.Address, param.Id, executor.AppSecret)
	if err != nil {
		global.SingletonPool().Log.Error("终止任务失败",
			"jobLogId", param.Id,
			"err", err)
		return errors.New("终止任务失败")
	}

	// 发送 手动终止任务 操作日志消息
	go func() {
		job, _ := my.hzmJobDao.FindById(*jobLog.JobId)
		if job == nil {
			return
		}
		desc := fmt.Sprintf("手动终止了任务[%s], 调度日志id: %d", *job.Name, *jobLog.Id)
		global.SingletonPool().MessageBus.SendMsg(&vo.OperateLogMsg{
			OperatorId:  userId,
			Description: desc,
			OperateTime: time.Now(),
			NewValue:    job,
		})
	}()
	return nil
}
