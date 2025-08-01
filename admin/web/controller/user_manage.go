package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hongzhaomin/hzm-job/admin/internal/tool"
	"github.com/hongzhaomin/hzm-job/admin/service"
	"github.com/hongzhaomin/hzm-job/admin/vo"
	"github.com/hongzhaomin/hzm-job/admin/vo/req"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"io"
	"net/http"
	"strconv"
)

// UserManage 用户管理
type UserManage struct {
	hzmUserService service.HzmUserService
}

// Add 新增用户
// @Post /admin/user/add
func (my *UserManage) Add(ctx *gin.Context) {
	var param req.User
	if err := ctx.ShouldBind(&param); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}
	if param.Password == nil || *param.Password == "" {
		ctx.JSON(http.StatusOK, sdk.Fail("密码不能为空"))
		return
	}

	if err := my.hzmUserService.Add(param); err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, sdk.Ok())
}

// Edit 修改用户
// @Post /admin/user/edit
func (my *UserManage) Edit(ctx *gin.Context) {
	var param req.User
	if err := ctx.ShouldBind(&param); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}
	if param.Id == nil || *param.Id <= 0 {
		ctx.JSON(http.StatusOK, sdk.Fail("用户id不能为空"))
		return
	}

	if err := my.hzmUserService.Edit(param); err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, sdk.Ok())
}

func (my *UserManage) DeleteBatch(ctx *gin.Context) {
	var userIds []int64
	if err := ctx.ShouldBind(&userIds); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	if err := my.hzmUserService.DeleteBatch(userIds); err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, sdk.Ok())
}

// Current 获取当前登录用户
// @Get /admin/user/current
func (my *UserManage) Current(ctx *gin.Context) {
	user := my.hzmUserService.FindUserById(tool.GetUserId(ctx))
	ctx.JSON(http.StatusOK, sdk.Ok2[vo.User](*user))
}

// PageUsers 用户分页列表
// @Get /admin/user/page
func (my *UserManage) PageUsers(ctx *gin.Context) {
	var param req.UserPage
	if err := ctx.ShouldBind(&param); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	count, users := my.hzmUserService.PageUsers(param)

	ctx.JSON(http.StatusOK, sdk.Ok4Page[vo.User](count, users))
}

// DataPermsAll 查询所有可配置权限的数据
// @Get /admin/user/data-perms/all
func (my *UserManage) DataPermsAll(ctx *gin.Context) {
	dataPerms := my.hzmUserService.DataPermsAll()

	ctx.JSON(http.StatusOK, sdk.Ok2[[]*vo.DataPermsTransfer](dataPerms))
}

// QueryDataPermsByUserId 查询用户已配置权限的数据
// @Get /admin/user/data-perms/by-userid
func (my *UserManage) QueryDataPermsByUserId(ctx *gin.Context) {
	userIdStr := ctx.Query("userId")
	if userIdStr == "" {
		ctx.JSON(http.StatusOK, sdk.Fail("用户id不能为空"))
		return
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	dataPerms := my.hzmUserService.DataPermsByUserId(int64(userId))

	ctx.JSON(http.StatusOK, sdk.Ok2[vo.UserDataPerms](*dataPerms))
}

// EditPassword 当前登录用户修改密码
// @Post /admin/user/edit/password
func (my *UserManage) EditPassword(ctx *gin.Context) {
	var param req.EditPasswordParam
	if err := ctx.ShouldBind(&param); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	if err := my.hzmUserService.EditPassword(tool.GetUserId(ctx), param); err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, sdk.Ok())
}
