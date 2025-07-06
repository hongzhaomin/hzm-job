package properties

import "github.com/hongzhaomin/hzm-job/core/ezconfig"

type HzmJobConfigBean struct {
	ezconfig.ConfigurationProperties `prefix:"hzm.job.config"`

	AccessToken string // 鉴权token
}
