package routers

import (
	"durl/app/exec/backend/controllers"
	"github.com/beego/beego/v2/server/web"
)

// RouterHandler 路由跳转
func RouterHandler() {

	// backendapi初始化
	controllers.InitCon()

	// 设置短链
	web.Router("/url", &controllers.BackendController{}, "post:SetShortUrl")
	// 修改短链
	web.Router("/url/:id([0-9a-zA-Z]+)", &controllers.BackendController{}, "put:UpdateShortUrl")
	// 删除短链
	web.Router("/url/:id([0-9a-zA-Z]+)", &controllers.BackendController{}, "delete:DelShortUrl")
	// 短链详情
	web.Router("/url/:id([0-9a-zA-Z]+)", &controllers.BackendController{}, "get:GetShortUrlInfo")
	// 批量删除短链
	web.Router("/url", &controllers.BackendController{}, "delete:BatchDelShortUrl")
	// 短链列表
	web.Router("/url/list", &controllers.BackendController{}, "get:GetShortUrlList")
	// 冻结Url
	web.Router("/url/frozen/:id([0-9a-zA-Z]+)", &controllers.BackendController{}, "put:FrozenShortUrl")
	// 批量冻结Url
	web.Router("/url/frozen", &controllers.BackendController{}, "put:BatchFrozenShortUrl")

	// 设置黑名单
	web.Router("/blacklist", &controllers.BackendController{}, "post:SetBlacklist")
	// 修改黑名单
	web.Router("/blacklist/:id([0-9a-zA-Z]+)", &controllers.BackendController{}, "put:UpdateBlacklist")
	// 删除黑名单
	web.Router("/blacklist/:id([0-9a-zA-Z]+)", &controllers.BackendController{}, "delete:DelBlacklist")
	// 黑名单详情
	web.Router("/blacklist/:id([0-9]+)", &controllers.BackendController{}, "get:GetBlacklistInfo")
	// 黑名单列表
	web.Router("/blacklist/list", &controllers.BackendController{}, "get:GetBlacklistList")

}
