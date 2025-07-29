package req

import (
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/admin/vo"
)

type UserPage struct {
	BasePage

	Role     string `json:"role" form:"role"`
	UserName string `json:"userName" form:"userName"`
}

type User struct {
	Id        *int64                  `json:"id,omitempty"`                          // 用户id
	UserName  *string                 `json:"userName,omitempty" binding:"required"` // 用户名
	Password  *string                 `json:"password,omitempty"`                    // 密码
	Role      *po.UserRole            `json:"role,omitempty" binding:"required"`     // 角色：0-管理员；1-普通用户
	Email     *string                 `json:"email,omitempty"`                       // 邮件
	DataPerms []*vo.DataPermsTransfer `json:"dataPerms,omitempty"`                   // 配置的权限数据
}

type EditPasswordParam struct {
	OldPassword   string `json:"oldPassword,omitempty" form:"oldPassword" binding:"required"`     // 旧密码
	NewPassword   string `json:"newPassword,omitempty" form:"newPassword" binding:"required"`     // 新密码
	AgainPassword string `json:"againPassword,omitempty" form:"againPassword" binding:"required"` // 确认新密码
}
