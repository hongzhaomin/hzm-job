package po

import "time"

// HzmOperateLog 操作日志表
type HzmOperateLog struct {
	BasePo
	OperatorId  *int64     // 操作人id
	Description *string    // 操作描述
	Detail      *string    // 操作内容详情：字符串数组json
	OperateTime *time.Time // 操作时间
}
