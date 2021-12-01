package routers

import (
	"durl/app/exec/backend/controllers"
	"github.com/beego/beego/v2/server/web"
)

// RouterHandler 路由跳转
func RouterHandler() {

	// backendapi初始化
	controllers.InitCon()

	// 获取xsrfToken
	web.Router("/xsrf-token", &controllers.Controller{}, "get:GetXsrfToken")

	// 设置短链
	web.Router("/url", &controllers.Controller{}, "post:SetShortUrl")

	// 修改短链
	web.Router("/url/:id([0-9a-zA-Z]+)", &controllers.Controller{}, "put:UpdateShortUrl")

	// 删除短链
	web.Router("/url/:id([0-9a-zA-Z]+)", &controllers.Controller{}, "delete:DelShortUrl")

	// 批量删除短链
	web.Router("/url", &controllers.Controller{}, "delete:BatchDelShortUrl")

	// 短链列表
	web.Router("/url/list", &controllers.Controller{}, "get:GetShortUrlList")

	// 短链详情
	web.Router("/url/info/:id([0-9a-zA-Z]+)", &controllers.Controller{}, "get:GetShortUrlInfo")

	// 冻结Url
	web.Router("/url/frozen/:id([0-9a-zA-Z]+)", &controllers.Controller{}, "put:FrozenShortUrl")

	// 批量冻结Url
	web.Router("/url/frozen", &controllers.Controller{}, "put:BatchFrozenShortUrl")

}
