package hzmjob

import (
	"context"
	"fmt"
	"github.com/hongzhaomin/hzm-job/client/internal"
	"github.com/hongzhaomin/hzm-job/client/internal/api"
	"github.com/hongzhaomin/hzm-job/client/internal/consts"
	"github.com/hongzhaomin/hzm-job/client/internal/global"
	"github.com/hongzhaomin/hzm-job/client/internal/prop"
	"github.com/hongzhaomin/hzm-job/core/ezconfig"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"path"
	"sync/atomic"
	"time"
)

var defaultJobStarter atomic.Pointer[Starter]

func init() {
	defaultJobStarter.Store(&Starter{})
}

func DefaultJobStarter() *Starter {
	return defaultJobStarter.Load()
}

func Enable() {
	DefaultJobStarter().Enable()
}

func Close() {
	DefaultJobStarter().Close()
}

type Starter struct {
	stared atomic.Bool
	client *api.HttpJobClient
	cancel context.CancelFunc
}

func (my *Starter) Enable() {
	swapped := my.stared.CompareAndSwap(false, true)
	if !swapped {
		global.SingletonPool().Log.Info("hzm-job is already started")
		return
	}
	// 启动内嵌 http 服务
	my.client = &api.HttpJobClient{}
	go my.client.Start()
	global.SingletonPool().Log.Info("embed http serve started")

	ctx, cancel := context.WithCancel(context.Background())
	my.cancel = cancel
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
	global.SingletonPool().Log.Info("hzm-job is enabled")
}

func (my *Starter) Close() {
	global.SingletonPool().JobCancelCtx.CancelAndRemoveAll()
	if my.cancel != nil {
		my.cancel()
	}

	if my.client != nil {
		my.client.Stop()
	}
	my.stared.Swap(false)
	global.SingletonPool().Log.Info("hzm-job closed")
}
