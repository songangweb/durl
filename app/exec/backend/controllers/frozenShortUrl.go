package controllers

import (
	"durl/app/share/comm"
	"durl/app/share/dao/db"
	"strconv"
)

// FrozenShortUrl
// 函数名称: FrozenShortUrl
// 功能: 冻结ShortUrl
// 输入参数:
//     id: 数据id
// 输出参数:
// 返回: 冻结操作结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/26 1:56 下午 #
func (c *BackendController) FrozenShortUrl() {

	id := c.Ctx.Input.Param(":id")
	intId, _ := strconv.Atoi(id)

	// 查询此短链
	fields := map[string]interface{}{"id": intId}
	engine := db.NewDbService()
	urlInfo := engine.GetShortUrlInfo(fields)
	if urlInfo.ShortNum == 0 {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}

	// 冻结/解冻ShortUrl
	if urlInfo.IsFrozen == 1 {
		c.FormatResp(comm.OK, comm.OK, comm.MsgOk)
		return
	}

	updateData := make(map[string]interface{})
	updateData["is_frozen"] = 1
	_, err := engine.UpdateUrlById(intId, urlInfo.ShortNum, updateData)
	if err != nil {
		c.ErrorMessage(comm.ErrSysDb, comm.MsgNotOk)
		return
	}
	c.FormatResp(comm.OK, comm.OK, comm.MsgOk)
	return
}
