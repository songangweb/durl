package controllers

import (
	"github.com/beego/beego/v2/core/config"
)

// Index 首页入口
func (c *Controller) Index() {
	// xsrf 值
	c.Data["xsrf_token"] = c.XSRFToken()

	// 百度统计key
	runmode, _ := config.String("runmode")
	sConf, _ := config.String(runmode + "::Baidu")
	if sConf != "" {
		c.Data["Statistical_Baidu_Key"] = sConf
	}

	// jump服务域名
	jumpUrlConf, _ := config.String("jumpUrl")
	if jumpUrlConf != "" {
		c.Data["JumpUrl"] = jumpUrlConf
	}

	c.TplName = "index.html"
}
