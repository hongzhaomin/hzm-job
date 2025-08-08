package hzmjob

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hongzhaomin/hzm-job/client/internal/global"
	"github.com/hongzhaomin/hzm-job/client/internal/prop"
	"github.com/hongzhaomin/hzm-job/core/config"
	"github.com/hongzhaomin/hzm-job/core/ezconfig"
	"github.com/hongzhaomin/hzm-job/core/sdk"
)

var (
	filePath = flag.String("f", "hzm-job.yaml", "config file path")
)

func init() {
	// 解析flag命令行参数
	flag.Parse()

	// 初始化配置文件
	ezconfig.Builder().
		AddFiles(*filePath).
		AddConfigBeans(new(config.LogBean), new(prop.HzmJobConfigBean)).
		AddWatcher(configWatcher).
		Build()

	// 创建日志对象
	log, logLevelVar := sdk.NewSlog()

	// 创建远程请求工具
	remotingUtil := sdk.NewRemotingUtil()

	// 存储全局对象池
	global.Store(&global.Obj{
		Log:          log,
		LogLevelVar:  logLevelVar,
		RemotingUtil: remotingUtil,
		JobCancelCtx: global.NewJobCancelCtx(),
	})
}

func configWatcher(params []ezconfig.WatcherParam, _ map[ezconfig.ConfigurationBean]ezconfig.ConfigurationBeanDefinition) {
	jsonParam, _ := json.Marshal(params)
	fmt.Println("配置发生改变：", string(jsonParam))
	for _, param := range params {
		key := param.Key
		val := param.NewVal
		switch key {
		case "hzm.job.log.level":
			if logLevel, ok := val.(string); ok && logLevel != "" {
				global.SingletonPool().LogLevelVar.Set(config.ConvLogLevel(logLevel).ToSlogLevel())
			}
		default:
			// nothing to do
		}
	}
}
