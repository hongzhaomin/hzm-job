package main

import (
	"github.com/hongzhaomin/hzm-job/client/internal/global"
	"github.com/hongzhaomin/hzm-job/core/ezconfig"
	"github.com/hongzhaomin/hzm-job/core/sdk"
)

func init() {
	// 初始化配置文件
	ezconfig.Builder().
		AddFiles("hzm-job.yaml").
		//AddConfigBeans(mysqlBean).
		Build()

	// 创建远程请求工具
	remotingUtil := sdk.NewRemotingUtil()

	// 存储全局对象池
	global.Store(&global.Obj{
		RemotingUtil: remotingUtil,
	})
}
