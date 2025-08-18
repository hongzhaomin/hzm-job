package service

import (
	"github.com/hongzhaomin/hzm-job/admin/dao"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/internal/tool"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/admin/vo"
	"time"
)

type HzmHomeService struct {
	hzmJobDao                dao.HzmJobDao
	hzmExecutorDao           dao.HzmExecutorDao
	hzmScheduleStatisticsDao dao.HzmScheduleStatisticsDao
	hzmJobLogDao             dao.HzmJobLogDao
}

func (my *HzmHomeService) DateBlock() *vo.DataBlock {
	jobTotalNum, jobRunningNum, err := my.hzmJobDao.CountStatistics()
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return &vo.DataBlock{}
	}
	dataBlock := &vo.DataBlock{
		JobTotalNum:   jobTotalNum,
		RunningJobNum: jobRunningNum,
	}

	executorTotalNum, executorOfflineNum, err := my.hzmExecutorDao.CountStatistics()
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return dataBlock
	}
	dataBlock.ExecutorTotalNum = executorTotalNum
	dataBlock.ExecutorOfflineNum = executorOfflineNum
	return dataBlock
}

func (my *HzmHomeService) ScheduleTrend() []*vo.ScheduleTrend {
	days := 15
	endDay := time.Now()
	startDay := endDay.AddDate(0, 0, -days+1)
	scheduleStaList, err := my.hzmScheduleStatisticsDao.FindByGeStartDay(startDay.Format(time.DateOnly))
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return nil
	}

	day2StaMap := make(map[string]*po.HzmScheduleStatistics)
	for _, sta := range scheduleStaList {
		d := sta.Day.Format("01-02")
		day2StaMap[d] = sta
	}

	var staDays []string
	tool.TimeForEach(startDay, endDay, func(d time.Time) {
		staDays = append(staDays, d.Format("01-02"))
	})

	return tool.BeanConv4Basic[string, vo.ScheduleTrend](staDays,
		func(d string) (*vo.ScheduleTrend, bool) {
			sta, ok := day2StaMap[d]
			if !ok {
				return &vo.ScheduleTrend{StatisticsDate: d}, true
			}

			return &vo.ScheduleTrend{
				StatisticsDate: d,
				TotalNum:       *sta.TotalNum,
				SuccessNum:     *sta.SuccessNum,
				FailNum:        *sta.FailNum,
			}, true
		})
}

// SyncScheduleStatisticsJob 同步调度统计任务
func (my *HzmHomeService) SyncScheduleStatisticsJob(dateTime time.Time) {
	day := dateTime.Format(time.DateOnly)
	if dateTime.Hour() == 0 {
		// 零点，再跑一次昨天的数据
		day = dateTime.AddDate(0, 0, -1).Format(time.DateOnly)
	}
	total, success, fail, err := my.hzmJobLogDao.CountStatistics(day)
	if err != nil {
		global.SingletonPool().Log.Error("== SyncScheduleStatisticsJob ==> 同步调度统计任务失败", "err", err)
		return
	}

	if total <= 0 && success <= 0 && fail <= 0 {
		return
	}

	sta, err := my.hzmScheduleStatisticsDao.FindByDayIfAbsentCreate(day)
	if err != nil {
		global.SingletonPool().Log.Error("== SyncScheduleStatisticsJob ==> 同步调度统计任务失败", "err", err)
		return
	}
	if sta != nil {
		// 更新
		sta.TotalNum = &total
		sta.SuccessNum = &success
		sta.FailNum = &fail
		if err = my.hzmScheduleStatisticsDao.Update(sta); err != nil {
			global.SingletonPool().Log.Error("== SyncScheduleStatisticsJob ==> 同步调度统计任务失败", "err", err)
		}
	}

	global.SingletonPool().MessageBus.SendMsg(vo.SseScheduleTrend)
}

// ReceiveMsg 消费调度统计消息
func (my *HzmHomeService) ReceiveMsg(msg *vo.ScheduleStaMsg) {
	global.SingletonPool().Log.Info("ReceiveMsg ==> 收到调度统计消息", "msg", *msg)

	sta, err := my.hzmScheduleStatisticsDao.FindByDayIfAbsentCreate(msg.StaDay)
	if err != nil {
		global.SingletonPool().Log.Error("== ReceiveMsg ==> 查询调度统计任务失败", "err", err)
		return
	}
	if sta != nil {
		// 更新
		if err = my.hzmScheduleStatisticsDao.Increment(*sta.Id, msg.TotalIncr, msg.SuccessIncr, msg.FailIncr); err != nil {
			global.SingletonPool().Log.Error("== ReceiveMsg ==> 累计调度统计任务失败", "err", err)
		}
	}

	global.SingletonPool().MessageBus.SendMsg(vo.SseScheduleTrend)
}
