package controllers

import (
	"durl/app/share/comm"
	"durl/app/share/dao/db"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
)

// DelShortUrl
// 函数名称: DelShortUrl
// 功能: 根据单个短链删除短链接
// 输入参数:
//     id: 数据id
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/18 5:44 下午 #
func (c *OpenApiController) DelShortUrl() {

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

	// 删除此短链
	_, err := engine.DelUrlById(intId, urlInfo.ShortNum)
	if err != nil {
		logs.Error("Action DelShortKey, err: ", err.Error())
		c.ErrorMessage(comm.ErrSysDb, comm.MsgNotOk)
		return
	}
	c.FormatResp(comm.OK, comm.OK, comm.MsgOk)
	return
}
