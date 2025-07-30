package consts

import (
	"errors"
	"time"
)

const (
	DefaultPort = "8888"

	// 客户端url
	HeartBeat = "/api/heart-beat"
	JobHandle = "/api/job-handle"
	JobCancel = "/api/job-cancel"
)

var (
	ServerError = errors.New("服务异常")

	JwtIssuer               = "hzm-job"
	JwtSecret               = []byte("your-secret-key")
	JwtTokenExpiresDuration = time.Hour * 24 * 7
)
