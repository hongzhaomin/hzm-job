package config

import "github.com/hongzhaomin/hzm-job/core/ezconfig"

type CommonConfigBean struct {
	ezconfig.ConfigurationProperties `prefix:"hzm.job.common"`

	AccessToken string // 鉴权token
}
