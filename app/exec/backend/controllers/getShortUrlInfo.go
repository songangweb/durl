package controllers

import (
	"durl/app/share/comm"
	"durl/app/share/dao/db"
	"durl/app/share/tool"
	"strconv"
)

type ShortUrlInfoRes struct {
	Id             uint32 `json:"id"`
	ShortKey       string `json:"shortKey"`
	FullUrl        string `json:"fullUrl"`
	ExpirationTime uint32 `json:"expirationTime"`
	IsFrozen       uint8  `json:"isFrozen"`
	CreateTime     uint32 `json:"createTime"`
	UpdateTime     uint32 `json:"updateTime"`
}

// 函数名称: GetShortUrlInfo
// 功能: 获取url详情
// 输入参数:
//     id: 数据id
// 输出参数:
// 返回: 短链详情
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/26 10:59 上午 #

func (c *BackendController) GetShortUrlInfo() {
	id := c.Ctx.Input.Param(":id")
	intId, _ := strconv.ParseUint(id, 10, 32)
	uint32Id := uint32(intId)

	// 查询此短链
	fields := map[string]interface{}{"id": uint32Id}

	urlInfo := db.NewDbService().GetShortUrlInfo(fields)
	if urlInfo.ShortNum == 0 {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}
	// 短链转化
	shortKey := tool.Base62Encode(urlInfo.ShortNum)
	c.FormatInterfaceResp(comm.OK, comm.OK, comm.MsgOk, ShortUrlInfoRes{
		Id:             urlInfo.Id,
		ShortKey:       shortKey,
		FullUrl:        urlInfo.FullUrl,
		ExpirationTime: urlInfo.ExpirationTime,
		IsFrozen:       urlInfo.IsFrozen,
		CreateTime:     urlInfo.CreateTime,
		UpdateTime:     urlInfo.UpdateTime,
	})
	return
}
