package po

// HzmUserDataPermission 用户数据权限表(执行器维度)
type HzmUserDataPermission struct {
	BasePo
	UserId     *int64 // 用户id
	ExecutorId *int64 // 执行器id
}
