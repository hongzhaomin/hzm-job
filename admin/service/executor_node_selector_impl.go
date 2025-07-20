package service

import (
	"encoding/json"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/internal/global/iface"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"math/rand"
	"sync"
	"time"
)

// ===================================================================================
// ===================================================================================
// =================================== 轮询选择器 =====================================
// ===================================================================================
// ===================================================================================
var _ iface.ExecutorNodeSelector = (*PollExecutorNodeSelector)(nil)

func NewPollExecutorNodeSelector() *PollExecutorNodeSelector {
	return &PollExecutorNodeSelector{
		lastSelectMap: make(map[int64]int64, 16),
		lock:          &sync.Mutex{},
	}
}

// PollExecutorNodeSelector 轮询选择器
type PollExecutorNodeSelector struct {
	// key: 执行器id
	// value: 上次选择的节点id
	lastSelectMap map[int64] /* executorId */ int64 /* nodeId */
	lock          *sync.Mutex
}

func (my *PollExecutorNodeSelector) StrategyType() po.JobRouterStrategy {
	return po.JobPoll
}

func (my *PollExecutorNodeSelector) nodeSelected(nodes []*po.HzmExecutorNode) *po.HzmExecutorNode {
	firstNode := nodes[0]
	executorId := *firstNode.ExecutorId
	firstNodeId := *firstNode.Id
	my.lock.Lock()
	defer my.lock.Unlock()
	lastSelectedNodeId, ok := my.lastSelectMap[executorId]
	if !ok {
		my.lastSelectMap[executorId] = firstNodeId
		return firstNode
	}
	// 前面已经保证了，nodes是按照id从小到大的顺序排过序的（从数据库查询 order by id asc 了）
	// 如果上次选择的是最后一个，从第一个开始
	if lastSelectedNodeId >= *nodes[len(nodes)-1].Id {
		my.lastSelectMap[executorId] = firstNodeId
		return firstNode
	}

	// 循环遍历，取出第一个比上次选择的id大的节点
	for _, node := range nodes {
		nodeId := *node.Id
		if nodeId > lastSelectedNodeId {
			my.lastSelectMap[executorId] = nodeId
			return node
		}
	}
	// 理论上不会走到这里
	// 但是如果以上都没找到，取第一个
	my.lastSelectMap[executorId] = firstNodeId
	return firstNode
}

func (my *PollExecutorNodeSelector) NodeSchedule(nodes []*po.HzmExecutorNode, doSchedule func(*po.HzmExecutorNode) error) {
	if len(nodes) <= 0 || doSchedule == nil {
		return
	}
	node := my.nodeSelected(nodes)
	err := doSchedule(node)
	if err != nil {
		global.SingletonPool().Log.Error("任务调度失败", "executorId", *node.ExecutorId,
			"nodeAddress", *node.Address, "err", err)
		return
	}
}

// ===================================================================================
// ===================================================================================
// =================================== 随机选择器 =====================================
// ===================================================================================
// ===================================================================================
var _ iface.ExecutorNodeSelector = (*RandomExecutorNodeSelector)(nil)

// RandomExecutorNodeSelector 随机选择器
type RandomExecutorNodeSelector struct{}

func (my *RandomExecutorNodeSelector) StrategyType() po.JobRouterStrategy {
	return po.JobRandom
}

func (my *RandomExecutorNodeSelector) NodeSchedule(nodes []*po.HzmExecutorNode, doSchedule func(*po.HzmExecutorNode) error) {
	if len(nodes) <= 0 || doSchedule == nil {
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	node := nodes[r.Intn(len(nodes))]
	err := doSchedule(node)
	if err != nil {
		global.SingletonPool().Log.Error("任务调度失败", "executorId", *node.ExecutorId,
			"nodeAddress", *node.Address, "err", err)
		return
	}
}

// ===================================================================================
// ===================================================================================
// =================================== 故障转移选择器 ==================================
// ===================================================================================
// ===================================================================================
var _ iface.ExecutorNodeSelector = (*ErrNextExecutorNodeSelector)(nil)

// ErrNextExecutorNodeSelector 故障转移选择器
type ErrNextExecutorNodeSelector struct{}

func (my *ErrNextExecutorNodeSelector) StrategyType() po.JobRouterStrategy {
	return po.JobErrNext
}

func (my *ErrNextExecutorNodeSelector) NodeSchedule(nodes []*po.HzmExecutorNode, doSchedule func(*po.HzmExecutorNode) error) {
	if len(nodes) <= 0 || doSchedule == nil {
		return
	}

	// 所谓故障转移，就是如果调度失败，要每一个都去重试一下，直到尝试到最后一个节点
	for _, node := range nodes {
		err := doSchedule(node)
		if err == nil {
			// 成功则结束
			return
		}
		global.SingletonPool().Log.Error("任务调度失败，故障转移，重试下一个节点", "executorId", *node.ExecutorId,
			"nodeAddress", *node.Address, "err", err)
	}
	marshal, _ := json.Marshal(nodes)
	global.SingletonPool().Log.Error("任务调度失败，故障转移，所有节点全部失败", "nodes", string(marshal))
}
