package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hongzhaomin/hzm-job/admin/api"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/internal/global/iface"
	"github.com/hongzhaomin/hzm-job/admin/internal/prop"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/admin/service"
	"github.com/hongzhaomin/hzm-job/admin/web/controller/openapi"
	"github.com/hongzhaomin/hzm-job/core/config"
	"github.com/hongzhaomin/hzm-job/core/ezconfig"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"github.com/robfig/cron/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	filePath = flag.String("f", "hzm-job.yaml", "config file path")
)

func main() {
	log := global.SingletonPool().Log
	log.Debug("初始化注册任务", "file", *filePath)
	// 注册所有任务
	cronRegister := global.SingletonPool().CronFuncRegister
	cronRegister.RegistryHeatBeatFunc()
	cronRegister.RegistryJobs()

	// 启动 cron 定时任务
	c := global.SingletonPool().Cron
	c.Start()
	defer c.Stop()

	// 启动web服务
	// 创建开放接口controller
	openApi := openapi.NewJobServerOpenApi(&api.JobServerApiImpl{})
	router := NewGinRouter(openApi)
	router.Start()
}

func init() {
	// 解析flag命令行参数
	flag.Parse()

	// 初始化配置文件
	mysqlBean := new(prop.MysqlProperties)
	ezconfig.Builder().
		AddFiles(*filePath).
		AddConfigBeans(mysqlBean, new(prop.HzmJobConfigBean), new(config.LogBean), new(config.CommonConfigBean)).
		AddWatcher(configWatcher).
		Build()

	// 初始化mysql，gorm
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlBean.UserName, mysqlBean.Password, mysqlBean.Host, mysqlBean.Port, mysqlBean.Dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	// 创建日志对象
	log, logLevelVar := sdk.NewSlog()

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
		Log:              log,
		LogLevelVar:      logLevelVar,
		Mysql:            db,
		Cron:             c,
		RemotingUtil:     remotingUtil,
		CronFuncRegister: cronRegister,
		NodeSelectorMap:  nodeSelectorMap,
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
