package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
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

// SseEvent 统计事件推送
// @Get home/sseEvent
func (my *HomeManage) SseEvent(ctx *gin.Context) {
	// 设置响应头为SSE
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*") // 如果需要跨域支持

	// 发送事件流，这里使用flusher来确保数据即时发送到客户端
	flusher := ctx.Writer.Flush

	finish := func() {
		ctx.SSEvent("message", vo.SseDone)
		flusher()
	}

	for {
		select {
		case sseMsg, ok := <-global.SingletonPool().MessageBus.GetSseMsgChan():
			if !ok {
				finish()
				return
			}
			ctx.SSEvent("message", sseMsg)
			flusher()
		case <-ctx.Done():
			finish()
			return
		}
	}
}
