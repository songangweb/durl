package routers

import (
	"durl/app/exec/openApi/controllers"
	"github.com/beego/beego/v2/server/web"
)

type Conf struct {
	OpenApi bool
}

// RouterHandler 路由跳转
func (c Conf) RouterHandler() {

	// openApi初始化
	controllers.InitCon()

	// 获取xsrfToken
	web.Router("/xsrf-token", &controllers.Controller{}, "get:GetXsrfToken")

	// 设置短链
	web.Router("/url", &controllers.Controller{}, "post:SetShortUrl")

	// 修改短链
	web.Router("/url", &controllers.Controller{}, "put:UpdateShortUrl")

	// 删除短链
	web.Router("/url", &controllers.Controller{}, "delete:DelShortKey")

}
