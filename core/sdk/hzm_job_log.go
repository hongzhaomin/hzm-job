package sdk

import (
	"github.com/hongzhaomin/hzm-job/core/config"
	"github.com/hongzhaomin/hzm-job/core/ezconfig"
	"log/slog"
	"os"
	"strings"
)

func NewSlog() (*slog.Logger, *slog.LevelVar) {
	logConfig := ezconfig.Get[*config.LogBean]()
	var levelVar slog.LevelVar
	levelVar.Set(slog.LevelInfo)
	if logConfig.Level != "" {
		level := config.ConvLogLevel(logConfig.Level).ToSlogLevel()
		levelVar.Set(level)
	}
	opt := &slog.HandlerOptions{
		Level: &levelVar,
	}

	var handler slog.Handler
	logType := logConfig.Type
	if logType == "" {
		logType = string(config.TextLog)
	}
	switch config.LogType(strings.ToLower(logConfig.Type)) {
	case config.TextLog:
		handler = slog.NewTextHandler(os.Stdout, opt)
	case config.JsonLog:
		handler = slog.NewJSONHandler(os.Stdout, opt)
	default:
		panic("invalid log type")
	}

	return slog.New(handler), &levelVar
}
