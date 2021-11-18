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
	Durl           string `json:"durl"`
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
			Durl:           durl,
			ExpirationTime: req.ExpirationTime,
		},
	}
	_ = c.ServeJSON()
	return
}

type delShortKeyReq struct {
	Key string `form:"key" valid:"Required"`
}

type delShortKeyResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// 效验请求过来的参数
func (c *Controller) delShortKeyParam(req *delShortKeyReq) bool {
	if err := c.ParseForm(req); err != nil {
		logs.Info("Action delShortKeyParam, err: ", err.Error())
		c.Data["json"] = &delShortKeyResp{
			Code: comm.ErrParamMiss,
			Msg:  comm.MsgParseFormErr,
		}
		return false
	}

	valid := validation.Validation{}
	b, err := valid.Valid(req)
	if err != nil {
		logs.Info("Action delShortKeyParam, err: ", err.Error())
		c.Data["json"] = &delShortKeyResp{
			Code: comm.ErrParamMiss,
			Msg:  comm.MsgParseFormErr,
		}
		return false
	}
	if !b {
		for _, err := range valid.Errors {
			logs.Info("Action delShortKeyParam, err: ", err.Key, err.Message)
		}
		c.Data["json"] = &delShortKeyResp{
			Code: comm.ErrParamInvalid,
			Msg:  comm.MsgParseFormErr,
		}
		return false
	}
	return true
}

// DelShortKey 删除某个短链
func (c *Controller) DelShortKey() {

	req := delShortKeyReq{}
	// 效验请求参数格式
	if !c.delShortKeyParam(&req) {
		_ = c.ServeJSON()
		return
	}

	shortKey := req.Key

	if !tool.DisposeShortKey(req.Key) {
		c.Data["json"] = &delShortKeyResp{
			Code: comm.ErrParamInvalid,
			Msg:  comm.MsgParseFormErr,
		}
		_ = c.ServeJSON()
		return
	}

	shortNum := tool.Base62Decode(shortKey)

	// 删除此短链
	_, err := db.DelUrlByShortNum(shortNum)
	if err != nil {
		logs.Error("Action DelShortKey, err: ", err.Error())
		c.Data["json"] = &delShortKeyResp{
			Code: comm.ErrSysDb,
			Msg:  comm.MsgNotOk,
		}
		_ = c.ServeJSON()
		return
	}

	c.Data["json"] = &delShortKeyResp{
		Code: comm.OK,
		Msg:  comm.MsgOk,
	}
	_ = c.ServeJSON()
	return
}

type updateShortUrlReq struct {
	Key            string `form:"key"  valid:"Required"`
	Url            string `form:"url"`
	IsFrozen       int    `form:"isFrozen"`
	ExpirationTime int64  `form:"expirationTime"`
}

type updateShortUrlResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// 效验请求过来的参数
func (c *Controller) updateShortKeyParam(req *updateShortUrlReq) bool {
	if err := c.ParseForm(req); err != nil {
		logs.Info("Action updateShortKeyParam, err: ", err.Error())
		c.Data["json"] = &delShortKeyResp{
			Code: comm.ErrParamMiss,
			Msg:  comm.MsgParseFormErr,
		}
		return false
	}

	valid := validation.Validation{}
	b, err := valid.Valid(req)
	if err != nil {
		logs.Info("Action updateShortKeyParam, err: ", err.Error())
		c.Data["json"] = &delShortKeyResp{
			Code: comm.ErrParamMiss,
			Msg:  comm.MsgParseFormErr,
		}
		return false
	}
	if !b {
		for _, err := range valid.Errors {
			logs.Info("Action updateShortKeyParam, err: ", err.Key, err.Message)
		}
		c.Data["json"] = &delShortKeyResp{
			Code: comm.ErrParamInvalid,
			Msg:  comm.MsgParseFormErr,
		}
		return false
	}
	return true
}

