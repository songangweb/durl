package routers

import (
	"durl/app/exec/openapi/controllers"

	"github.com/beego/beego/v2/server/web"
)

// RouterHandler 路由跳转
func RouterHandler() {
	// openapi初始化
	controllers.InitCon()

	// 设置短链
	web.Router("/url", &controllers.OpenApiController{}, "post:SetShortUrl")
	// 修改短链
	web.Router("/url/:id([0-9a-zA-Z]+)", &controllers.OpenApiController{}, "put:UpdateShortUrl")
	// 删除短链
	web.Router("/url/:id([0-9a-zA-Z]+)", &controllers.OpenApiController{}, "delete:DelShortUrl")
	// 短链详情
	web.Router("/url/:id([0-9a-zA-Z]+)", &controllers.OpenApiController{}, "get:GetShortUrlInfo")
	// 批量删除短链
	web.Router("/url", &controllers.OpenApiController{}, "delete:BatchDelShortUrl")
	// 短链列表
	web.Router("/url/list", &controllers.OpenApiController{}, "get:GetShortUrlList")
	// 冻结Url
	web.Router("/url/frozen/:id([0-9a-zA-Z]+)", &controllers.OpenApiController{}, "put:FrozenShortUrl")
	// 批量冻结Url
	web.Router("/url/frozen", &controllers.OpenApiController{}, "put:BatchFrozenShortUrl")

	// 设置黑名单
	web.Router("/blacklist", &controllers.OpenApiController{}, "post:SetBlacklist")
	// 修改黑名单
	web.Router("/blacklist/:id([0-9a-zA-Z]+)", &controllers.OpenApiController{}, "put:UpdateBlacklist")
	// 删除黑名单
	web.Router("/blacklist/:id([0-9a-zA-Z]+)", &controllers.OpenApiController{}, "delete:DelBlacklist")
	// 黑名单详情
	web.Router("/blacklist/:id([0-9]+)", &controllers.OpenApiController{}, "get:GetBlacklistInfo")
	// 黑名单列表
	web.Router("/blacklist/list", &controllers.OpenApiController{}, "get:GetBlacklistList")

}
