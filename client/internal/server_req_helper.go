package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hongzhaomin/hzm-job/client/internal/consts"
	"github.com/hongzhaomin/hzm-job/client/internal/global"
	"github.com/hongzhaomin/hzm-job/client/internal/prop"
	"github.com/hongzhaomin/hzm-job/core/config"
	"github.com/hongzhaomin/hzm-job/core/ezconfig"
	"github.com/hongzhaomin/hzm-job/core/sdk"
)

// Registry2Admin 调用调度中心注册接口
func Registry2Admin() error {
	commonConfig := ezconfig.Get[*config.CommonConfigBean]()
	accessToken := commonConfig.AccessToken // token
	ip := GetHostIp()
	// 服务端地址
	clientConfig := ezconfig.Get[*prop.HzmJobConfigBean]()
	address := clientConfig.AdminAddress
	appKey := clientConfig.AppKey

	clientAddress := fmt.Sprintf("http://%s:%s", ip, ezconfig.Get[*prop.HzmJobConfigBean]().Port)
	req := &RegistryReq{
		BaseParam:       sdk.NewBaseParam[sdk.Result[*bool]](address+consts.Registry, accessToken),
		AppKey:          &appKey,
		ExecutorAddress: &clientAddress,
	}
	_, err := Post[bool](req)
	return err
}

// Offline2Admin 调用调度中心下线接口
func Offline2Admin() error {
	commonConfig := ezconfig.Get[*config.CommonConfigBean]()
	accessToken := commonConfig.AccessToken // token
	ip := GetHostIp()
	// 服务端地址
	clientConfig := ezconfig.Get[*prop.HzmJobConfigBean]()
	address := clientConfig.AdminAddress
	appKey := clientConfig.AppKey

	clientAddress := fmt.Sprintf("http://%s:%s", ip, ezconfig.Get[*prop.HzmJobConfigBean]().Port)
	req := &RegistryReq{
		BaseParam:       sdk.NewBaseParam[sdk.Result[*bool]](address+consts.Offline, accessToken),
		AppKey:          &appKey,
		ExecutorAddress: &clientAddress,
	}
	_, err := Post[bool](req)
	return err
}

// Callback2Admin 调用调度中心回调接口
func Callback2Admin(logId int64, code int, msg string) error {
	commonConfig := ezconfig.Get[*config.CommonConfigBean]()
	accessToken := commonConfig.AccessToken // token
	// 服务端地址
	clientConfig := ezconfig.Get[*prop.HzmJobConfigBean]()
	address := clientConfig.AdminAddress

	req := &JobResultReq{
		BaseParam:   sdk.NewBaseParam[sdk.Result[*bool]](address+consts.Callback, accessToken),
		LogId:       &logId,
		HandlerCode: &code,
		HandlerMsg:  &msg,
	}
	_, err := Post[bool](req)
	return err
}

func Post[T any](param sdk.Param[T]) (*T, error) {
	url := param.GetUrl()
	ctx := context.Background()
	client := global.SingletonPool().RemotingUtil
	jsonStr, err := client.PostJSON(ctx, url, param.GetAccessToken(), param)
	if err != nil {
		global.SingletonPool().Log.Error("请求http异常", "url", url, "err", err)
		return nil, err
	}

	var result sdk.Result[T]
	err = json.Unmarshal(jsonStr, &result)
	if result.Success {
		return result.Data, nil
	}
	return nil, err
}

// RegistryReq 执行器注册、下线参数
type RegistryReq struct {
	*sdk.BaseParam[sdk.Result[*bool]]
	AppKey          *string `json:"appKey,omitempty"`          // 执行器服务名称标识
	ExecutorAddress *string `json:"executorAddress,omitempty"` // 执行器地址（ip+端口）
}

// JobResultReq 任务处理回调参数
type JobResultReq struct {
	*sdk.BaseParam[sdk.Result[*bool]]
	//JobId       *int64  `json:"jobId,omitempty"`       // 任务id
	LogId       *int64  `json:"logId,omitempty"`       // 本次调度日志id
	HandlerCode *int    `json:"handlerCode,omitempty"` // 任务处理编码，200标识成功，其他失败
	HandlerMsg  *string `json:"handlerMsg,omitempty"`  // 任务处理结果消息
}
