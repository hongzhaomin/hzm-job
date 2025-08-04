package hzmjob

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hongzhaomin/hzm-job/client/internal"
	"github.com/hongzhaomin/hzm-job/client/internal/api"
	"github.com/hongzhaomin/hzm-job/client/internal/consts"
	"github.com/hongzhaomin/hzm-job/client/internal/global"
	"github.com/hongzhaomin/hzm-job/client/internal/prop"
	"github.com/hongzhaomin/hzm-job/core/config"
	"github.com/hongzhaomin/hzm-job/core/ezconfig"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"os"
	"os/signal"
	"path"
	"time"
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

	// 启动内嵌 http 服务
	client := &api.HttpJobClient{}
	go client.Start()
	global.SingletonPool().Log.Info("embed http serve started")

	ctx, cancel := context.WithCancel(context.Background())
	// 延迟检测服务状态
	time.Sleep(100 * time.Millisecond)
	clientConfig := ezconfig.Get[*prop.HzmJobConfigBean]()
	accessToken := clientConfig.AppSecret // token
	url := fmt.Sprintf("http://localhost:%s", ezconfig.Get[*prop.HzmJobConfigBean]().Port) + path.Join(consts.BaseUrl, consts.HeartBeatUrl)
	ok, err := internal.Post[bool](sdk.NewBaseParam[sdk.Result[*bool]](url, accessToken))
	if err == nil && ok != nil && *ok {
		// 服务活跃，启动成功，向服务端发起注册请求
		err = internal.Registry2Admin()
		if err != nil {
			global.SingletonPool().Log.Error("hzm-job: 执行器自动注册失败", "err", err)
			// 注册失败，10s后进行无限重试
			go func() {
				ticker := time.NewTicker(10 * time.Second)
				defer ticker.Stop()
				for {
					select {
					case <-ticker.C:
						err = internal.Registry2Admin()
						if err == nil {
							break
						}
						// 注册失败，10s后进行无限重试
						global.SingletonPool().Log.Error("hzm-job: 执行器自动注册失败", "err", err)
					case <-ctx.Done():
						global.SingletonPool().Log.Error("hzm-job: 执行器自动注册失败重试退出")
						return
					}
				}
			}()
		}
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		cancel()
		client.Stop()
	}()
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
