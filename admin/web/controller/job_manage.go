package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hongzhaomin/hzm-job/admin/internal/tool"
	"github.com/hongzhaomin/hzm-job/admin/po/sipcron"
	"github.com/hongzhaomin/hzm-job/admin/service"
	"github.com/hongzhaomin/hzm-job/admin/vo"
	"github.com/hongzhaomin/hzm-job/admin/vo/req"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"io"
	"net/http"
	"strconv"
)

// JobManage 任务管理
type JobManage struct {
	hzmJobService service.HzmJobService
}

// PageJobs 任务分页列表
// @Get /admin/job/page
func (my *JobManage) PageJobs(ctx *gin.Context) {
	var param req.JobPage
	if err := ctx.ShouldBind(&param); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	count, jobs := my.hzmJobService.PageJobs(param)

	ctx.JSON(http.StatusOK, sdk.Ok4Page[vo.Job](count, jobs))
}

// QuerySimpleCronSelectBox 极简表达式下拉框
// @Get /admin/job/simple-cron/select-box
func (my *JobManage) QuerySimpleCronSelectBox(ctx *gin.Context) {
	selectBox := tool.BeanConv4Basic[sipcron.SimpleCron, vo.SimpleCronSelectBox](sipcron.Values(),
		func(express sipcron.SimpleCron) (*vo.SimpleCronSelectBox, bool) {
			return &vo.SimpleCronSelectBox{
				Name:  express.GetDesc(),
				Value: string(express),
			}, true
		})
	ctx.JSON(http.StatusOK, sdk.Ok2[[]*vo.SimpleCronSelectBox](selectBox))
}

// JobSwitch 任务开关
// @Post /admin/job/switch
func (my *JobManage) JobSwitch(ctx *gin.Context) {
	jobIdStr := ctx.PostForm("jobId")
	if jobIdStr == "" {
		ctx.JSON(http.StatusOK, sdk.Fail("任务id不能为空"))
		return
	}
	jobId, err := strconv.Atoi(jobIdStr)
	if err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	if err = my.hzmJobService.JobSwitch(int64(jobId)); err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, sdk.Ok())
}

// Add 新增任务
// @Post /admin/job/add
func (my *JobManage) Add(ctx *gin.Context) {
	var param req.Job
	if err := ctx.ShouldBind(&param); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	if err := my.hzmJobService.Add(param); err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, sdk.Ok())
}

// Edit 编辑任务
// @Post /admin/job/edit
func (my *JobManage) Edit(ctx *gin.Context) {
	var param req.Job
	if err := ctx.ShouldBind(&param); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	if param.Id == nil || *param.Id <= 0 {
		ctx.JSON(http.StatusOK, sdk.Fail("任务id不能为空"))
		return
	}

	if err := my.hzmJobService.Edit(param); err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, sdk.Ok())
}

// DeleteBatch 删除任务
// @Post /admin/job/del
func (my *JobManage) DeleteBatch(ctx *gin.Context) {
	var jobIds []int64
	if err := ctx.ShouldBind(&jobIds); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	if err := my.hzmJobService.DeleteBatch(jobIds); err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, sdk.Ok())
}

// RunOnce 手动执行一次任务
// @Post /admin/job/run-once
func (my *JobManage) RunOnce(ctx *gin.Context) {
	var param req.JobRunOnce
	if err := ctx.ShouldBind(&param); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	if err := my.hzmJobService.RunOnce(param); err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, sdk.Ok())
}

// QuerySelectBox 任务下拉框
// @Get /admin/job/select-box
func (my *JobManage) QuerySelectBox(ctx *gin.Context) {
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

	jobSelectBox := my.hzmJobService.QuerySelectBox(int64(executorId))
	ctx.JSON(http.StatusOK, sdk.Ok2[[]*vo.JobSelectBox](jobSelectBox))
}
