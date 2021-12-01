package controllers

import (
	comm "durl/app/share/comm"
	"durl/app/share/dao/db"
)

// 函数名称: GetShortUrlInfo
// 功能: 获取url详情
// 输入参数:
//     id: 短链id
// 输出参数:
// 返回: 短链详情
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/26 10:59 上午 #

func (c *Controller) GetShortUrlInfo() {
	id := c.Ctx.Input.Param(":id")
	// 查询此短链
	where := make(map[string][]interface{})
	where["id"] = append(where["id"], "=", id)
	urlInfo := db.GetShortUrlInfo(where)
	if urlInfo.ShortNum == 0 {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}
	c.FormatInterfaceResp(comm.OK, comm.OK, comm.MsgOk, urlInfo)
	return
}
