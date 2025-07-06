package po

// HzmExecutorNode 执行器节点表
type HzmExecutorNode struct {
	BasePo
	ExecutorId *int64  // 执行器id
	Address    *string // 节点地址
	Status     *byte   // 节点状态：0-离线；1-在线
}

// NodeStatus 节点状态枚举
type NodeStatus byte

func (my NodeStatus) Is(status *byte) bool {
	return my == NodeStatus(*status)
}

func (my NodeStatus) ToPtr() *NodeStatus {
	p := new(NodeStatus)
	*p = my
	return p
}

const (
	NodeOffline NodeStatus = 0 // 离线
	NodeOnline  NodeStatus = 1 // 在线
)
