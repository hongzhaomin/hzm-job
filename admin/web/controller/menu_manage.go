package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongzhaomin/hzm-job/admin/vo"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"net/http"
)

// MenuManage 菜单管理
type MenuManage struct {
}

// Init 菜单初始化
// @Get /admin/menu/init
func (my *MenuManage) Init(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, sdk.Ok2[gin.H](gin.H{
		"homeInfo": vo.MenuInfo{
			Id:    1,
			Title: "首页",
			Href:  "home",
		},
		"logoInfo": map[string]any{
			"title": "HZM-JOB",
			"image": "../static/images/logo.jpg",
			"href":  "",
		},
		"menuInfo": []vo.MenuInfo{
			{
				Id:     10,
				Title:  "任务管理",
				Href:   "job",
				Icon:   "fa fa-th-large",
				Target: "_self",
			},
			my.GetMenu4Log(),
			{
				Id:     30,
				Title:  "执行器管理",
				Href:   "executor",
				Icon:   "fa fa-laptop",
				Target: "_self",
			},
			{
				Id:     40,
				Title:  "用户管理",
				Href:   "user",
				Icon:   "fa fa-user",
				Target: "_self",
				//Children: []vo.MenuInfo{
				//	{
				//		Id:       41,
				//		ParentId: 40,
				//		Title:    "用户管理",
				//		Href:     "user",
				//		Target:   "_self",
				//		Icon:     "fa fa-user",
				//	},
				//},
			},
		},
	}))
}

func (my *MenuManage) GetLogMenu(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, sdk.Ok2[vo.MenuInfo](my.GetMenu4Log()))
}

func (my *MenuManage) GetMenu4Log() vo.MenuInfo {
	return vo.MenuInfo{
		Id:     20,
		Title:  "调度日志",
		Href:   "log",
		Icon:   "fa fa-tachometer",
		Target: "_self",
	}
}
