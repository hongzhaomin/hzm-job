package prop

import (
	"errors"
	"github.com/hongzhaomin/hzm-job/core/ezconfig"
)

// LdapProperties ldap配置文件结构体
// 理论上是可以不启用 ldap 的
// 如果 addr 未配置，则表示不启用 ldap，配置则表示启用
type LdapProperties struct {
	ezconfig.ConfigurationProperties `prefix:"hzm.job.admin.ldap"`

	Addr  string `alias:"addr" default:""`    // ldap地址
	Dc    string `alias:"dc" default:""`      // 域组件，表示域名的一部分，通常用于构建目录树的根节点或顶级结构（必填）
	Group string `alias:"ou" default:""`      // 组织单元（可选）
	CnKey string `alias:"cnkey" default:"cn"` // 用户条目名称（cn/uid，可选，默认cn）
}

func (my *LdapProperties) Enabled() bool {
	return my.Addr != ""
}

func (my *LdapProperties) CheckProperties() {
	// 如果 addr 未配置，则表示不启用 ldap，配置则表示启用
	if my.Addr != "" {
		// 启用
		if my.Dc == "" {
			panic(errors.New("ldap dc is empty"))
		}
	}
}
