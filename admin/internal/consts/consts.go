package consts

import "errors"

const (
	DefaultPort = "8888"

	// 客户端url
	HeartBeat = "/api/heart-beat"
	JobHandle = "/api/job-handle"
	JobCancel = "/api/job-cancel"
)

var (
	ServerError = errors.New("服务异常")
)
