package vo

import "github.com/hongzhaomin/hzm-job/admin/po"

type User struct {
	Id       *int64       `json:"id,omitempty"`       // 用户id
	UserName *string      `json:"userName,omitempty"` // 用户名
	Role     *po.UserRole `json:"role,omitempty"`     // 角色：0-管理员；1-普通用户
	RoleName *string      `json:"roleName,omitempty"` // 角色名称
	Email    *string      `json:"email,omitempty"`    // 邮件
}

type UserDataPerms struct {
	AllDataPerms      []*DataPermsTransfer `json:"allDataPerms"`      // 所有权限数据列表
	SelectedDataPerms []*DataPermsTransfer `json:"selectedDataPerms"` // 选中的权限数据列表
}

type DataPermsTransfer struct {
	Value    int64  `json:"value"`    // 执行器id
	Title    string `json:"title"`    // 执行器名称+执行器标识
	Disabled bool   `json:"disabled"` // 是否禁用
	Checked  bool   `json:"checked"`  // 是否选中
}
