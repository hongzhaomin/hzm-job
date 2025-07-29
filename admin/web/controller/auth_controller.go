package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"net/http"
)

// AuthController 鉴权控制器
type AuthController struct {
}

// Login 登录接口
// @Post /admin/login
func (my *AuthController) Login(ctx *gin.Context) {
	// fixme 登录逻辑待完善
	ctx.JSON(http.StatusOK, sdk.Ok2[bool](true))
}

// CheckLoginStatus 校验登录状态
// @Post /admin/check-login-status
func (my *AuthController) CheckLoginStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, sdk.Ok2[bool](true))
}

// LoginOut 退出登录
// @Post /admin/login-out
func (my *AuthController) LoginOut(ctx *gin.Context) {
	// fixme 退出登录逻辑待完善
	ctx.JSON(http.StatusOK, sdk.Ok2[bool](true))
}
