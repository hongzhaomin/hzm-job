package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hongzhaomin/hzm-job/admin/internal/consts"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/internal/prop"
	"github.com/hongzhaomin/hzm-job/admin/web/controller"
	"github.com/hongzhaomin/hzm-job/admin/web/controller/openapi"
	"github.com/hongzhaomin/hzm-job/core/ezconfig"
)

func NewGinRouter(openApi *openapi.JobServerOpenApi) *GinRouter {
	return &GinRouter{openApi: openApi}
}

type GinRouter struct {
	openApi        *openapi.JobServerOpenApi
	ToHtml         controller.ToHtml
	UserManage     controller.UserManage
	ExecutorManage controller.ExecutorManage
}

func (my *GinRouter) webGroup(webGroup *gin.RouterGroup) {
	webGroup.POST("/html/index", my.ToHtml.ToIndex)

	webGroup.POST("/user/page", my.UserManage.PageUsers)

	webGroup.POST("/executor/page", my.ExecutorManage.PageExecutors)
}

func (my *GinRouter) apiGroup(apiGroup *gin.RouterGroup) {
	// 客户端注册
	apiGroup.POST("/admin/registry", my.openApi.Registry)

	// 客户端下线
	apiGroup.POST("/admin/offline", my.openApi.Offline)

	// 任务完成回调
	apiGroup.POST("/admin/callback", my.openApi.Callback)
}

func (my *GinRouter) Start() {
	adminConfig := ezconfig.Get[*prop.HzmJobConfigBean]()

	//gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.LoadHTMLGlob("web/templates/*")
	engine.Static("/web/static", "./static")

	// 注册路由: openapi
	apiGroup := engine.Group("/api")
	my.apiGroup(apiGroup)

	// 注册路由: web
	webGroup := engine.Group("/admin")
	my.webGroup(webGroup)

	port := consts.DefaultPort
	if adminConfig.Port != "" {
		port = adminConfig.Port
	}
	global.SingletonPool().Log.Info(fmt.Sprintf("server started on port(s): %s (http) with context path '/'", port))
	_ = engine.Run(":" + port)
}
