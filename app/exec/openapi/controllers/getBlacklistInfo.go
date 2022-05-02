package controllers

import (
	"durl/app/share/comm"
	"durl/app/share/dao/db"
	"strconv"
)

type BlacklistInfoRes struct {
	Id         int    `json:"id"`
	Ip         string `json:"ip"`
	CreateTime int    `json:"createTime"`
	UpdateTime int    `json:"updateTime"`
}

// GetBlacklistInfo
// 函数名称: GetBlacklistInfo
// 功能: 获取url详情
// 输入参数:
//     id: 数据id
// 输出参数:
// 返回: 短链详情
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #
func (c *OpenApiController) GetBlacklistInfo() {
	id := c.Ctx.Input.Param(":id")
	intId, _ := strconv.Atoi(id)

	fields := map[string]interface{}{"id": intId}
	BlacklistInfo := db.NewDbService().GetBlacklistInfo(fields)
	if BlacklistInfo.Ip == "" {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
	}

	c.FormatInterfaceResp(comm.OK, comm.OK, comm.MsgOk, &BlacklistInfo)
	return
}
