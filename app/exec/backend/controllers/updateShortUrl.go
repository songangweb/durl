package controllers

import (
	"durl/app/share/comm"
	"durl/app/share/dao/db"
	"strconv"
)

type updateShortUrlReq struct {
	FullUrl        string `form:"fullUrl" valid:"Required"`
	IsFrozen       uint8  `form:"isFrozen" valid:"Range(0,1)"`
	ExpirationTime uint32 `form:"expirationTime" valid:"Match(/([0-9]{10}$)|([0])/);Max(9999999999)"`
}

// 函数名称: UpdateShortUrl
// 功能: 根据短链修改短链接信息
// 输入参数:
//	   fullUrl: 原始url
//	   isFrozen: 是否冻结
//	   expirationTime: 过期时间
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/18 5:46 下午 #

func (c *BackendController) UpdateShortUrl() {

	req := updateShortUrlReq{}
	// 效验请求参数格式
	c.BaseCheckParams(&req)

	id := c.Ctx.Input.Param(":id")
	intId, _ := strconv.ParseUint(id, 10, 32)
	uint32Id := uint32(intId)

	// 查询此短链
	fields := map[string]interface{}{"id": uint32Id}
	engine := db.NewDbService()
	urlInfo := engine.GetShortUrlInfo(fields)
	if urlInfo.ShortNum == 0 {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}

	// 初始化需要更新的内容
	updateData := make(map[string]interface{})
	updateData["expiration_time"] = req.ExpirationTime
	updateData["full_url"] = req.FullUrl
	updateData["is_frozen"] = req.IsFrozen

	// 修改此短链信息
	_, err := engine.UpdateUrlById(uint32Id, urlInfo.ShortNum, updateData)
	if err != nil {
		c.ErrorMessage(comm.ErrSysDb, comm.MsgNotOk)
		return
	}

	c.FormatResp(comm.OK, comm.OK, comm.MsgOk)
	return
}
