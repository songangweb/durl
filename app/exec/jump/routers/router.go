package routers

import (
	"html/template"
	"net/http"

	"durl/app/exec/jump/controllers"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/server/web"
)

// RouterHandler 路由跳转
func RouterHandler() {
	// 链接跳转
	web.ErrorHandler("404", pageNotFound)

	web.Router("/:jump([0-9a-zA-Z]+)", &controllers.Controller{}, "*:Jump")
}

// 定义404页面
func pageNotFound(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("404.html").ParseFiles(web.BConfig.WebConfig.ViewsPath + "/404.html")
	data := make(map[string]interface{})
	data["content"] = "page not found"
	runmode, _ := config.String("runmode")
	sConf, _ := config.String(runmode + "::Baidu")
	if sConf != "" {
		data["Statistical_Baidu_Key"] = sConf
	}
	_ = t.Execute(rw, data)
}
