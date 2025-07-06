package internal

import (
	"context"
	"encoding/json"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/core/sdk"
)

const (
	heartBeat = "/api/heart-beat"
	jobHandle = "/api/job-handle"
)

func Post[T any](param sdk.Param[T]) (*T, error) {
	url := param.GetUrl()
	ctx := context.Background()
	client := global.SingletonPool().RemotingUtil
	jsonStr, err := client.PostJSON(ctx, url, param.GetAccessToken(), param)
	if err != nil {
		//slog.Error(fmt.Sprintf("请求http异常，url: %s", url), err.Error())
		return nil, err
	}

	var result sdk.Result[T]
	err = json.Unmarshal(jsonStr, &result)
	if result.Success {
		return result.Data, nil
	}
	return nil, err
}

// JobHandleReq 任务处理调度请求参数
type JobHandleReq struct {
	*sdk.BaseParam[sdk.Result[*bool]]
	JobId     *int64  `json:"jobId,omitempty"`     // 任务id
	JobName   *string `json:"jobName,omitempty"`   // 任务名称
	JobParams *string `json:"jobParams,omitempty"` // 任务参数
}
