package controllers

import (
	comm "durl/app/share/comm"
	"durl/app/share/dao/db"
	"durl/app/share/tool"
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
)

type setShortUrlReq struct {
	Url            string `form:"url" valid:"Required"`
	ExpirationTime int    `form:"expirationTime"`
}

type setShortUrlDataResp struct {
	Key            string `json:"key"`
	Url            string `json:"url"`
	ExpirationTime int    `json:"expirationTime"`
}

// 函数名称: setShortUrlParam
// 功能: 效验SetShortUrl接口请求参数
// 输入参数:
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021-11-17 15:15:42 #

func (c *Controller) setShortUrlParam(req *setShortUrlReq) {
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, req); err != nil {
		logs.Info("Action setShortUrlParam, err: ", err.Error())
		c.ErrorMessage(comm.ErrParamMiss, comm.MsgParseFormErr)
	}

	if err := c.ParseForm(req); err != nil {
		logs.Info("Action setShortUrlParam, err: ", err.Error())
		c.ErrorMessage(comm.ErrParamMiss, comm.MsgParseFormErr)
	}

	valid := validation.Validation{}
	b, err := valid.Valid(req)
	if err != nil {
		logs.Info("Action setShortUrlParam, err: ", err.Error())
		c.ErrorMessage(comm.ErrParamMiss, comm.MsgParseFormErr)
	}
	if !b {
		for _, err := range valid.Errors {
			logs.Info("Action setShortUrlParam, err: ", err.Key, err.Message)
		}
		c.ErrorMessage(comm.ErrParamInvalid, comm.MsgParseFormErr)
	}
}

// 函数名称: SetShortUrl
// 功能: 根据 单个url 设置短链
// 输入参数:
//		url: 原始url
//		expirationTime: 过期时间
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021-11-17 15:15:42 #

func (c *Controller) SetShortUrl() {

	req := setShortUrlReq{}
	// 效验请求参数格式
	c.setShortUrlParam(&req)

	// 处理url
	req.Url = tool.DisposeUrlProto(req.Url)

	// 消耗池中的短链
	shortNum := ReturnShortNumOne()

	// 数据放入数据库
	var UrlOne db.InsertUrlOneReq
	UrlOne.ShortNum = shortNum
	UrlOne.FullUrl = req.Url
	UrlOne.ExpirationTime = req.ExpirationTime
	err := db.InsertUrlOne(&UrlOne)
	if err != nil {
		logs.Error("Action SetShortUrl, err: ", err.Error())
		c.ErrorMessage(comm.ErrSysDb, comm.MsgNotOk)
		return
	}

	// 拼接url
	shortKey := tool.Base62Encode(shortNum)

	data := &setShortUrlDataResp{
		Url:            req.Url,
		Key:            shortKey,
		ExpirationTime: req.ExpirationTime,
	}

	c.FormatInterfaceResp(comm.OK, comm.OK, comm.MsgOk, data)
	return
}
