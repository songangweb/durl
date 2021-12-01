package controllers

import (
	comm "durl/app/share/comm"
	"durl/app/share/dao/db"
	"github.com/beego/beego/v2/core/logs"
)

// 函数名称: DelShortUrl
// 功能: 根据单个短链删除短链接
// 输入参数:
//     key: 短链结果
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/18 5:44 下午 #

func (c *Controller) DelShortUrl() {

	id := c.Ctx.Input.Param(":id")

	// 查询此短链
	where := make(map[string][]interface{})
	where["id"] = append(where["id"], "=", id)
	where["is_del"] = append(where["is_del"], "=", 0)
	urlInfo := db.GetShortUrlInfo(where)
	if urlInfo.ShortNum == 0 {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}

	// 删除此短链
	_, err := db.DelUrlById(id, urlInfo.ShortNum)
	if err != nil {
		logs.Error("Action DelShortKey, err: ", err.Error())
		c.ErrorMessage(comm.ErrSysDb, comm.MsgNotOk)
		return
	}
	c.FormatResp(comm.OK, comm.OK, comm.MsgOk)
	return
}
