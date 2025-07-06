package po

import "time"

type BasePo struct {
	Id         *int64     `gorm:"primaryKey"` // 主键
	Valid      *bool      // 是否可用
	CreateTime *time.Time // 创建时间
	UpdateTime *time.Time // 更新时间
}
