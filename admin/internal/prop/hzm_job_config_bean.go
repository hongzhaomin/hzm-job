package prop

import "github.com/hongzhaomin/hzm-job/core/ezconfig"

type HzmJobConfigBean struct {
	ezconfig.ConfigurationProperties `prefix:"hzm.job.admin"`

	Port string
}