// UpdateShortUrl 修改某个短链
func (c *Controller) UpdateShortUrl() {

	req := updateShortUrlReq{}
	// 效验请求参数格式
	if !c.updateShortKeyParam(&req) {
		_ = c.ServeJSON()
		return
	}

	// 效验 key
	shortKey := req.Key
	if !tool.DisposeShortKey(req.Key) {
		c.Data["json"] = &updateShortUrlResp{
			Code: comm.ErrParamInvalid,
			Msg:  comm.MsgParseFormErr,
		}
		_ = c.ServeJSON()
		return
	}
	shortNum := tool.Base62Decode(shortKey)

	// 初始化需要更新的内容
	updateData := make(map[string]interface{})

	// 接收请求数据
	expirationTimeInt, err := c.GetInt("expirationTime")
	if err == nil {
		updateData["expirationTime"] = expirationTimeInt
	}

	urlStr := c.GetString("url")
	if urlStr != "" {
		updateData["url"] = urlStr
	}

	isFrozenInt, err := c.GetInt("isFrozen")
	if err == nil {
		updateData["isFrozen"] = isFrozenInt
	}

	if len(updateData) == 0 {
		c.Data["json"] = &delShortKeyResp{
			Code: comm.ErrParamMiss,
			Msg:  comm.MsgParseFormErr,
		}
		_ = c.ServeJSON()
		return
	}

	// 修改此短链信息
	_, err = db.UpdateUrlByShortNum(shortNum, &updateData)
	if err != nil {
		c.Data["json"] = &updateShortUrlResp{
			Code: comm.ErrSysDb,
			Msg:  comm.MsgNotOk,
		}
		_ = c.ServeJSON()
		return
	}

	c.Data["json"] = &updateShortUrlResp{
		Code: comm.OK,
		Msg:  comm.MsgOk,
	}
	_ = c.ServeJSON()
	return
}

type getShortUrlListReq struct {
	Url  string `form:"url"`
	Page int    `form:"page"`
	Size int    `form:"size"`
}

type getShortUrlListResp struct {
	Code  int                       `json:"code"`
	Msg   string                    `json:"msg"`
	Data  []*getShortUrlListDataResp `json:"data"`
	Total int64                     `json:"total""`
}

type getShortUrlListDataResp struct {
	Id             int    `json:"id"`
	ShortNum       int    `json:"shor_url"`
	FullUrl        string `json:"full_url"`
	ExpirationTime int    `json:"expiration_time"`
	IsFrozen       int8   `json:"is_frozen"`
	CreateTime     int    `json:"create_time"`
}

// 效验请求过来的参数
func (c *Controller) getShortUrlListParam(req *getShortUrlListReq) bool {
	if err := c.ParseForm(req); err != nil {
		logs.Info("Action getShortUrlListParam, err: ", err.Error())
		c.Data["json"] = &getShortUrlListResp{
			Code: comm.ErrParamMiss,
			Msg:  comm.MsgParseFormErr,
		}
		return false
	}

	valid := validation.Validation{}
	b, err := valid.Valid(req)
	if err != nil {
		logs.Info("Action getShortUrlListParam, err: ", err.Error())
		c.Data["json"] = &getShortUrlListResp{
			Code: comm.ErrParamMiss,
			Msg:  comm.MsgParseFormErr,
		}
		return false
	}
	if !b {
		for _, err := range valid.Errors {
			logs.Info("Action getShortUrlListResp, err: ", err.Key, err.Message)
		}
		c.Data["json"] = &getShortUrlListResp{
			Code: comm.ErrParamInvalid,
			Msg:  comm.MsgParseFormErr,
		}
		return false
	}
	return true
}

// 获取url 列表数据
func (c *Controller) GetShortUrlList() {
	req := getShortUrlListReq{}
	// 效验请求参数格式
	if !c.getShortUrlListParam(&req) {
		_ = c.ServeJSON()
		return
	}
	// page 需要给默认值
	if req.Page == 0 {
		req.Page = 1
	}
	// size 需要默认值且不可大于500
	if req.Size == 0 || req.Size > 500 {
		req.Size = 20
	}

	where := " is_del =? "
	bindValue :=[]interface{}{0}
	// key 是原url模糊搜索
	// 拼接条件进行检索
	if req.Url != "" {
		where += " and full_url like ? "
		bindValue = append(bindValue,"%"+req.Url+"%")

	}
	data := db.GetShortUrlList(where, req.Page, req.Size,bindValue...)
	var total int64
	// 有数据且当page=1时计算结果总条数
	if data != nil && req.Page == 1 {
		// 统计结果总条数
		total = db.GetShortUrlListTotal(where,bindValue...)
	}

	var list []*getShortUrlListDataResp
	for _, queueStruct := range data {
		var One getShortUrlListDataResp
		One.Id = queueStruct.Id
		One.ShortNum = queueStruct.ShortNum
		One.FullUrl = queueStruct.FullUrl
		One.ExpirationTime = queueStruct.ExpirationTime
		One.IsFrozen = queueStruct.IsFrozen
		One.CreateTime = queueStruct.CreateTime
		list = append(list, &One)
	}
	c.Data["json"] = &getShortUrlListResp{
		Code:  comm.OK,
		Msg:   comm.MsgOk,
		Data:  list,
		Total: total,
	}
	_ = c.ServeJSON()
	return

}
