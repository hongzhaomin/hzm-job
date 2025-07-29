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
	toHtml         controller.ToHtml
	authController controller.AuthController
	jobManage      controller.JobManage
	logManage      controller.LogManage
	executorManage controller.ExecutorManage
	userManage     controller.UserManage
	menuManage     controller.MenuManage
}

func (my *GinRouter) webGroup(webGroup *gin.RouterGroup) {
	// 鉴权相关
	webGroup.POST("/login", my.authController.Login)
	webGroup.POST("/check-login-status", my.authController.CheckLoginStatus)
	webGroup.POST("/login-out", my.authController.LoginOut)

	// 菜单初始化
	webGroup.GET("/menu/init", my.menuManage.Init)
	webGroup.GET("/menu/get-log-menu", my.menuManage.GetLogMenu)

	// 任务管理
	webGroup.GET("/job/page", my.jobManage.PageJobs)
	webGroup.GET("/job/simple-cron/select-box", my.jobManage.QuerySimpleCronSelectBox)
	webGroup.POST("/job/switch", my.jobManage.JobSwitch)
	webGroup.POST("/job/add", my.jobManage.Add)
	webGroup.POST("/job/edit", my.jobManage.Edit)
	webGroup.POST("/job/del", my.jobManage.DeleteBatch)
	webGroup.POST("/job/run-once", my.jobManage.RunOnce)
	webGroup.GET("/job/select-box", my.jobManage.QuerySelectBox)

	// 任务调度日志
	webGroup.GET("/log/page", my.logManage.PageLogs)
	webGroup.GET("/log/clear-type/select-box", my.logManage.QueryClearTypeSelectBox)
	webGroup.POST("/log/del", my.logManage.DeleteByQuery)
	webGroup.POST("/log/stop-job", my.logManage.StopJob)

	// 执行器管理
	webGroup.GET("/executor/page", my.executorManage.PageExecutors)
	webGroup.POST("/executor/add", my.executorManage.Add)
	webGroup.POST("/executor/edit", my.executorManage.Edit)
	webGroup.POST("/executor/del", my.executorManage.DeleteBatch)
	webGroup.GET("/executor/nodes/by-executorid", my.executorManage.QueryNodesByExecutorId)
	webGroup.GET("/executor/select-box", my.executorManage.QuerySelectBox)

	// 用户管理
	webGroup.POST("/user/add", my.userManage.Add)
	webGroup.POST("/user/edit", my.userManage.Edit)
	webGroup.POST("/user/del", my.userManage.DeleteBatch)
	webGroup.POST("/user/edit/password", my.userManage.EditPassword)
	webGroup.GET("/user/current", my.userManage.Current)
	webGroup.GET("/user/page", my.userManage.PageUsers)
	webGroup.GET("/user/data-perms/all", my.userManage.DataPermsAll)
	webGroup.GET("/user/data-perms/by-userid", my.userManage.QueryDataPermsByUserId)
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
	// 修改Gin模板定界符
	engine.Delims("{[{", "}]}")
	engine.LoadHTMLGlob("web/templates/*")
	engine.Static("/static", "./web/static")

	// 页面跳转路由
	my.htmlJump(engine)

	// 注册路由: openapi
	apiGroup := engine.Group("/api")
	my.apiGroup(apiGroup)

	// 注册路由: web（后端接口）
	webGroup := engine.Group("/admin")
	my.webGroup(webGroup)

	port := consts.DefaultPort
	if adminConfig.Port != "" {
		port = adminConfig.Port
	}
	global.SingletonPool().Log.Info(fmt.Sprintf("server started on port(s): %s (http) with context path '/'", port))
	_ = engine.Run(":" + port)
}

// 页面跳转路由
func (my *GinRouter) htmlJump(engine *gin.Engine) {
	// 首页
	engine.GET("/", my.toHtml.ToIndex)
	engine.GET("/home", my.toHtml.ToHome)

	// 任务管理
	engine.GET("/job", my.toHtml.ToJob)
	engine.GET("/job-add-layer", my.toHtml.ToJobAddLayer)
	engine.GET("/job-edit-layer", my.toHtml.ToJobEditLayer)
	//engine.GET("/job-run-once-layer", my.toHtml.ToJobRunOnceLayer)

	// 任务调度日志
	engine.GET("/log", my.toHtml.ToLog)
	engine.GET("/log-clear-layer", my.toHtml.ToLogClearLayer)

	// 执行器管理
	engine.GET("/executor", my.toHtml.ToExecutor)
	engine.GET("/executor-add-layer", my.toHtml.ToExecutorAddLayer)
	engine.GET("/executor-edit-layer", my.toHtml.ToExecutorEditLayer)
	engine.GET("/executor-nodes-info-layer", my.toHtml.ToExecutorNodesInfoLayer)

	// 用户管理
	engine.GET("/user", my.toHtml.ToUser)
	engine.GET("/user-add-layer", my.toHtml.ToUserAddLayer)
	engine.GET("/user-edit-layer", my.toHtml.ToUserEditLayer)
	engine.GET("/user-password", my.toHtml.ToUserPassword)
}
