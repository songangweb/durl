package routers

import (
	"durl/controllers/jump"
	"durl/controllers/openApi"
	"github.com/beego/beego/v2/server/web"
	"html/template"
	"net/http"
)

type Conf struct {
	OpenApi bool
}

// RouterHandler 路由跳转
func (c Conf) RouterHandler() {

	// 判断是否开启 openApi 接口
	if c.OpenApi {

		// openApi初始化
		openApi.InitCon()

		// 首页
		web.Router("/", &openApi.Controller{}, "get:Index")
		web.Router("/index", &openApi.Controller{}, "get:Index")

		// 获取xsrfToken
		web.Router("/xsrf-token", &openApi.Controller{}, "get:GetXsrfToken")

		// 设置短链
		web.Router("/url", &openApi.Controller{}, "post:SetShortUrl")

		// 修改短链
		web.Router("/url", &openApi.Controller{}, "put:UpdateShortUrl")

		// 删除短链
		web.Router("/url", &openApi.Controller{}, "delete:DelShortKey")

	}

	// 链接跳转
	web.Router("/:jump([0-9a-zA-Z]+)", &jump.Controller{}, "*:Jump")

	web.ErrorHandler("404", pageNotFound)

}

// 定义404页面
func pageNotFound(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("404.html").ParseFiles(web.BConfig.WebConfig.ViewsPath + "/404.html")
	data := make(map[string]interface{})
	data["content"] = "page not found"
	//// 百度统计key
	//sConf, _ := config.String("Statistical_Baidu")
	//if sConf != "" {
	//	data["Statistical_Baidu_Key"] = sConf
	//}
	_ = t.Execute(rw, data)
}
