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
	"github.com/hongzhaomin/hzm-job/core/tools"
	"strings"
)

type HzmExecutorService struct {
	hzmExecutorDao           dao.HzmExecutorDao
	hzmExecutorNodeDao       dao.HzmExecutorNodeDao
	hzmUserDataPermissionDao dao.HzmUserDataPermissionDao
	hzmJobDao                dao.HzmJobDao
	hzmJobLogDao             dao.HzmJobLogDao
}

func (my *HzmExecutorService) PageExecutors(param req.ExecutorPage) (int64, []*vo.Executor) {
	count, executors, err := my.hzmExecutorDao.Page(param)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return 0, nil
	}
	voExecutors := tool.BeanConv[po.HzmExecutor, vo.Executor](executors,
		func(executor *po.HzmExecutor) (*vo.Executor, bool) {
			return &vo.Executor{
				Id:           executor.Id,
				Name:         executor.Name,
				AppKey:       executor.AppKey,
				RegistryType: executor.RegistryType,
			}, true
		})
	if len(voExecutors) <= 0 {
		return 0, nil
	}

	// 查询在线节点数量
	executorIds := tools.GetIds4Slice(executors, func(executor *po.HzmExecutor) int64 {
		return *executor.Id
	})
	executorId2CountMap, err := my.hzmExecutorNodeDao.CountOnlineByExecutorIds(executorIds)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
	} else {
		if len(executorId2CountMap) > 0 {
			for _, executor := range voExecutors {
				executor.OnlineNodeCount = executorId2CountMap[*executor.Id]
			}
		}
	}
	return count, voExecutors
}

func (my *HzmExecutorService) Add(param req.Executor) error {
	// 校验是否存在 AppKey 相同的执行器
	sameAppKeyExecutor, err := my.hzmExecutorDao.FindByAppKey(*param.AppKey)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	if sameAppKeyExecutor != nil {
		return errors.New(fmt.Sprintf("执行器[%s]已存在", *param.AppKey))
	}

	// 新增执行器
	executor := &po.HzmExecutor{
		Name:         param.Name,
		AppKey:       param.AppKey,
		RegistryType: param.RegistryType,
	}
	if err = my.hzmExecutorDao.Save(executor); err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}

	// 手动注册，新增执行器节点信息
	if po.ManualRegistry.Is(param.RegistryType) {
		nodes := tool.BeanConv4Basic[string, po.HzmExecutorNode](strings.Split(*param.Addresses, ","),
			func(adds string) (*po.HzmExecutorNode, bool) {
				adds = strings.TrimSpace(adds)
				if adds == "" {
					return nil, false
				}
				return &po.HzmExecutorNode{
					ExecutorId: executor.Id,
					Address:    &adds,
					Status:     (*byte)(po.NodeOffline.ToPtr()),
				}, true
			})
		if err = my.hzmExecutorNodeDao.SaveBatch(nodes...); err != nil {
			global.SingletonPool().Log.Error(err.Error())
			return consts.ServerError
		}
	}
	return nil
}

func (my *HzmExecutorService) Edit(param req.Executor) error {
	executor, err := my.hzmExecutorDao.FindById(*param.Id)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	if executor == nil {
		return errors.New("执行器不存在")
	}

	// 校验是否存在 AppKey 相同的执行器
	sameAppKeyExecutor, err := my.hzmExecutorDao.FindByAppKey(*param.AppKey)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	if sameAppKeyExecutor != nil && *sameAppKeyExecutor.Id != *executor.Id {
		return errors.New(fmt.Sprintf("执行器[%s]已存在", *param.AppKey))
	}

	// 更新执行器
	executor.AppKey = param.AppKey
	executor.Name = param.Name
	executor.RegistryType = param.RegistryType
	if err = my.hzmExecutorDao.Update(executor); err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}

	// 更新前：自动注册；更新后：手动注册；删除
	// 更新前：手动注册；更新后：手动注册；删除
	// 更新前：手动注册；更新后：自动注册；删除
	// 更新前：自动注册；更新后：自动注册；不删除
	if !(po.AutoRegistry.Is(param.RegistryType) && po.AutoRegistry.Is(executor.RegistryType)) {
		if err = my.hzmExecutorNodeDao.DeleteByExecutorId(executor.Id); err != nil {
			global.SingletonPool().Log.Error(err.Error())
			return consts.ServerError
		}
	}

	// 手动注册，新增执行器节点信息
	if po.ManualRegistry.Is(param.RegistryType) {
		nodes := tool.BeanConv4Basic[string, po.HzmExecutorNode](strings.Split(*param.Addresses, ","),
			func(adds string) (*po.HzmExecutorNode, bool) {
				adds = strings.TrimSpace(adds)
				if adds == "" {
					return nil, false
				}
				return &po.HzmExecutorNode{
					ExecutorId: executor.Id,
					Address:    &adds,
					Status:     (*byte)(po.NodeOffline.ToPtr()),
				}, true
			})
		if err = my.hzmExecutorNodeDao.SaveBatch(nodes...); err != nil {
			global.SingletonPool().Log.Error(err.Error())
			return consts.ServerError
		}
	}
	return nil
}

func (my *HzmExecutorService) QueryNodesByExecutorId(executorId int64) []*vo.ExecutorNode {
	nodes, err := my.hzmExecutorNodeDao.FindByExecutorIds(nil, executorId)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return nil
	}
	return tool.BeanConv[po.HzmExecutorNode, vo.ExecutorNode](nodes,
		func(node *po.HzmExecutorNode) (*vo.ExecutorNode, bool) {
			return &vo.ExecutorNode{
				Id:      node.Id,
				Address: node.Address,
				Status:  node.Status,
			}, true
		})
}

func (my *HzmExecutorService) LogicDeleteBatch(executorIds []int64) error {
	// 删除执行器（逻辑删）
	if err := my.hzmExecutorDao.LogicDeleteBatch(executorIds); err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	// 删除执行器节点数据（逻辑删）
	if err := my.hzmExecutorNodeDao.LogicDeleteBatchByExecutorIds(executorIds); err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	// 删除执行器相关权限数据（逻辑删）
	if err := my.hzmUserDataPermissionDao.LogicDeleteBatchByExecutorIds(executorIds); err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	// 删除执行器相关的job（逻辑删）
	if err := my.hzmJobDao.LogicDeleteBatchByExecutorIds(executorIds); err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	// 删除任务调度日志（逻辑删）
	if err := my.hzmJobLogDao.LogicDeleteBatchByExecutorIds(executorIds); err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	return nil
}

func (my *HzmExecutorService) QuerySelectBox() []*vo.ExecutorSelectBox {
	// fixme 鉴权功能弄好后，增加数据权限
	executors, err := my.hzmExecutorDao.FindAll()
	if len(executors) <= 0 {
		if err != nil {
			global.SingletonPool().Log.Error(err.Error())
		}
		return []*vo.ExecutorSelectBox{}
	}

	return tool.BeanConv[po.HzmExecutor, vo.ExecutorSelectBox](executors,
		func(executor *po.HzmExecutor) (*vo.ExecutorSelectBox, bool) {
			return &vo.ExecutorSelectBox{
				Name:  executor.Name,
				Value: executor.Id,
			}, true
		})
}
