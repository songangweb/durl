package controllers

import (
	comm "durl/app/share/comm"
	"durl/app/share/dao/db"
	"durl/app/share/tool"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
)

type GetXsrfTokenResp struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Token string `json:"token"`
}

// GetXsrfToken 获取token
func (c *Controller) GetXsrfToken() {
	c.Data["json"] = &GetXsrfTokenResp{
		Code:  comm.OK,
		Msg:   comm.MsgOk,
		Token: c.XSRFToken(),
	}
	_ = c.ServeJSON()
	return
}

type setShortUrlReq struct {
	Url            string `form:"url" valid:"Required"`
	ExpirationTime int    `form:"expirationTime"`
}

type setShortUrlResp struct {
	Code int                  `json:"code"`
	Msg  string               `json:"msg"`
	Data *setShortUrlDataResp `json:"data"`
}

type setShortUrlDataResp struct {
	Key            string `json:"key"`
	Url            string `json:"url"`
	ExpirationTime int    `json:"expirationTime"`
}

// 效验请求过来的参数
func (c *Controller) setShortUrlParam(req *setShortUrlReq) bool {
	if err := c.ParseForm(req); err != nil {
		logs.Info("Action setShortUrlParam, err: ", err.Error())
		c.Data["json"] = &setShortUrlResp{
			Code: comm.ErrParamMiss,
			Msg:  comm.MsgParseFormErr,
		}
		return false
	}

	valid := validation.Validation{}
	b, err := valid.Valid(req)
	if err != nil {
		logs.Info("Action setShortUrlParam, err: ", err.Error())
		c.Data["json"] = &setShortUrlResp{
			Code: comm.ErrParamMiss,
			Msg:  comm.MsgParseFormErr,
		}
		return false
	}
	if !b {
		for _, err := range valid.Errors {
			logs.Info("Action setShortUrlParam, err: ", err.Key, err.Message)
		}
		c.Data["json"] = &setShortUrlResp{
			Code: comm.ErrParamInvalid,
			Msg:  comm.MsgParseFormErr,
		}
		return false
	}
	return true
}

// SetShortUrl 根据 单个url 设置短链
func (c *Controller) SetShortUrl() {

	req := setShortUrlReq{}
	// 效验请求参数格式
	if !c.setShortUrlParam(&req) {
		_ = c.ServeJSON()
		return
	}

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
		c.Data["json"] = &setShortUrlResp{
			Code: comm.ErrSysDb,
			Msg:  comm.MsgNotOk,
		}
		_ = c.ServeJSON()
		return
	}

	// 拼接url
	shortKey := tool.Base62Encode(shortNum)
	durl := c.Ctx.Request.Host + "/" + shortKey

	c.Data["json"] = &setShortUrlResp{
		Code: comm.OK,
		Msg:  comm.MsgOk,
		Data: &setShortUrlDataResp{
			Url:            req.Url,
			Key:            shortKey,
			ExpirationTime: req.ExpirationTime,
		},
	}
	_ = c.ServeJSON()
	return
}
