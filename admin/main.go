package main

import (
	"flag"
	"fmt"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/internal/global/iface"
	"github.com/hongzhaomin/hzm-job/admin/internal/properties"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/admin/service"
	"github.com/hongzhaomin/hzm-job/core/ezconfig"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"github.com/robfig/cron/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	filePath = flag.String("f", "hzm-job.yaml", "config file path")
)

func main() {
	// 注册所有任务
	cronRegister := global.SingletonPool().CronFuncRegister
	cronRegister.RegistryHeatBeatFunc()
	cronRegister.RegistryJobs()

	// 启动 cron
	c := global.SingletonPool().Cron
	c.Start()
	defer c.Stop()

	// todo 启动web服务

	select {}
}

func init() {
	// 解析flag命令
	flag.Parse()

	// 初始化配置文件
	mysqlBean := new(properties.MysqlProperties)
	ezconfig.Builder().
		AddFiles(*filePath).
		AddConfigBeans(mysqlBean).
		Build()

	// 初始化mysql，gorm
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlBean.UserName, mysqlBean.Password, mysqlBean.Host, mysqlBean.Port, mysqlBean.Dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 创建 cron 定时任务对象
	c := cron.New(cron.WithSeconds())

	// 创建远程请求工具
	remotingUtil := sdk.NewRemotingUtil()

	// 创建定时任务注册器
	cronRegister := new(service.CronFuncRegister)

	// 创建执行器节点选择器
	pollSelector := service.NewPollExecutorNodeSelector()
	randSelector := new(service.RandomExecutorNodeSelector)
	errNextSelector := new(service.ErrNextExecutorNodeSelector)
	nodeSelectorMap := map[po.JobRouterStrategy]iface.ExecutorNodeSelector{
		pollSelector.StrategyType():    pollSelector,
		randSelector.StrategyType():    randSelector,
		errNextSelector.StrategyType(): errNextSelector,
	}

	// 存储全局对象池
	global.Store(&global.Obj{
		Mysql:            db,
		Cron:             c,
		RemotingUtil:     remotingUtil,
		CronFuncRegister: cronRegister,
		NodeSelectorMap:  nodeSelectorMap,
	})
}
