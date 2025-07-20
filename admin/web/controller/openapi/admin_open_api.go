package openapi

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hongzhaomin/hzm-job/admin/api"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"io"
	"net/http"
)

func NewJobServerOpenApi(jobServerApi api.JobServerApi) *JobServerOpenApi {
	return &JobServerOpenApi{jobServerApi: jobServerApi}
}

// JobServerOpenApi 调度中心开放api接口
type JobServerOpenApi struct {
	jobServerApi api.JobServerApi
}

// Registry 客户端注册接口
func (my *JobServerOpenApi) Registry(ctx *gin.Context) {
	var req api.RegistryReq
	if err := ctx.ShouldBind(&req); err != nil && !errors.Is(err, io.EOF) {
		global.SingletonPool().Log.Error("客户端注册失败", "err", err.Error())
		ctx.JSON(http.StatusOK, sdk.Fail[any](err.Error()))
		return
	}

	my.checkRegistryReq(&req)
	my.jobServerApi.Registry(&req)
	ctx.JSON(http.StatusOK, sdk.Ok2[bool](true))
}

// Offline 客户端下线接口
func (my *JobServerOpenApi) Offline(ctx *gin.Context) {
	var req api.RegistryReq
	if err := ctx.ShouldBind(&req); err != nil && !errors.Is(err, io.EOF) {
		global.SingletonPool().Log.Error("客户端下线失败", "err", err.Error())
		ctx.JSON(http.StatusOK, sdk.Fail[any](err.Error()))
		return
	}

	my.checkRegistryReq(&req)
	my.jobServerApi.Offline(&req)

	ctx.JSON(http.StatusOK, sdk.Ok2[bool](true))
}

// Callback 回调接口
func (my *JobServerOpenApi) Callback(ctx *gin.Context) {
	var req api.JobResultReq
	if err := ctx.ShouldBind(&req); err != nil && !errors.Is(err, io.EOF) {
		global.SingletonPool().Log.Error("任务完成回调失败", "err", err.Error())
		ctx.JSON(http.StatusOK, sdk.Fail[any](err.Error()))
		return
	}

	if req.LogId == nil {
		panic(errors.New("LogId is nil"))
	}
	if req.HandlerCode == nil {
		panic(errors.New("HandlerCode is nil"))
	}
	my.jobServerApi.Callback(&req)
	ctx.JSON(http.StatusOK, sdk.Ok2[bool](true))
}

func (my *JobServerOpenApi) checkRegistryReq(req *api.RegistryReq) {
	if req.AppKey == nil {
		panic(errors.New("appKey is nil"))
	}
	if req.ExecutorAddress == nil {
		panic(errors.New("executor address is nil"))
	}
}
