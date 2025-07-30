package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ToHtml 页面跳转
type ToHtml struct {
}

// ToLogin 跳转登录页面
// @Get /login
func (my *ToHtml) ToLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", nil)
}

// ToIndex 跳转index页面
// @Get /
func (my *ToHtml) ToIndex(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}

// ToHome 首页
// @Get /home
func (my *ToHtml) ToHome(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "home.html", nil)
}

// ToUser 用户管理页面
// @Get /user
func (my *ToHtml) ToUser(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "user.html", nil)
}

// ToUserAddLayer 新增用户弹框
// @Get /user-add-layer
func (my *ToHtml) ToUserAddLayer(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "user-add-layer.html", nil)
}

// ToUserEditLayer 编辑用户弹框
// @Get /user-edit-layer
func (my *ToHtml) ToUserEditLayer(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "user-edit-layer.html", nil)
}

// ToUserPassword 修改密码页面
// @Get /user-password
func (my *ToHtml) ToUserPassword(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "user-password.html", nil)
}

// ToExecutor 执行器管理页面
// @Get /executor
func (my *ToHtml) ToExecutor(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "executor.html", nil)
}

// ToExecutorAddLayer 执行器添加弹窗
// @Get /executor-add-layer
func (my *ToHtml) ToExecutorAddLayer(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "executor-add-layer.html", nil)
}

// ToExecutorEditLayer 执行器编辑弹窗
// @Get /executor-edit-layer
func (my *ToHtml) ToExecutorEditLayer(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "executor-edit-layer.html", nil)
}

// ToExecutorNodesInfoLayer 查看执行器节点信息弹窗
// @Get /executor-nodes-info-layer
func (my *ToHtml) ToExecutorNodesInfoLayer(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "executor-nodes-info-layer.html", nil)
}

// ToJob 任务管理页面
// @Get /job
func (my *ToHtml) ToJob(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "job.html", nil)
}

// ToJobAddLayer 任务添加弹窗
// @Get /job-add-layer
func (my *ToHtml) ToJobAddLayer(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "job-add-layer.html", nil)
}

// ToJobEditLayer 任务编辑弹窗
// @Get /job-edit-layer
func (my *ToHtml) ToJobEditLayer(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "job-edit-layer.html", nil)
}

//// ToJobRunOnceLayer 任务执行一次弹窗
//// @Get /job-run-once-layer
//func (my *ToHtml) ToJobRunOnceLayer(ctx *gin.Context) {
//	ctx.HTML(http.StatusOK, "job-run-once-layer.html", nil)
//}

// ToLog 调度日志管理页面
// @Get /log
func (my *ToHtml) ToLog(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "log.html", nil)
}

// ToLogClearLayer 调度日志清理弹窗
// @Get /log-clear-layer
func (my *ToHtml) ToLogClearLayer(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "log-clear-layer.html", nil)
}
