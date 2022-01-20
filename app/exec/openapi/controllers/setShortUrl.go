package controllers

import (
	"durl/app/share/comm"
	"durl/app/share/dao/db"
	"durl/app/share/tool"

	"github.com/beego/beego/v2/core/logs"
)

type setShortUrlReq struct {
	FullUrl        string `form:"fullUrl" valid:"Required"`
	IsFrozen       int    `form:"isFrozen"`
	ExpirationTime int    `form:"expirationTime"`
}

type setShortUrlDataResp struct {
	ShortKey       string `json:"shortKey"`
	FullUrl        string `json:"fullUrl"`
	ExpirationTime int    `json:"expirationTime"`
}

// 函数名称: SetShortUrl
// 功能: 根据 单个fullUrl设置短链
// 输入参数:
//		fullUrl: 原始url
//		expirationTime: 过期时间
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021-11-17 15:15:42 #

func (c *OpenApiController) SetShortUrl() {
	req := setShortUrlReq{}
	// 效验请求参数格式
	c.BaseCheckParams(&req)

	// 处理url
	req.FullUrl = tool.DisposeUrlProto(req.FullUrl)

	// 消耗池中的短链
	shortNum := ReturnShortNumOne()

	// 数据放入数据库
	var UrlOne db.InsertUrlOneReq
	UrlOne.ShortNum = shortNum
	UrlOne.FullUrl = req.FullUrl
	UrlOne.ExpirationTime = req.ExpirationTime
	err := db.NewDbService().InsertUrlOne(&UrlOne)
	if err != nil {
		logs.Error("Action SetShortUrl, err: ", err.Error())
		c.ErrorMessage(comm.ErrSysDb, comm.MsgNotOk)
		return
	}

	// 拼接url
	shortKey := tool.Base62Encode(shortNum)

	data := &setShortUrlDataResp{
		FullUrl:        req.FullUrl,
		ShortKey:       shortKey,
		ExpirationTime: req.ExpirationTime,
	}

	c.FormatInterfaceResp(comm.OK, comm.OK, comm.MsgOk, data)
	return
}
