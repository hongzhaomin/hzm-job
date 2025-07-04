package properties

import "github.com/hongzhaomin/hzm-job/admin/internal/ezconfig"

// MysqlProperties mysql配置文件结构体
type MysqlProperties struct {
	ezconfig.ConfigurationProperties `prefix:"hzmJob.mysql"`

	Host     string // 域名
	Port     int    // 端口
	UserName string // 用户名
	Password string // 用户密码
	Dbname   string // 库名
}
