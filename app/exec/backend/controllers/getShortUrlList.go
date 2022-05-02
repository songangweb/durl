package controllers

import (
	"durl/app/share/comm"
	"durl/app/share/dao/db"
	"durl/app/share/tool"
)

type getShortUrlListReq struct {
	FullUrl     string `form:"fullUrl"`
	ShortKey    string `form:"shortKey"`
	IsFrozen    int    `form:"isFrozen" valid:"Range(-1,1)"`
	Page        int    `form:"page" valid:"Min(1)"`
	Size        int    `form:"size" valid:"Range(1,500)"`
	CreateTimeL int    `form:"createTimeL" valid:"Match(/([0-9]{10}$)|([0])/);Max(9999999999)"`
	CreateTimeR int    `form:"createTimeR" valid:"Match(/([0-9]{10}$)|([0])/);Max(9999999999)"`
}

type getShortUrlListDataResp struct {
	Id             int    `json:"id"`
	ShortKey       string `json:"shortKey"`
	FullUrl        string `json:"fullUrl"`
	ExpirationTime int    `json:"expirationTime"`
	IsFrozen       int    `json:"isFrozen"`
	CreateTime     int    `json:"createTime"`
	UpdateTime     int    `json:"updateTime"`
}

// GetShortUrlList
// 函数名称: GetShortUrlList
// 功能: 分页获取url数据
// 输入参数:
//   	shortUrl: 原始url
//		page: 页码  默认1
//		size: 每页展示条数 默认 20  最大500
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/18 6:41 下午 #
func (c *BackendController) GetShortUrlList() {
	req := getShortUrlListReq{}
	// 效验请求参数格式
	c.BaseCheckParams(&req)

	// 透传业务搜索字段
	fields := make(map[string]interface{})
	if req.FullUrl != "" {
		fields["fullUrl"] = req.FullUrl
	}
	if req.ShortKey != "" {
		fields["shortKey"] = req.ShortKey
	}
	if req.IsFrozen != 0 {
		if req.IsFrozen == -1 {
			fields["isFrozen"] = 0
		}
		fields["isFrozen"] = req.IsFrozen
	}
	if req.CreateTimeL != 0 {
		fields["createTimeL"] = req.CreateTimeL
	}
	if req.CreateTimeR != 0 {
		fields["createTimeR"] = req.CreateTimeR
	}

	engine := db.NewDbService()

	var total int
	// 统计结果总条数
	total = engine.GetShortUrlListTotal(fields)

	var list []*getShortUrlListDataResp
	if total != 0 {
		data := engine.GetShortUrlList(fields, req.Page, req.Size)
		for _, queueStruct := range data {
			var One getShortUrlListDataResp
			One.Id = queueStruct.Id
			One.ShortKey = tool.Base62Encode(queueStruct.ShortNum)
			One.FullUrl = queueStruct.FullUrl
			One.ExpirationTime = queueStruct.ExpirationTime
			One.IsFrozen = queueStruct.IsFrozen
			One.CreateTime = queueStruct.CreateTime
			One.UpdateTime = queueStruct.UpdateTime
			list = append(list, &One)
		}
	}

	c.FormatInterfaceListResp(comm.OK, comm.OK, total, comm.MsgOk, list)
	return
}
