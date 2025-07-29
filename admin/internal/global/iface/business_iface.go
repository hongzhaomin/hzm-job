package iface

import "github.com/hongzhaomin/hzm-job/admin/po"

// CronFuncRegister 定时任务注册器接口
// 为了避免将单例注册到单例池中出现包引入循环依赖，因此在这里定义一个接口
type CronFuncRegister interface {

	// RegistryHeatBeatFunc 注册心跳任务
	RegistryHeatBeatFunc()

	// RegistryJobs 注册所有配置的任务
	RegistryJobs()

	// WrapperRegistryJobFunc 封装注册任务函数
	WrapperRegistryJobFunc(job *po.HzmJob, executorNodeId *int64)
}

// ExecutorNodeSelector 执行器节点选择器接口
type ExecutorNodeSelector interface {

	// StrategyType 策略类型
	StrategyType() po.JobRouterStrategy

	// NodeSchedule 节点任务调度
	NodeSchedule(nodes []*po.HzmExecutorNode, doSchedule func(*po.HzmExecutorNode) error)
}
