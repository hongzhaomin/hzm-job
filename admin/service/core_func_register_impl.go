package service

import (
	"fmt"
	"github.com/hongzhaomin/hzm-job/admin/dao"
	"github.com/hongzhaomin/hzm-job/admin/internal"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/internal/global/iface"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"github.com/hongzhaomin/hzm-job/core/tools"
	"github.com/robfig/cron/v3"
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
	// fixme 加个本地缓存进行优化
	_, err := global.SingletonPool().Cron.AddFunc("0/5 * * * * ?", func() {
		defer func() {
			if e := recover(); e != nil {
				global.SingletonPool().Log.Error("心跳检测任务异常", "err", e)
			}
		}()
		executors, err := my.hzmExecutorDao.FindAll()
		if err != nil {
			global.SingletonPool().Log.Error("[心跳检测]查询执行器错误", "err", err)
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
				global.SingletonPool().Log.Warn("[心跳检测]执行器节点不存在", "executorId", *executorId)
				continue
			}

			for _, node := range nodes {
				go func() {
					defer func() {
						if e := recover(); e != nil {
							global.SingletonPool().Log.Error("[心跳检测]执行器节点心跳检测异常",
								"ExecutorId", node.ExecutorId,
								"Address", node.Address,
								"err", e)
						}
					}()
					// 当执行器是自动录入时，离线节点尝试再连接一次，连不上会被删除
					if po.AutoRegistry.Is(executor.RegistryType) && po.NodeOffline.Is(node.Status) {
						if err = internal.HeartBeat2Client(*node.Address); err != nil {
							global.SingletonPool().Log.Error("心跳检测失败", "nodeAddress", *node.Address, "err", err)
							// 连不上，就删除掉
							if err = my.hzmExecutorNodeDao.Delete(node.Id); err != nil {
								global.SingletonPool().Log.Error("[心跳检测]执行器节点删除失败",
									"nodeAddress", *node.Address, "err", err)
								return
							}
							// 将该节点存在执行中的任务调度日志状态结束掉，更改为执行失败
							runningJobLogIds, err2 := my.hzmJobLogDao.FindRunningLogIdsByAddress(node.Address)
							if err2 != nil {
								global.SingletonPool().Log.Error("[心跳检测]节点执行中的任务日志查询失败",
									"nodeAddress", *node.Address, "err", err)
								return
							}
							for _, logId := range runningJobLogIds {
								if err = my.jobLogErrorFinish(logId, "节点离线，异常结束"); err != nil {
									global.SingletonPool().Log.Error("[心跳检测]节点执行中的任务日志结束更新失败",
										"nodeAddress", *node.Address,
										"jobLogId", logId,
										"err", err)
								}
							}
						}
						return
					}
					if err = internal.HeartBeat2Client(*node.Address); err != nil {
						global.SingletonPool().Log.Error("心跳检测失败", "nodeAddress", *node.Address, "err", err)
						if po.NodeOnline.Is(node.Status) {
							// 标记该节点为离线
							if err = my.hzmExecutorNodeDao.UpdateStatus(node.Id, po.NodeOffline); err != nil {
								global.SingletonPool().Log.Error("[心跳检测]标记节点为离线失败",
									"nodeAddress", *node.Address, "err", err)
								return
							}
							// 将该节点存在执行中的任务调度日志状态结束掉，更改为执行失败
							runningJobLogIds, err2 := my.hzmJobLogDao.FindRunningLogIdsByAddress(node.Address)
							if err2 != nil {
								global.SingletonPool().Log.Error("[心跳检测]节点执行中的任务日志查询失败",
									"nodeAddress", *node.Address, "err", err)
								return
							}
							for _, logId := range runningJobLogIds {
								if err = my.jobLogErrorFinish(logId, "节点离线，异常结束"); err != nil {
									global.SingletonPool().Log.Error("[心跳检测]节点执行中的任务日志结束更新失败",
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
func (my *CronFuncRegister) jobLogErrorFinish(logId int64, handleMsg string) error {
	handleCode := 500
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
		if po.JobStop.Is(job.Status) {
			// 任务停止了，不注册
			continue
		}
		spec := job.ScheduleValue
		if spec == nil {
			global.SingletonPool().Log.Error("job spec not exist", "jobId", *job.Id, "err", err)
			continue
		}

		var entryId cron.EntryID
		entryId, err = global.SingletonPool().Cron.AddFunc(*spec, func() {
			my.WrapperRegistryJobFunc(job, nil)
		})
		if err != nil {
			global.SingletonPool().Log.Error("job registry failed", "jobId", *job.Id, "err", err)
		} else {
			// 将entryId更新到job中，方便后续删除注册的任务
			if err = my.hzmJobDao.UpdateCronEntryId(*job.Id, int(entryId)); err != nil {
				global.SingletonPool().Log.Error("update job cron entry id error",
					"jobId", *job.Id,
					"cronEntryId", entryId,
					"err", err)
			}
		}
	}
}

// WrapperRegistryJobFunc 封装注册任务函数
func (my *CronFuncRegister) WrapperRegistryJobFunc(job *po.HzmJob, executorNodeId *int64) {
	defer func() {
		if e := recover(); e != nil {
			global.SingletonPool().Log.Error("job执行异常",
				"JobId", job.Id,
				"JobName", job.Name,
				"err", e)
		}
	}()

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
			online := my.hzmExecutorNodeDao.IsOnline(jobLog.ExecutorNodeAddress)
			if !online {
				if err = my.jobLogErrorFinish(*jobLog.Id, "节点离线，异常结束"); err != nil {
					global.SingletonPool().Log.Error("执行中的任务日志结束更新失败",
						"executorNodeAddress", jobLog.ExecutorNodeAddress,
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
		global.SingletonPool().Log.Error("[任务调度]执行器查询失败", "executorId", *executorId, "err", err)
		return
	}
	if executor == nil {
		global.SingletonPool().Log.Error("[任务调度]执行器不存在", "executorId", *executorId)
		return
	}

	// 定义调度函数
	doSchedule := func(node *po.HzmExecutorNode) error {
		// 更新log状态为执行中
		runningJobLog := &po.HzmJobLog{
			BasePo: po.BasePo{
				Id: jobLog.Id,
			},
			ExecutorNodeAddress: node.Address,
			Parameters:          job.Parameters,
		}
		if err = my.hzmJobLogDao.UpdateLog4JobRunningById(runningJobLog); err != nil {
			return err
		}

		// 远程调度执行器任务
		success, err2 := internal.JobHandle2Client(func(url, accessToken string) *internal.JobHandleReq {
			return &internal.JobHandleReq{
				BaseParam: sdk.NewBaseParam[sdk.Result[*bool]](*node.Address+url, accessToken),
				LogId:     jobLog.Id,
				JobId:     job.Id,
				JobName:   job.Name,
				JobParams: job.Parameters,
			}
		})
		if success == nil || !*success || err2 != nil {
			// 请求失败，结束任务日志，调度失败
			if err = my.jobLogErrorFinish(*jobLog.Id, fmt.Sprintf("调度失败: %s", err2.Error())); err != nil {
				global.SingletonPool().Log.Error("结束任务日志失败[调度失败]",
					"executorNodeAddress", jobLog.ExecutorNodeAddress,
					"jobLogId", jobLog.Id,
					"err", err)
			}
		}
		return err2
	}

	// 传入的 executorNodeId 存在，则使用传入的节点地址调度
	// 否则，根据路由策略获取节点地址调度
	if executorNodeId != nil && *executorNodeId > 0 {
		node, findNodeErr := my.hzmExecutorNodeDao.FindById(executorNodeId)
		if node == nil {
			// 查询不到则node，跳转到路由调度
			errMsg := "数据不存在"
			if findNodeErr != nil {
				errMsg = findNodeErr.Error()
			}
			global.SingletonPool().Log.Error("[任务调度]查询指定执行器节点失败", "executorId", *executorId,
				"executorNodeId", *executorNodeId, "err", errMsg)
			goto routeDoSchedule
		}
		if err = doSchedule(node); err != nil {
			global.SingletonPool().Log.Error("任务调度失败", "executorId", *executorId,
				"nodeAddress", *node.Address, "err", err)
		}
		return
	}

routeDoSchedule:
	// 查询执行器对应的在线节点机器列表
	nodes, err := my.hzmExecutorNodeDao.FindByExecutorIds(po.NodeOnline.ToPtr(), *executorId)
	if len(nodes) <= 0 {
		msg := "执行器无可用节点"
		if err != nil {
			global.SingletonPool().Log.Error("[任务调度]执行器节点查询失败", "executorId", *executorId, "err", err)
			msg = err.Error()
		} else {
			global.SingletonPool().Log.Error("[任务调度]执行器无可用节点", "executorId", *executorId)
		}
		// 任务调度失败
		if err = my.jobLogErrorFinish(*jobLog.Id, fmt.Sprintf("调度失败: %s", msg)); err != nil {
			global.SingletonPool().Log.Error("结束任务日志失败[调度失败]",
				"executorNodeAddress", jobLog.ExecutorNodeAddress,
				"jobLogId", jobLog.Id,
				"err", err)
		}
		return
	}

	// 获取路由策略
	routerStrategy := po.JobPoll
	if job.RouterStrategy != nil {
		routerStrategy = po.JobRouterStrategy(*job.RouterStrategy)
	}
	// 根据路由策略执行调度
	global.SingletonPool().NodeSelectorMap[routerStrategy].NodeSchedule(nodes, doSchedule)
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

	nodes, err := my.hzmExecutorNodeDao.FindByExecutorIds(status, executorIds...)
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
