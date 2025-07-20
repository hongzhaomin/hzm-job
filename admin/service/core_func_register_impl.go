package service

import (
	"github.com/hongzhaomin/hzm-job/admin/dao"
	"github.com/hongzhaomin/hzm-job/admin/internal"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/internal/global/iface"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"github.com/hongzhaomin/hzm-job/core/tools"
)

var _ iface.CronFuncRegister = (*CronFuncRegister)(nil)

// CronFuncRegister 定时任务注册器
type CronFuncRegister struct {
	hzmJobDao          dao.HzmJobDao
	hzmExecutorDao     dao.HzmExecutorDao
	hzmExecutorNodeDao dao.HzmExecutorNodeDao
	hzmJobLogDao       dao.HzmJobLogDao
}

// RegistryHeatBeatFunc 注册心跳任务
func (my *CronFuncRegister) RegistryHeatBeatFunc() {
	// 注册心跳检测任务，每 5s 执行一次
	_, err := global.SingletonPool().Cron.AddFunc("0/5 * * * * ?", func() {
		executors, err := my.hzmExecutorDao.FindAll()
		if err != nil {
			global.SingletonPool().Log.Error("查询执行器错误", "err", err)
			return
		}

		executorIds := tools.GetIds4Slice(executors, func(executor *po.HzmExecutor) int64 {
			return *executor.Id
		})
		id2NodesMap := my.findExecutorNodesMap(nil, executorIds)

		for _, executor := range executors {
			executorId := executor.Id
			nodes, ok := id2NodesMap[*executorId]
			if !ok {
				global.SingletonPool().Log.Warn("执行器节点不存在", "executorId", *executorId)
				continue
			}

			for _, node := range nodes {
				go func() {
					// 当执行器是自动录入时，离线节点尝试再连接一次，连不上会被删除
					if po.AutoRegistry.Is(executor.RegistryType) && po.NodeOffline.Is(node.Status) {
						if err = internal.HeartBeat2Client(*node.Address); err != nil {
							global.SingletonPool().Log.Error("心跳检测失败", "nodeAddress", *node.Address, "err", err)
							// 连不上，就删除掉
							err = my.hzmExecutorNodeDao.Delete(node.Id)
							if err != nil {
								global.SingletonPool().Log.Error("执行器节点删除失败", "nodeAddress", *node.Address, "err", err)
							}
						}
						return
					}
					if err = internal.HeartBeat2Client(*node.Address); err != nil {
						global.SingletonPool().Log.Error("心跳检测失败", "nodeAddress", *node.Address, "err", err)
						if po.NodeOnline.Is(node.Status) {
							// 标记该节点为离线
							_ = my.hzmExecutorNodeDao.UpdateStatus(node.Id, po.NodeOffline)
							// 将该节点存在执行中的任务调度日志状态结束掉，更改为执行失败
							runningJobLogIds, err2 := my.hzmJobLogDao.FindRunningLogIdsByNodeId(node.Id)
							if err2 != nil {
								global.SingletonPool().Log.Error("节点执行中的任务日志查询失败", "nodeAddress", *node.Address, "err", err)
								return
							}
							for _, logId := range runningJobLogIds {
								if err := my.jobLogErrorFinish(logId); err != nil {
									global.SingletonPool().Log.Error("节点执行中的任务日志结束更新失败",
										"nodeAddress", *node.Address,
										"jobLogId", logId,
										"err", err)
								}
							}
						}
					}
				}()
			}
		}
	})
	if err != nil {
		global.SingletonPool().Log.Error("心跳检测任务注册失败", "err", err)
	}
}

// 任务调度日志异常结束
func (my *CronFuncRegister) jobLogErrorFinish(logId int64) error {
	handleCode := 500
	handleMsg := "节点离线，异常结束"
	jobLog := &po.HzmJobLog{
		BasePo: po.BasePo{
			Id: &logId,
		},
		HandleCode: &handleCode,
		HandleMsg:  &handleMsg,
	}
	return my.hzmJobLogDao.FinishJobLogById(jobLog)
}

// RegistryJobs 注册所有配置的任务
func (my *CronFuncRegister) RegistryJobs() {
	jobs, err := my.hzmJobDao.FindRunningJobs()
	if err != nil {
		global.SingletonPool().Log.Error("select jobs error", "err", err)
		return
	}

	// 查询所有执行器
	//id2ExecutorMap := my.findExecutorMap(executorIds)
	//// 查询所有执行器对应的在线节点机器列表
	//id2NodesMap := my.findExecutorNodesMap(executorIds)

	for _, job := range jobs {
		spec := job.ScheduleValue
		if spec == nil {
			global.SingletonPool().Log.Error("job spec not exist", "jobId", *job.Id, "err", err)
			continue
		}

		_, err = global.SingletonPool().Cron.AddFunc(*spec, func() {
			my.WrapperRegistryJobFunc(job, nil)
		})
		if err != nil {
			global.SingletonPool().Log.Error("job registry failed", "jobId", *job.Id, "err", err)
		}
	}
}

