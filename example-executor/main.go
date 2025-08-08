package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hongzhaomin/hzm-job/client/annotation"
	"github.com/hongzhaomin/hzm-job/client/hzmjob"
	"os"
	"os/signal"
	"time"
)

func main() {
	hzmjob.AddJob("cancelableJobFuncTest", func(ctx context.Context, param *string) error {
		fmt.Println("====== cancelableJobFuncTest ========> 任务开始执行:", param)
		var count int
		for {
			select {
			case <-ctx.Done():
				fmt.Println("====== cancelableJobFuncTest ========> 任务取消:", param)
				return errors.New("任务[cancelableJobFuncTest]被调度中心终止")
			default:
				count++
				time.Sleep(time.Second * 5)
				fmt.Println(fmt.Sprintf("=== cancelableJobFuncTest ===> 模拟数据库操作，执行 %d 次", count))
			}
		}
	})

	hzmjob.AddJobs(&JobTest{})

	hzmjob.Enable()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	hzmjob.Close()
}

type JobTest struct {
	annotation.HzmJob[Req] `name:"commonJobFuncTest"`
}

func (my *JobTest) DoHandle(_ context.Context, param *Req) error {
	fmt.Println("====== commonJobFuncTest ========> 任务开始执行")
	paramJson, err := json.Marshal(param)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("====== commonJobFuncTest ========> 模拟任务执行:", string(paramJson))
	time.Sleep(time.Second * 3)

	fmt.Println("====== commonJobFuncTest ========> 任务执行结束:", string(paramJson))
	return nil
}

type Req struct {
	Name string `json:"name"`
}
