package prop

import "github.com/hongzhaomin/hzm-job/core/ezconfig"

type HzmJobConfigBean struct {
	ezconfig.ConfigurationProperties `prefix:"hzm.job.client"`

	Localhost    string `default:"" notAutoRefresh:""`     // 配置本机ip地址（关闭自动刷新）
	Port         string `default:"8888" notAutoRefresh:""` // 配置客户端启动端口，默认8888（关闭自动刷新）
	AppKey       string `notAutoRefresh:""`                // 客户端服务名（关闭自动刷新）
	AdminAddress string // 配置服务端访问地址
}
