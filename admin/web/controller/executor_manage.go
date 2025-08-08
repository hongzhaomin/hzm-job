package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/internal/tool"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/admin/service"
	"github.com/hongzhaomin/hzm-job/admin/vo"
	"github.com/hongzhaomin/hzm-job/admin/vo/req"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ExecutorManage 执行器管理
type ExecutorManage struct {
	hzmExecutorService service.HzmExecutorService
}

// PageExecutors 执行器分页列表
// @Get /admin/executor/page
func (my *ExecutorManage) PageExecutors(ctx *gin.Context) {
	var param req.ExecutorPage
	if err := ctx.ShouldBind(&param); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	count, executors := my.hzmExecutorService.PageExecutors(tool.GetLoginUser(ctx), param)

	ctx.JSON(http.StatusOK, sdk.Ok4Page[vo.Executor](count, executors))
}

// GenerateSecret 如需鉴权，生成appSecret
// @Post /admin/executor/generate-secret
func (my *ExecutorManage) GenerateSecret(ctx *gin.Context) {
	secret := tool.RandStr(32)
	ctx.JSON(http.StatusOK, sdk.Ok2(secret))
}

// Add 新增执行器
// @Post /admin/executor/add
func (my *ExecutorManage) Add(ctx *gin.Context) {
	var param req.Executor
	if err := ctx.ShouldBind(&param); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	if po.ManualRegistry.Is(param.RegistryType) {
		if param.Addresses == nil {
			ctx.JSON(http.StatusOK, sdk.Fail("手动注册节点地址不能为空"))
			return
		}
		trimSpaceAddress := strings.TrimSpace(*param.Addresses)
		param.Addresses = &trimSpaceAddress
		if trimSpaceAddress == "" {
			ctx.JSON(http.StatusOK, sdk.Fail("手动注册节点地址不能为空"))
			return
		}
	}

	if err := my.hzmExecutorService.Add(param, tool.GetUserId(ctx)); err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, sdk.Ok())
}

// Edit 编辑执行器
// @Post /admin/executor/edit
func (my *ExecutorManage) Edit(ctx *gin.Context) {
	var param req.Executor
	if err := ctx.ShouldBind(&param); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	if param.Id == nil || *param.Id <= 0 {
		ctx.JSON(http.StatusOK, sdk.Fail("执行器id不能为空"))
		return
	}

	if po.ManualRegistry.Is(param.RegistryType) {
		if param.Addresses == nil {
			ctx.JSON(http.StatusOK, sdk.Fail("手动注册节点地址不能为空"))
			return
		}
		trimSpaceAddress := strings.TrimSpace(*param.Addresses)
		param.Addresses = &trimSpaceAddress
		if trimSpaceAddress == "" {
			ctx.JSON(http.StatusOK, sdk.Fail("手动注册节点地址不能为空"))
			return
		}
	}

	if err := my.hzmExecutorService.Edit(param, tool.GetUserId(ctx)); err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, sdk.Ok())
}

// QueryNodesByExecutorId 根据执行器id查询节点信息
// @Get /admin/executor/nodes/by-executorid
func (my *ExecutorManage) QueryNodesByExecutorId(ctx *gin.Context) {
	executorIdStr := ctx.Query("executorId")
	if executorIdStr == "" {
		ctx.JSON(http.StatusOK, sdk.Fail("执行器id不能为空"))
		return
	}
	executorId, err := strconv.Atoi(executorIdStr)
	if err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	executorNodes := my.hzmExecutorService.QueryNodesByExecutorId(int64(executorId))

	ctx.JSON(http.StatusOK, sdk.Ok2[[]*vo.ExecutorNode](executorNodes))
}

// DeleteBatch 删除执行器，逻辑删除
// @Post /admin/executor/del
func (my *ExecutorManage) DeleteBatch(ctx *gin.Context) {
	var executorIds []int64
	if err := ctx.ShouldBind(&executorIds); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	if err := my.hzmExecutorService.LogicDeleteBatch(executorIds); err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	// 发送 删除执行器 操作日志消息
	go func() {
		jsonByte, _ := json.Marshal(executorIds)
		desc := fmt.Sprintf("删除了执行器(执行器id:%s)", string(jsonByte))
		global.SingletonPool().MessageBus.SendMsg(&vo.OperateLogMsg{
			OperatorId:  tool.GetUserId(ctx),
			Description: desc,
			OperateTime: time.Now(),
		})
	}()

	ctx.JSON(http.StatusOK, sdk.Ok())
}

// QuerySelectBox 执行器下拉框查询
// @Get /admin/executor/select-box
func (my *ExecutorManage) QuerySelectBox(ctx *gin.Context) {
	executorSelectBox := my.hzmExecutorService.QuerySelectBox(tool.GetLoginUser(ctx))
	ctx.JSON(http.StatusOK, sdk.Ok2[[]*vo.ExecutorSelectBox](executorSelectBox))
}
