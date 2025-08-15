package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/internal/tool"
	"github.com/hongzhaomin/hzm-job/admin/po/cleartype"
	"github.com/hongzhaomin/hzm-job/admin/service"
	"github.com/hongzhaomin/hzm-job/admin/vo"
	"github.com/hongzhaomin/hzm-job/admin/vo/req"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"io"
	"net/http"
	"strings"
	"time"
)

// LogManage 任务调度日志管理
type LogManage struct {
	hzmLogService service.HzmLogService
}

// PageLogs 调度日志分页列表
// @Get /admin/log/page
func (my *LogManage) PageLogs(ctx *gin.Context) {
	var param req.JobLogPage
	if err := ctx.ShouldBind(&param); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	count, jobLogs := my.hzmLogService.PageLogs(tool.GetLoginUser(ctx), param)

	ctx.JSON(http.StatusOK, sdk.Ok4Page[vo.JobLog](count, jobLogs))
}

// QueryClearTypeSelectBox 日志清理策略下拉框
// @Get /admin/log/clear-type/select-box
func (my *LogManage) QueryClearTypeSelectBox(ctx *gin.Context) {
	selectBox := tool.BeanConv4Basic[cleartype.ClearType, vo.ClearTypeSelectBox](cleartype.Values(),
		func(typ cleartype.ClearType) (*vo.ClearTypeSelectBox, bool) {
			return &vo.ClearTypeSelectBox{
				Name:  typ.GetDesc(),
				Value: int(typ),
			}, true
		})
	ctx.JSON(http.StatusOK, sdk.Ok2[[]*vo.ClearTypeSelectBox](selectBox))
}

// DeleteByQuery 根据条件删除日志
// @Post /admin/log/del
func (my *LogManage) DeleteByQuery(ctx *gin.Context) {
	var param req.LogDelParam
	if err := ctx.ShouldBind(&param); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	if err := my.hzmLogService.DeleteByQuery(param); err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	// 发送 清理调度日志 操作日志消息
	go func() {
		desc := fmt.Sprintf("清理了%s调度日志",
			strings.ReplaceAll(cleartype.ConvClearType(*param.ClearType).GetDesc(), "清理", ""))
		global.SingletonPool().MessageBus.SendMsg(&vo.OperateLogMsg{
			OperatorId:  tool.GetUserId(ctx),
			Description: desc,
			OperateTime: time.Now(),
			NewValue:    &param,
		})
	}()

	ctx.JSON(http.StatusOK, sdk.Ok())
}

// StopJob 终止本次正在进行中的任务
// @Post /admin/log/stop-job
func (my *LogManage) StopJob(ctx *gin.Context) {
	var param req.StopJobParam
	if err := ctx.ShouldBind(&param); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	if err := my.hzmLogService.StopJob(param, tool.GetUserId(ctx)); err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, sdk.Ok())
}
