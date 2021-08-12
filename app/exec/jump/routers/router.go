package routers

import (
	"durl/app/exec/jump/controllers"
	"github.com/beego/beego/v2/server/web"
	"html/template"
	"net/http"
)

type Conf struct {
	OpenApi bool
}

// RouterHandler 路由跳转
func (c Conf) RouterHandler() {

	// 链接跳转
	web.Router("/:jump([0-9a-zA-Z]+)", &controllers.Controller{}, "*:Jump")

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
