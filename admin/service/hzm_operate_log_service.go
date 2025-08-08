package service

import (
	"encoding/json"
	"fmt"
	"github.com/hongzhaomin/hzm-job/admin/dao"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/internal/tool"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/admin/vo"
	"github.com/hongzhaomin/hzm-job/core/tools"
	"time"
)

type HzmOperateLogService struct {
	hzmOperateLogDao dao.HzmOperateLogDao
	hzmUserDao       dao.HzmUserDao
}

// ReceiveMsg 消费调度统计消息
func (my *HzmOperateLogService) ReceiveMsg(msg *vo.OperateLogMsg) {
	global.SingletonPool().Log.Info("ReceiveMsg ==> 收到操作日志消息", "msg", *msg)
	opeLog := &po.HzmOperateLog{
		OperatorId:  &msg.OperatorId,
		Description: &msg.Description,
		OperateTime: &msg.OperateTime,
	}

	oldValue := msg.OldValue
	newValue := msg.NewValue
	if newValue != nil {
		var details []string
		switch nv := newValue.(type) {
		case *po.HzmJob:
			if oldValue != nil {
				ov := oldValue.(*po.HzmJob)
				details = append(details, "字段名:操作前:操作后")
				details = append(details, fmt.Sprintf("任务id:%d:%d", *ov.Id, *nv.Id))
				details = append(details, fmt.Sprintf("任务名称:%s:%s", *ov.Name, *nv.Name))
				details = append(details, fmt.Sprintf("调度类型:%s:%s",
					po.GetJobScheduleNameByType(ov.ScheduleType), po.GetJobScheduleNameByType(nv.ScheduleType)))
				details = append(details, fmt.Sprintf("调度值:%s:%s", *ov.ScheduleValue, *nv.ScheduleValue))
				details = append(details, fmt.Sprintf("任务参数:%s:%s", *ov.Parameters, *nv.Parameters))
				details = append(details, fmt.Sprintf("任务描述:%s:%s", *ov.Description, *nv.Description))
				details = append(details, fmt.Sprintf("负责人:%s:%s", *ov.Head, *nv.Head))
				details = append(details, fmt.Sprintf("任务状态:%s:%s",
					po.GetJobStatusNameByType(ov.Status), po.GetJobStatusNameByType(nv.Status)))
				details = append(details, fmt.Sprintf("路由策略:%s:%s",
					po.GetJobRouterStrategyNameByType(ov.RouterStrategy), po.GetJobRouterStrategyNameByType(nv.RouterStrategy)))
			} else {
				details = append(details, fmt.Sprintf("任务id:%d", *nv.Id))
				details = append(details, fmt.Sprintf("任务名称:%s", *nv.Name))
				details = append(details, fmt.Sprintf("任务描述:%s", *nv.Description))
				details = append(details, fmt.Sprintf("任务参数:%s", *nv.Parameters))
				details = append(details, fmt.Sprintf("执行器id:%d", *nv.ExecutorId))
			}
		case *po.HzmExecutor:
			if oldValue != nil {
				ov := oldValue.(*po.HzmExecutor)
				details = append(details, "字段名:操作前:操作后")
				details = append(details, fmt.Sprintf("执行器id:%d:%d", *ov.Id, *nv.Id))
				details = append(details, fmt.Sprintf("执行器名称:%s:%s", *ov.Name, *nv.Name))
				details = append(details, fmt.Sprintf("AppKey:%s:%s", *ov.AppKey, *nv.AppKey))
				details = append(details, fmt.Sprintf("注册方式:%s:%s",
					po.GetExeRegistryTypeNameByType(ov.RegistryType), po.GetExeRegistryTypeNameByType(nv.RegistryType)))
			} else {
				details = append(details, fmt.Sprintf("执行器id:%d", *nv.Id))
				details = append(details, fmt.Sprintf("执行器名称:%s", *nv.Name))
				details = append(details, fmt.Sprintf("AppKey:%s", *nv.AppKey))
				details = append(details, fmt.Sprintf("注册方式:%s", po.GetExeRegistryTypeNameByType(nv.RegistryType)))
			}
		default:
			global.SingletonPool().Log.Error("ReceiveMsg ==> 不支持的操作日志对象", "msg", *msg)
		}
		opeLog.Detail = my.convertDetail(details)
	}

	if err := my.hzmOperateLogDao.Insert(opeLog); err != nil {
		global.SingletonPool().Log.Info("ReceiveMsg ==> 操作日志添加失败", "msg", *msg, "err", err)
	}
}

func (my *HzmOperateLogService) convertDetail(details []string) *string {
	var jsonStr string
	if len(details) <= 0 {
		return &jsonStr
	}
	jsonBytes, err := json.Marshal(details)
	if err != nil {
		return &jsonStr
	}
	jsonStr = string(jsonBytes)
	return &jsonStr
}

func (my *HzmOperateLogService) OperateLogs() []*vo.OperateLog {
	list, err := my.hzmOperateLogDao.FindList()
	if err != nil {
		global.SingletonPool().Log.Info("操作日志查询失败", "err", err)
		return nil
	}

	userIds := tools.GetIds4DistinctSlice(list, func(opeLog *po.HzmOperateLog) *int64 {
		return opeLog.OperatorId
	})
	userId2NameMap := my.hzmUserDao.FindUserNameMap(userIds)

	return tool.BeanConv[po.HzmOperateLog, vo.OperateLog](list, func(opeLog *po.HzmOperateLog) (*vo.OperateLog, bool) {
		var details []string
		if opeLog.Detail != nil && *opeLog.Detail != "" {
			if e := json.Unmarshal([]byte(*opeLog.Detail), &details); e != nil {
				global.SingletonPool().Log.Error("操作日志详情json转换异常", "err", e)
			}
		}
		return &vo.OperateLog{
			Operator:    userId2NameMap[*opeLog.OperatorId],
			Description: *opeLog.Description,
			OperateTime: opeLog.OperateTime.Format(time.DateTime),
			Details:     details,
		}, true
	})
}
