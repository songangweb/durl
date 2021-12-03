package controllers

import (
	comm "durl/app/share/comm"
	"durl/app/share/dao/db"
	"durl/app/share/tool"
)

type getShortUrlListReq struct {
	Url       string `form:"url"`
	Page      int    `form:"page" valid:"Min(1)"`
	Size      int    `form:"size" valid:"Range(1,500)"`
	StartTime int    `from:"StartTime" valid:"Match(/([0-9]{10}$)|([0])/);Range(0,9999999999)"`
	EndTime   int    `from:"EndTime" valid:"Match(/([0-9]{10}$)|([0])/);Range(0,9999999999)"`
}

type getShortUrlListDataResp struct {
	Id             interface{} `json:"id"`
	ShortNum       string      `json:"shorUrl"`
	FullUrl        string      `json:"fullUrl"`
	ExpirationTime int         `json:"expirationTime"`
	IsFrozen       int8        `json:"isFrozen"`
	CreateTime     int         `json:"createTime"`
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

	// 透传业务搜索字段
	fields := make(map[string]interface{})
	if req.Url != "" {
		fields["fullUrl"] = req.Url
	}
	if req.StartTime != 0 {
		fields["startTime"] = req.StartTime
	}
	if req.EndTime != 0 {
		fields["endTime"] = req.EndTime
	}

	data := db.GetShortUrlList(fields, req.Page, req.Size)

	var total int64
	// 有数据且当page=1时计算结果总条数
	if data != nil && req.Page == 1 {
		// 统计结果总条数
		total = db.GetShortUrlListTotal(fields)
	}

	var list []*getShortUrlListDataResp
	for _, queueStruct := range data {
		var One getShortUrlListDataResp
		One.Id = queueStruct.Id
		One.ShortNum = tool.Base62Encode(queueStruct.ShortNum)
		One.FullUrl = queueStruct.FullUrl
		One.ExpirationTime = queueStruct.ExpirationTime
		One.IsFrozen = queueStruct.IsFrozen
		One.CreateTime = queueStruct.CreateTime
		list = append(list, &One)
	}
	c.FormatInterfaceListResp(comm.OK, comm.OK, total, comm.MsgOk, list)
	return

}
