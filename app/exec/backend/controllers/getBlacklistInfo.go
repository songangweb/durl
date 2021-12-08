package controllers

import (
	comm "durl/app/share/comm"
	"durl/app/share/dao/db"
)

type BlacklistInfoRes struct {
	Id         interface{} `json:"id"`
	Ip         string      `json:"ip"`
	CreateTime int         `json:"createTime"`
	UpdateTime int         `json:"updateTime"`
}

// 函数名称: GetBlacklistInfo
// 功能: 获取url详情
// 输入参数:
//     id: 数据id
// 输出参数:
// 返回: 短链详情
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #

func (c *BackendController) GetBlacklistInfo() {
	id := c.Ctx.Input.Param(":id")

	fields := map[string]interface{}{"id": id}
	BlacklistInfo := db.GetBlacklistInfo(fields)
	if BlacklistInfo.Ip == "" {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
	}

	c.FormatInterfaceResp(comm.OK, comm.OK, comm.MsgOk, BlacklistInfoRes{
		Id:             BlacklistInfo.Id,
		Ip:        BlacklistInfo.Ip,
		CreateTime:     BlacklistInfo.CreateTime,
		UpdateTime:     BlacklistInfo.UpdateTime,
	})
	return
}
