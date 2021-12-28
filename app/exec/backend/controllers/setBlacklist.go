package controllers

import (
	comm "durl/app/share/comm"
	"durl/app/share/dao/db"
	"durl/app/share/dao/db/xormDb"
	"github.com/beego/beego/v2/core/logs"
)

type setBlacklistReq struct {
	Ip string `valid:"IP"`
}

// 函数名称: SetBlacklist
// 功能: 根据 单个fullUrl设置短链
// 输入参数:
//		Ip: Ip地址
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #

func (c *BackendController) SetBlacklist() {
	req := setBlacklistReq{}
	// 效验请求参数格式
	c.BaseCheckParams(&req)

	// 数据放入数据库
	var BlacklistOne db.InsertBlacklistOneReq
	BlacklistOne.Ip = req.Ip
	err := db.NewDbService(xormDb.Engine).InsertBlacklistOne(&BlacklistOne)
	if err != nil {
		logs.Error("Action SetBlacklist, err: ", err.Error())
		c.ErrorMessage(comm.ErrSysDb, comm.MsgNotOk)
		return
	}

	c.FormatResp(comm.OK, comm.OK, comm.MsgOk)
	return
}