// WrapperRegistryJobFunc 封装注册任务函数
func (my *CronFuncRegister) WrapperRegistryJobFunc(job *po.HzmJob, jobParameters *string) {
	executorId := job.ExecutorId
	if executorId == nil {
		return
	}

	// 查询是否存在未完成任务日志
	jobLog, err := my.hzmJobLogDao.FindUnfinishLogForUpdate(job.Id, executorId)
	if err != nil {
		global.SingletonPool().Log.Error("[任务调度]查询未完成任务日志异常", "jobId", *job.Id, "err", err)
	}
	if jobLog != nil && !po.LogToSchedule.Is(jobLog.Status) {
		if po.LogJobRunning.Is(jobLog.Status) {
			// 查询执行中的执行器节点是否在线，不在线则结束任务
			online := my.hzmExecutorNodeDao.IsOnline(jobLog.ExecutorNodeId)
			if !online {
				if err := my.jobLogErrorFinish(*jobLog.Id); err != nil {
					global.SingletonPool().Log.Error("执行中的任务日志结束更新失败",
						"executorNodeId", jobLog.ExecutorNodeId,
						"jobLogId", jobLog.Id,
						"err", err)
				}
				return
			}
		}
		global.SingletonPool().Log.Warn("[任务调度]存在未完成任务日志，本次调度忽略", "jobId", *job.Id)
		return
	}

	// 查询执行器
	executor, err := my.hzmExecutorDao.FindById(*executorId)
	if err != nil {
		global.SingletonPool().Log.Error("[任务调度]执行器不存在", "executorId", *executorId, "err", err)
		return
	}
	if executor == nil {
		global.SingletonPool().Log.Error("[任务调度]执行器不存在", "executorId", *executorId)
		return
	}

	// 查询执行器对应的在线节点机器列表
	nodes, err := my.hzmExecutorNodeDao.FindOnlineByExecutorIds(po.NodeOnline.ToPtr(), *executorId)
	if err != nil {
		global.SingletonPool().Log.Error("[任务调度]执行器节点不存在", "executorId", *executorId, "err", err)
		return
	}
	if len(nodes) <= 0 {
		global.SingletonPool().Log.Error("[任务调度]执行器无可用节点", "executorId", *executorId)
		return
	}

	// 获取路由策略
	routerStrategy := po.JobPoll
	if job.RouterStrategy != nil {
		routerStrategy = po.JobRouterStrategy(*job.RouterStrategy)
	}
	// 根据路由策略执行调度
	global.SingletonPool().NodeSelectorMap[routerStrategy].NodeSchedule(nodes,
		func(node *po.HzmExecutorNode) error {
			address := node.Address
			if jobParameters == nil {
				jobParameters = job.Parameters
			}

			// 更新log状态为执行中
			runningJobLog := &po.HzmJobLog{
				BasePo: po.BasePo{
					Id: jobLog.Id,
				},
				ExecutorNodeId: node.Id,
				Parameters:     jobParameters,
			}
			if err = my.hzmJobLogDao.UpdateLog4JobRunningById(runningJobLog); err != nil {
				return err
			}

			// 远程调度执行器任务
			success, err2 := internal.JobHandle2Client(func(url, accessToken string) *internal.JobHandleReq {
				return &internal.JobHandleReq{
					BaseParam: sdk.NewBaseParam[sdk.Result[*bool]](*address+url, accessToken),
					LogId:     jobLog.Id,
					JobId:     job.Id,
					JobName:   job.Name,
					JobParams: jobParameters,
				}
			})
			if success == nil || !*success || err2 != nil {
				// 请求失败，回滚任务日志状态为待调度
				err2 = my.hzmJobLogDao.RollbackToScheduleById(jobLog.Id)
			}
			return err2
		})
}

func (my *CronFuncRegister) findExecutorMap(executorIds []int64) map[int64]*po.HzmExecutor {
	if len(executorIds) <= 0 {
		return nil
	}

	executors, err := my.hzmExecutorDao.FindByIds(executorIds)
	if err != nil {
		global.SingletonPool().Log.Error("[任务调度]查询执行器失败", "executorIds", executorIds)
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

func (my *CronFuncRegister) findExecutorNodesMap(status *po.NodeStatus, executorIds []int64) map[int64][]*po.HzmExecutorNode {
	if len(executorIds) <= 0 {
		return nil
	}

	nodes, err := my.hzmExecutorNodeDao.FindOnlineByExecutorIds(status, executorIds...)
	if err != nil {
		global.SingletonPool().Log.Error("[任务调度]查询执行器节点信息失败", "executorIds", executorIds)
		return nil
	}

	id2NodesMap := make(map[int64][]*po.HzmExecutorNode, len(nodes))
	for _, node := range nodes {
		if node != nil {
			groupNodes, ok := id2NodesMap[*node.ExecutorId]
			if !ok {
				groupNodes = make([]*po.HzmExecutorNode, 0)
			}
			id2NodesMap[*node.ExecutorId] = append(groupNodes, node)
		}
	}
	return id2NodesMap
}
