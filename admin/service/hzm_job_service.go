package service

import (
	"errors"
	"fmt"
	"github.com/hongzhaomin/hzm-job/admin/dao"
	"github.com/hongzhaomin/hzm-job/admin/internal/consts"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/internal/tool"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/admin/vo"
	"github.com/hongzhaomin/hzm-job/admin/vo/req"
	"github.com/robfig/cron/v3"
	"time"
)

type HzmJobService struct {
	hzmJobDao                dao.HzmJobDao
	hzmJobLogDao             dao.HzmJobLogDao
	hzmUserDataPermissionDao dao.HzmUserDataPermissionDao
}

func (my *HzmJobService) PageJobs(loginUser *vo.User, param req.JobPage) (int64, []*vo.Job) {
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

	count, jobs, err := my.hzmJobDao.Page(param)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return 0, nil
	}

	voJobs := tool.BeanConv[po.HzmJob, vo.Job](jobs, func(job *po.HzmJob) (*vo.Job, bool) {
		return &vo.Job{
			Id:             job.Id,
			ExecutorId:     job.ExecutorId,
			Name:           job.Name,
			ScheduleType:   job.ScheduleType,
			ScheduleValue:  job.ScheduleValue,
			Parameters:     job.Parameters,
			Description:    job.Description,
			Head:           job.Head,
			Status:         job.Status,
			RouterStrategy: job.RouterStrategy,
		}, true
	})
	if len(voJobs) <= 0 {
		return 0, nil
	}

	return count, voJobs
}

func (my *HzmJobService) JobSwitch(jobId, userId int64) error {
	job, err := my.hzmJobDao.FindById(jobId)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	if job == nil {
		return errors.New("任务不存在")
	}

	if po.JobStop.Is(job.Status) {
		// 启动任务
		if err = my.hzmJobDao.UpdateJobStatus(*job.Id, po.JobRunning); err != nil {
			global.SingletonPool().Log.Error(err.Error())
			return consts.ServerError
		}
		// 注册cron
		var entryId cron.EntryID
		entryId, err = global.SingletonPool().Cron.AddFunc(*job.ScheduleValue, func() {
			global.SingletonPool().CronFuncRegister.WrapperRegistryJobFunc(job, nil)
		})
		if err != nil {
			global.SingletonPool().Log.Error("任务注册失败", "jobId", *job.Id, "err", err)
		} else {
			// 将entryId更新到job中，方便后续删除注册的任务
			if err = my.hzmJobDao.UpdateCronEntryId(*job.Id, int(entryId)); err != nil {
				global.SingletonPool().Log.Error("更新任务注册id失败",
					"jobId", *job.Id,
					"cronEntryId", entryId,
					"err", err)
			}
		}
	} else if po.JobRunning.Is(job.Status) {
		// 停止任务
		if err = my.hzmJobDao.UpdateJobStatus(*job.Id, po.JobStop); err != nil {
			global.SingletonPool().Log.Error(err.Error())
			return consts.ServerError
		}
		// 删除注册任务
		if job.CronEntryId != nil {
			global.SingletonPool().Cron.Remove(cron.EntryID(*job.CronEntryId))
		}
	}

	// 发送 开启/停止任务 操作日志消息
	go func() {
		var desc string
		if po.JobStop.Is(job.Status) {
			desc = fmt.Sprintf("启动了任务[%s]", *job.Name)
		} else if po.JobRunning.Is(job.Status) {
			desc = fmt.Sprintf("停止了任务[%s]", *job.Name)
		}
		global.SingletonPool().MessageBus.SendMsg(&vo.OperateLogMsg{
			OperatorId:  userId,
			Description: desc,
			OperateTime: time.Now(),
			NewValue:    job,
		})
	}()

	return nil
}

func (my *HzmJobService) Add(param req.Job, userId int64) error {
	sameNameJob, err := my.hzmJobDao.FindByExecutorIdAndName(param.ExecutorId, param.Name)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	if sameNameJob != nil {
		return errors.New(fmt.Sprintf("任务[%s]已存在", *param.Name))
	}

	// 新增任务
	newJob := &po.HzmJob{
		ExecutorId:     param.ExecutorId,
		Name:           param.Name,
		ScheduleType:   param.ScheduleType,
		ScheduleValue:  param.ScheduleValue,
		Parameters:     param.Parameters,
		Description:    param.Description,
		Head:           param.Head,
		RouterStrategy: param.RouterStrategy,
	}
	if err = my.hzmJobDao.Save(newJob); err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}

	// 发送 添加任务 操作日志消息
	go func() {
		desc := fmt.Sprintf("新增了任务[%s]", *newJob.Name)
		global.SingletonPool().MessageBus.SendMsg(&vo.OperateLogMsg{
			OperatorId:  userId,
			Description: desc,
			OperateTime: time.Now(),
			NewValue:    newJob,
		})
	}()
	return nil
}

