package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hongzhaomin/hzm-job/admin/service"
	"github.com/hongzhaomin/hzm-job/admin/vo"
	"github.com/hongzhaomin/hzm-job/admin/vo/req"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"io"
	"net/http"
)

// AuthController 鉴权控制器
type AuthController struct {
	hzmUserService service.HzmUserService
}

// Login 登录接口
// @Post /admin/login
func (my *AuthController) Login(ctx *gin.Context) {
	var param req.Login
	if err := ctx.ShouldBind(&param); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	loginUser, err := my.hzmUserService.Login(param)
	if err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, sdk.Ok2[vo.LoginUser](*loginUser))
}

// CheckLoginStatus 校验登录状态
// @Post /admin/check-login-status
func (my *AuthController) CheckLoginStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, sdk.Ok2[bool](true))
}

// LoginOut 退出登录
// @Post /admin/login-out
func (my *AuthController) LoginOut(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	if err := my.hzmUserService.LoginOut(userId.(int64)); err != nil {
		ctx.JSON(http.StatusOK, sdk.Fail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, sdk.Ok2[bool](true))
}
