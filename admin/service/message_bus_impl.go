package service

import (
	"context"
	"encoding/json"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/internal/global/iface"
	"github.com/hongzhaomin/hzm-job/admin/vo"
)

var _ iface.MessageBus = (*MessageBusImpl)(nil)

func NewMessageBus() *MessageBusImpl {
	return &MessageBusImpl{
		scheduleStaMsgChan: make(chan *vo.ScheduleStaMsg, 20),
		opeLogMsgChan:      make(chan *vo.OperateLogMsg, 20),
	}
}

type MessageBusImpl struct {
	// 消息总线监听取消方法
	cancelMsgBus context.CancelFunc

	// 调度统计消息相关
	scheduleStaMsgChan chan *vo.ScheduleStaMsg
	hzmHomeService     HzmHomeService

	// 操作日志消息相关
	opeLogMsgChan        chan *vo.OperateLogMsg
	hzmOperateLogService HzmOperateLogService
}

func (my *MessageBusImpl) SendMsg(message any) {
	switch msg := message.(type) {
	case *vo.ScheduleStaMsg:
		global.SingletonPool().Log.Info("SendMsg ==> 发送调度统计消息", "msg", *msg)
		my.scheduleStaMsgChan <- msg
	case *vo.OperateLogMsg:
		global.SingletonPool().Log.Info("SendMsg ==> 发送操作日志消息", "msg", *msg)
		my.opeLogMsgChan <- msg
	default:
		msgStr, _ := json.Marshal(msg)
		global.SingletonPool().Log.Info("SendMsg ==> 不支持的消息类型", "msg", string(msgStr))
	}
}

func (my *MessageBusImpl) ListenEnable() {
	ctx, cancelMsgBus := context.WithCancel(context.Background())
	my.cancelMsgBus = cancelMsgBus
	for {
		select {
		case scheduleStaMsg := <-my.scheduleStaMsgChan:
			go my.hzmHomeService.ReceiveMsg(scheduleStaMsg)
		case opeLogMsg := <-my.opeLogMsgChan:
			go my.hzmOperateLogService.ReceiveMsg(opeLogMsg)
		case <-ctx.Done():
			return
		}
	}
}

func (my *MessageBusImpl) Stop() {
	if my.cancelMsgBus != nil {
		my.cancelMsgBus()
		close(my.scheduleStaMsgChan)
		close(my.opeLogMsgChan)
	}
}
