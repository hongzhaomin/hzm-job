package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hongzhaomin/hzm-job/client/internal/global"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"log/slog"
)

const (
	registry = "/api/admin/registry"
	offline  = "/api/admin/offline"
	callback = "/api/admin/callback"
)

func Post[T any](param sdk.Param[T]) *T {
	url := param.GetUrl()
	ctx := context.Background()
	client := global.Pool().RemotingUtil
	jsonStr, err := client.PostJSON(ctx, url, param.GetAccessToken(), param)
	if err != nil {
		slog.Error(fmt.Sprintf("请求http异常，url: %s", url), err.Error())
		return nil
	}

	var result sdk.Result[T]
	_ = json.Unmarshal(jsonStr, &result)
	if result.Success {
		return result.Data
	}
	return nil
}

// RegistryReq 执行器注册、下线参数
type RegistryReq struct {
	sdk.BaseParam[sdk.Result[*bool]]
	ExecutorAddress *string `json:"executorAddress,omitempty"` // 执行器地址（ip+端口）
	ExecutorName    *string `json:"executorName,omitempty"`    // 执行器服务名称
}

// JobResultReq 任务处理回调参数
type JobResultReq struct {
	sdk.BaseParam[sdk.Result[*bool]]
	JobId       *int64  `json:"jobId,omitempty"`       // 任务id
	HandlerCode *int    `json:"handlerCode,omitempty"` // 任务处理编码，200标识成功，其他失败
	HandlerMsg  *string `json:"handlerMsg,omitempty"`  // 任务处理结果消息
}
