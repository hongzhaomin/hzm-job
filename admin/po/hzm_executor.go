package po

// HzmExecutor 执行器表
type HzmExecutor struct {
	BasePo
	Name         *string // 执行器名称
	AppKey       *string // 执行器标识
	RegistryType *byte   // 注册方式：0-自动；1-手动
}

// ExecutorRegistryType 注册方式枚举
type ExecutorRegistryType byte

func (my ExecutorRegistryType) Is(registryType *byte) bool {
	return my == ExecutorRegistryType(*registryType)
}

func (my ExecutorRegistryType) ToPtr() *ExecutorRegistryType {
	p := new(ExecutorRegistryType)
	*p = my
	return p
}

const (
	AutoRegistry   ExecutorRegistryType = 0 // 自动注册
	ManualRegistry ExecutorRegistryType = 1 // 手动注册
)
