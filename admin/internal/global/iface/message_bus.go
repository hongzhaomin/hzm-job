package iface

import "github.com/hongzhaomin/hzm-job/admin/vo"

// MessageBus 全局消息总线
type MessageBus interface {

	// SendMsg 发送消息
	SendMsg(msg any)

	// ListenEnable 开启监听
	ListenEnable()

	// Stop 关闭全局消息总线
	Stop()

	// GetSseMsgChan 获取sse事件消息通道
	GetSseMsgChan() <-chan vo.SseMsg
}
