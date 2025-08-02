package internal

import (
	"context"
	"encoding/json"
	"github.com/hongzhaomin/hzm-job/admin/internal/consts"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/core/sdk"
)

// HeartBeat2Client 调用执行器心跳检测接口
func HeartBeat2Client(address string, appSecret *string) error {
	var accessToken string
	if appSecret != nil {
		accessToken = *appSecret
	}
	req := sdk.NewBaseParam[sdk.Result[*bool]](address+consts.HeartBeat, accessToken)
	_, err := Post[bool](req)
	return err
}

// JobHandle2Client 调用执行器任务执行接口
func JobHandle2Client(appSecret *string, node *po.HzmExecutorNode, jobLog *po.HzmJobLog, job *po.HzmJob) (*bool, error) {
	var accessToken string
	if appSecret != nil {
		accessToken = *appSecret
	}
	req := &JobHandleReq{
		BaseParam: sdk.NewBaseParam[sdk.Result[*bool]](*node.Address+consts.JobHandle, accessToken),
		LogId:     jobLog.Id,
		JobId:     job.Id,
		JobName:   job.Name,
		JobParams: job.Parameters,
	}
	return Post[bool](req)
}

// JobCancel2Client 调用执行器任务取消接口
func JobCancel2Client(address string, jobLogId *int64, appSecret *string) error {
	var accessToken string
	if appSecret != nil {
		accessToken = *appSecret
	}
	req := &JobCancelReq{
		BaseParam: sdk.NewBaseParam[sdk.Result[*bool]](address+consts.JobCancel, accessToken),
		LogId:     jobLogId,
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
		global.SingletonPool().Log.Error("请求http失败", "url", url, "err", err)
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
	LogId     *int64  `json:"logId,omitempty"`     // 本次调度日志id
	JobId     *int64  `json:"jobId,omitempty"`     // 任务id
	JobName   *string `json:"jobName,omitempty"`   // 任务名称
	JobParams *string `json:"jobParams,omitempty"` // 任务参数
}

// JobCancelReq 任务取消请求参数
type JobCancelReq struct {
	*sdk.BaseParam[sdk.Result[*bool]]
	LogId *int64 `json:"logId,omitempty"` // 需要取消的调度日志id
}
