package service

import (
	"github.com/hongzhaomin/hzm-job/admin/dao"
	"github.com/hongzhaomin/hzm-job/admin/internal"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/internal/global/iface"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"github.com/hongzhaomin/hzm-job/core/tools"
	"path"
)

const (
	heartBeat = "api/heart-beat"
	jobHandle = "api/job-handle"
)

var _ iface.CronFuncRegister = (*CronFuncRegister)(nil)

// CronFuncRegister 定时任务注册器
type CronFuncRegister struct {
	hzmJobDao          dao.HzmJobDao
	hzmExecutorDao     dao.HzmExecutorDao
	hzmExecutorNodeDao dao.HzmExecutorNodeDao
}

// RegistryHeatBeatFunc 注册心跳任务
func (my *CronFuncRegister) RegistryHeatBeatFunc() {
	// 注册心跳检测任务，每 5s 执行一次
	_, err := global.SingletonPool().Cron.AddFunc("0/5 * * * * ?", func() {
		executors, err := my.hzmExecutorDao.FindAll()
		if err != nil {
			// todo 日志：查询执行器错误
			return
		}

		executorIds := tools.GetIds4Slice(executors, func(executor *po.HzmExecutor) int64 {
			return *executor.Id
		})
		id2NodesMap := my.findExecutorNodesMap(nil, executorIds)

		accessToken := ""
		for _, executor := range executors {
			executorId := executor.Id
			nodes, ok := id2NodesMap[*executorId]
			if !ok {
				// todo 日志：id为%d的执行器节点不存在
				continue
			}

			for _, node := range nodes {
				go func() {
					// 当执行器是自动录入时，离线节点尝试再连接一次，连不上会被删除
					if po.AutoRegistry.Is(executor.RegistryType) && po.NodeOffline.Is(node.Status) {
						req := sdk.NewBaseParam[sdk.Result[*bool]](path.Join(*node.Address, heartBeat), accessToken)
						_, err = internal.Post[bool](req)
						if err != nil {
							// todo 日志：失败
							// 连不上，就删除掉
							err = my.hzmExecutorNodeDao.Delete(node.Id)
							if err != nil {
								// todo 日志：删除离线节点失败
							}
						}
						return
					}
					req := sdk.NewBaseParam[sdk.Result[*bool]](path.Join(*node.Address, heartBeat), accessToken)
					_, err = internal.Post[bool](req)
					if err != nil && po.NodeOnline.Is(node.Status) {
						// todo 日志：失败
						// 标记该节点为离线
						_ = my.hzmExecutorNodeDao.UpdateStatus(node.Id, po.NodeOffline)
					}
				}()
			}
		}
	})
	if err != nil {
		// todo 日志：心跳检测任务注册失败
	}
}

// RegistryJobs 注册所有配置的任务
func (my *CronFuncRegister) RegistryJobs() {
	jobs, err := my.hzmJobDao.FindRunningJobs()
	if err != nil {
		// todo 打印日志: select jobs error
		return
	}

	// 查询所有执行器
	//id2ExecutorMap := my.findExecutorMap(executorIds)
	//// 查询所有执行器对应的在线节点机器列表
	//id2NodesMap := my.findExecutorNodesMap(executorIds)

	for _, job := range jobs {
		spec := job.ScheduleValue
		if spec == nil {
			// todo 打印日志: job spec not exist
			continue
		}

		_, err := global.SingletonPool().Cron.AddFunc(*spec, func() {
			my.WrapperRegistryJobFunc(job, nil)
		})
		if err != nil {
			// todo 日志：注册任务失败：id: {}
		}
	}
}

// WrapperRegistryJobFunc 封装注册任务函数
func (my *CronFuncRegister) WrapperRegistryJobFunc(job *po.HzmJob, jobParameters *string) {
	executorId := job.ExecutorId
	if executorId == nil {
		return
	}

	// 查询执行器
	executor, err := my.hzmExecutorDao.FindById(*executorId)
	if err != nil {
		// todo 日志：id为%d的执行器不存在
		return
	}
	if executor == nil {
		// todo 日志：id为%d的执行器不存在
		return
	}

	// 查询执行器对应的在线节点机器列表
	nodes, err := my.hzmExecutorNodeDao.FindOnlineByExecutorIds(po.NodeOnline.ToPtr(), *executorId)
	if err != nil {
		// todo 日志：id为%d的执行器节点不存在
		return
	}
	if len(nodes) <= 0 {
		// todo 日志：id为%d的执行器无可用节点
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
			accessToken := ""
			if jobParameters == nil {
				jobParameters = job.Parameters
			}
			req := &internal.JobHandleReq{
				BaseParam: sdk.NewBaseParam[sdk.Result[*bool]](path.Join(*address, jobHandle), accessToken),
				JobId:     job.Id,
				JobName:   job.Name,
				JobParams: jobParameters,
			}
			_, err = internal.Post[bool](req)
			return err
		})
}

func (my *CronFuncRegister) findExecutorMap(executorIds []int64) map[int64]*po.HzmExecutor {
	if len(executorIds) <= 0 {
		return nil
	}

	executors, err := my.hzmExecutorDao.FindByIds(executorIds)
	if err != nil {
		// todo 日志：根据id查询执行器异常
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
		// todo 日志：根据执行器id查询执行器节点信息异常
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
