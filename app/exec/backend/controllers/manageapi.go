package controllers

import (
	comm "durl/app/share/comm"
	"durl/app/share/dao/db"
	"durl/app/share/tool"
	"github.com/beego/beego/v2/core/logs"
)


type getXsrfTokenData struct {
	Token string `json:"token"`
}

// 函数名称: GetXsrfToken
// 功能: 获取XsrfToken
// 输入参数:
// 输出参数:
//		token
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021-11-17 15:15:42 #

func (c *Controller) GetXsrfToken() {
	data := &getXsrfTokenData{
		Token: c.XSRFToken(),
	}
	c.FormatInterfaceResp(comm.OK, comm.OK, comm.MsgOk, data)
	return
}

type setShortUrlReq struct {
	Url            string `form:"url" valid:"Required"`
	ExpirationTime int    `form:"expirationTime"`
}

type setShortUrlDataResp struct {
	Key            string `json:"key"`
	Url            string `json:"url"`
	ExpirationTime int    `json:"expirationTime"`
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
	c.BaseCheckParams(&req)

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

type delShortKeyReq struct {
	Key string `form:"key" valid:"Required"`
}

// 函数名称: DelShortKey
// 功能: 根据单个短链删除短链接
// 输入参数:
//     key: 短链结果
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/18 5:44 下午 #

func (c *Controller) DelShortKey() {

	req := delShortKeyReq{}
	// 效验请求参数格式
	c.BaseCheckParams(&req)

	shortKey := req.Key

	if !tool.DisposeShortKey(req.Key) {
		c.ErrorMessage(comm.ErrParamInvalid,comm.MsgParseFormErr)
		return
	}

	shortNum := tool.Base62Decode(shortKey)

	// 删除此短链
	_, err := db.DelUrlByShortNum(shortNum)
	if err != nil {
		logs.Error("Action DelShortKey, err: ", err.Error())
		c.ErrorMessage(comm.ErrSysDb,comm.MsgNotOk)
		return
	}
	c.FormatResp(comm.OK,comm.OK,comm.MsgOk)
	return
}

type updateShortUrlReq struct {
	Key            string `form:"key"  valid:"Required"`
	Url            string `form:"url"`
	IsFrozen       int    `form:"isFrozen"`
	ExpirationTime int64  `form:"expirationTime"`
}

//type updateShortUrlResp struct {
//	Code int    `json:"code"`
//	Msg  string `json:"msg"`
//}

// 函数名称: UpdateShortUrl
// 功能: 根据短链修改短链接信息
// 输入参数:
//     key: 短链内容
//	   url: 原始url
//	   isFrozen: 是否冻结
//	   expirationTime: 过期时间
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/18 5:46 下午 #

func (c *Controller) UpdateShortUrl() {

	req := updateShortUrlReq{}
	// 效验请求参数格式
	c.BaseCheckParams(&req)

	// 效验 key
	shortKey := req.Key
	if !tool.DisposeShortKey(req.Key) {
		c.ErrorMessage(comm.ErrParamInvalid,comm.MsgParseFormErr)
		return
	}
	shortNum := tool.Base62Decode(shortKey)

	// 初始化需要更新的内容
	updateData := make(map[string]interface{})

	if req.ExpirationTime !=0 {
		updateData["expirationTime"] = req.ExpirationTime
	}

	if req.Url !="" {
		updateData["url"] = req.Url
	}

	updateData["isFrozen"] = req.IsFrozen

	if len(updateData) == 0 {
		c.ErrorMessage(comm.ErrParamMiss,comm.MsgParseFormErr)
		return
	}

	// 修改此短链信息
	_, err := db.UpdateUrlByShortNum(shortNum, &updateData)
	if err != nil {
		c.ErrorMessage(comm.ErrSysDb,comm.MsgNotOk)
		return
	}

	c.FormatResp(comm.OK,comm.OK,comm.MsgOk)
	return
}


type getShortUrlListReq struct {
	Url  string `form:"url"`
	Page int    `form:"page"`
	Size int    `form:"size"`
}

type getShortUrlListDataResp struct {
	Id             int    `json:"id"`
	ShortNum       int    `json:"shor_url"`
	FullUrl        string `json:"full_url"`
	ExpirationTime int    `json:"expiration_time"`
	IsFrozen       int8   `json:"is_frozen"`
	CreateTime     int    `json:"create_time"`
}

// 函数名称: GetShortUrlList
// 功能: 分页获取url数据
// 输入参数:
//   	url: 原始url
//		page: 页码  默认0
//		size: 每页展示条数 默认 20  最大500
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/18 6:41 下午 #

func (c *Controller) GetShortUrlList() {
	req := getShortUrlListReq{}
	// 效验请求参数格式
	c.BaseCheckParams(&req)

	// page 需要给默认值
	if req.Page == 0 {
		req.Page = 1
	}
	// size 需要默认值且不可大于500
	if req.Size == 0 || req.Size > 500 {
		req.Size = 20
	}

	where := " is_del =? "
	bindValue := []interface{}{0}
	// key 是原url模糊搜索
	// 拼接条件进行检索
	if req.Url != "" {
		where += " and full_url like ? "
		bindValue = append(bindValue, "%"+req.Url+"%")

	}
	data := db.GetShortUrlList(where, req.Page, req.Size, bindValue...)
	var total int64
	// 有数据且当page=1时计算结果总条数
	if data != nil && req.Page == 1 {
		// 统计结果总条数
		total = db.GetShortUrlListTotal(where, bindValue...)
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
	c.FormatInterfaceListResp(comm.OK,comm.OK,total,comm.MsgOk,list)
	return

}
