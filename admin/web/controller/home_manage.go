package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongzhaomin/hzm-job/admin/service"
	"github.com/hongzhaomin/hzm-job/admin/vo"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"net/http"
)

// HomeManage 首页
type HomeManage struct {
	hzmHomeService       service.HzmHomeService
	hzmOperateLogService service.HzmOperateLogService
}

// DataBlock 统计数据块
// @Get /admin/home/data-block
func (my *HomeManage) DataBlock(ctx *gin.Context) {
	dataBlock := my.hzmHomeService.DateBlock()
	ctx.JSON(http.StatusOK, sdk.Ok2[vo.DataBlock](*dataBlock))
}

// ScheduleTrend 调度走势
// @Get /admin/home/schedule-trend
func (my *HomeManage) ScheduleTrend(ctx *gin.Context) {
	trends := my.hzmHomeService.ScheduleTrend()
	ctx.JSON(http.StatusOK, sdk.Ok2[[]*vo.ScheduleTrend](trends))
}

// OperateLogs 操作日志记录
// @Get /admin/home/operate-logs
func (my *HomeManage) OperateLogs(ctx *gin.Context) {
	opeLogs := my.hzmOperateLogService.OperateLogs()
	ctx.JSON(http.StatusOK, sdk.Ok2[[]*vo.OperateLog](opeLogs))
}
