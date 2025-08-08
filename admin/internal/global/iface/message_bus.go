package iface

// MessageBus 全局消息总线
type MessageBus interface {

	// SendMsg 发送消息
	SendMsg(msg any)

	// ListenEnable 开启监听
	ListenEnable()

	// Stop 关闭全局消息总线
	Stop()
}
