package controllers

import (
	comm "durl/app/share/comm"
	"durl/app/share/dao/db"
	"durl/app/share/dao/db/xormDb"
	"github.com/beego/beego/v2/core/logs"
)

// 函数名称: DelBlacklist
// 功能: 删除黑名单
// 输入参数:
//     id: 数据id
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #

func (c *BackendController) DelBlacklist() {

	id := c.Ctx.Input.Param(":id")

	fields := map[string]interface{}{"id": id}
	engine := db.NewDbService(xormDb.Engine)
	urlInfo := engine.GetBlacklistInfo(fields)
	if urlInfo.Id == 0 {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}

	_, err := engine.DelBlacklistById(id)
	if err != nil {
		logs.Error("Action DelBlacklist, err: ", err.Error())
		c.ErrorMessage(comm.ErrSysDb, comm.MsgNotOk)
		return
	}
	c.FormatResp(comm.OK, comm.OK, comm.MsgOk)
	return
}
