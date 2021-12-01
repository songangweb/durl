package controllers

import (
	comm "durl/app/share/comm"
)

type getXsrfTokenData struct {
	Token string `json:"token"`
}

// 函数名称: GetXsrfToken
// 功能: 获取XsrfToken
// 输入参数:
// 输出参数:
//		token
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021-11-17 15:15:42 #

func (c *Controller) GetXsrfToken() {
	data := &getXsrfTokenData{
		Token: c.XSRFToken(),
	}
	c.FormatInterfaceResp(comm.OK, comm.OK, comm.MsgOk, data)
	return
}
