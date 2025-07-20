package config

import (
	"github.com/hongzhaomin/hzm-job/core/ezconfig"
	"log/slog"
	"strings"
)

type LogBean struct {
	ezconfig.ConfigurationProperties `prefix:"hzm.job.common.log"`

	Level string
	Type  string
}

// LogType 日志类型类型枚举
type LogType string

const (
	TextLog LogType = "text" // text格式
	JsonLog LogType = "json" // json格式
)

// LogLevel 日志级别枚举
type LogLevel string

func ConvLogLevel(level string) LogLevel {
	level = strings.ToLower(level)
	switch LogLevel(level) {
	case Debug:
		return Debug
	case Info:
		return Info
	case Warn:
		return Warn
	case Error:
		return Error
	default:
		panic("invalid log level")
	}
}

func (my LogLevel) ToSlogLevel() slog.Level {
	return logLevelMap[my]
}

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
)

var logLevelMap = map[LogLevel]slog.Level{
	Debug: slog.LevelDebug,
	Info:  slog.LevelInfo,
	Warn:  slog.LevelWarn,
	Error: slog.LevelError,
}
