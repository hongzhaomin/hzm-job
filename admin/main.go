package main

import (
	"fmt"
	"github.com/hongzhaomin/hzm-job/admin/internal/ezconfig"
	"github.com/hongzhaomin/hzm-job/admin/internal/properties"
	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New(cron.WithSeconds())

	c.AddFunc("* * * * * ?", func() {
		fmt.Println("hello world")
	})

	c.Start()
	defer c.Stop()
	select {}
}

func init() {
	mysqlBean := new(properties.MysqlProperties)
	// 初始化配置文件
	ezconfig.Builder().
		AddFiles("hzmjob.yaml").
		AddConfigBeans(mysqlBean).
		Build()

	// 初始化mysql，gorm
}
