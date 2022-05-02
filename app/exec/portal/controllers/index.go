package controllers

import (
	"github.com/beego/beego/v2/core/config"
)

// Index
// 函数名称: 首页入口
// 功能:
// 输入参数:
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021-11-17 15:15:42 #
func (c *Controller) Index() {
	// xsrf 值
	c.Data["xsrf_token"] = c.XSRFToken()

	// jump服务域名
	jumpUrlConf, _ := config.String("jumpUrl")
	if jumpUrlConf != "" {
		c.Data["JumpUrl"] = jumpUrlConf
	}

	c.TplName = "index.html"
}
