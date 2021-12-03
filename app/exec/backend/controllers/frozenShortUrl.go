package controllers

import (
	comm "durl/app/share/comm"
	"durl/app/share/dao/db"
)

// 函数名称: FrozenShortUrl
// 功能: 冻结ShortUrl
// 输入参数:
//     id: 短链id
// 输出参数:
// 返回: 冻结操作结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/26 1:56 下午 #

func (c *Controller) FrozenShortUrl() {

	id := c.Ctx.Input.Param(":id")

	// 查询此短链
	fields := map[string]interface{}{"id": id}
	urlInfo := db.GetShortUrlInfo(fields)
	if urlInfo.ShortNum == 0 {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}

	// 冻结/解冻ShortUrl
	updateData := make(map[string]interface{})
	if urlInfo.IsFrozen == 0 {
		updateData["is_frozen"] = 1
	} else {
		updateData["is_frozen"] = 0
	}

	_, err := db.UpdateUrlById(id, urlInfo.ShortNum, updateData)
	if err != nil {
		c.ErrorMessage(comm.ErrSysDb, comm.MsgNotOk)
		return
	}
	c.FormatResp(comm.OK, comm.OK, comm.MsgOk)
	return
}