func (my *HzmJobService) Edit(param req.Job, userId int64) error {
	job, err := my.hzmJobDao.FindById(*param.Id)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	if job == nil {
		return errors.New("任务不存在")
	}

	sameNameJob, err := my.hzmJobDao.FindByExecutorIdAndName(param.ExecutorId, param.Name)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	if sameNameJob != nil && *sameNameJob.Id != *param.Id {
		return errors.New(fmt.Sprintf("任务[%s]已存在", *param.Name))
	}

	// 编辑任务
	oldJob := *job
	job.ExecutorId = param.ExecutorId
	job.Name = param.Name
	job.ScheduleType = param.ScheduleType
	job.ScheduleValue = param.ScheduleValue
	job.Parameters = param.Parameters
	job.Description = param.Description
	job.Head = param.Head
	job.RouterStrategy = param.RouterStrategy
	if po.JobRunning.Is(job.Status) {
		if *job.ExecutorId != *oldJob.ExecutorId ||
			*job.Name != *oldJob.Name ||
			*job.ScheduleValue != *oldJob.ScheduleValue ||
			*job.Parameters != *oldJob.Parameters ||
			*job.RouterStrategy != *oldJob.RouterStrategy {
			// 如果修改了以上属性并且任务状态为启动中，需要删除cron，并重新注册
			if job.CronEntryId != nil {
				// 删除注册任务
				global.SingletonPool().Cron.Remove(cron.EntryID(*oldJob.CronEntryId))
				// 注册cron
				var entryId cron.EntryID
				entryId, err = global.SingletonPool().Cron.AddFunc(*job.ScheduleValue, func() {
					global.SingletonPool().CronFuncRegister.WrapperRegistryJobFunc(job, nil)
				})
				if err != nil {
					global.SingletonPool().Log.Error("任务注册失败", "jobId", *job.Id, "err", err)
				} else {
					// 将entryId更新到job中，方便后续删除注册的任务
					job.CronEntryId = (*int)(&entryId)
				}
			}
		}
	}
	if err = my.hzmJobDao.Update(job); err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}

	// 发送 添加任务 操作日志消息
	go func() {
		desc := fmt.Sprintf("修改了任务[%s]", *job.Name)
		global.SingletonPool().MessageBus.SendMsg(&vo.OperateLogMsg{
			OperatorId:  userId,
			Description: desc,
			OperateTime: time.Now(),
			OldValue:    &oldJob,
			NewValue:    job,
		})
	}()
	return nil
}

func (my *HzmJobService) DeleteBatch(jobIds []int64) error {
	if err := my.hzmJobDao.DeleteBatch(jobIds); err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}

	if err := my.hzmJobLogDao.DeleteByJobIds(jobIds); err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	return nil
}

func (my *HzmJobService) RunOnce(param req.JobRunOnce, userId int64) error {
	job, err := my.hzmJobDao.FindById(*param.Id)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	if job == nil {
		return errors.New("任务不存在")
	}

	job.Parameters = &param.Parameters
	global.SingletonPool().CronFuncRegister.WrapperRegistryJobFunc(job, param.ExecutorNodeId)

	// 发送 手动执行一次任务 操作日志消息
	go func() {
		desc := fmt.Sprintf("手动执行了一次任务[%s]", *job.Name)
		global.SingletonPool().MessageBus.SendMsg(&vo.OperateLogMsg{
			OperatorId:  userId,
			Description: desc,
			OperateTime: time.Now(),
			NewValue:    job,
		})
	}()
	return nil
}

func (my *HzmJobService) QuerySelectBox(executorId int64) []*vo.JobSelectBox {
	jobs, err := my.hzmJobDao.FindByExecutorId(&executorId)
	if len(jobs) <= 0 {
		if err != nil {
			global.SingletonPool().Log.Error(err.Error())
		}
		return []*vo.JobSelectBox{}
	}

	return tool.BeanConv[po.HzmJob, vo.JobSelectBox](jobs, func(job *po.HzmJob) (*vo.JobSelectBox, bool) {
		return &vo.JobSelectBox{
			Name:  job.Description,
			Value: job.Id,
		}, true
	})
}
