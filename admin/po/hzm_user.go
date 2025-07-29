package po

// HzmUser 用户表
type HzmUser struct {
	BasePo
	UserName *string // 用户名
	Password *string // 密码
	Role     *byte   // 角色：0-管理员；1-普通用户
	Email    *string // 邮件
}

// UserRole 用户角色枚举
type UserRole byte

func GetNameByRole(role *byte) *string {
	if role == nil {
		return nil
	}

	var name string
	switch UserRole(*role) {
	case Admin:
		name = "管理员"
	case CommonUser:
		name = "普通用户"
	}
	return &name
}

func (my UserRole) Is(role *byte) bool {
	if role == nil {
		return false
	}
	return my == UserRole(*role)
}

const (
	Admin      UserRole = iota // 管理员
	CommonUser                 // 普通用户
)
