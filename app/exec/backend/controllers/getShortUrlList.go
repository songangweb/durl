package controllers

import (
	comm "durl/app/share/comm"
	"durl/app/share/dao/db"
)

type getShortUrlListReq struct {
	Url  string `form:"url"`
	Page int    `form:"page"`
	Size int    `form:"size"`
}

type getShortUrlListDataResp struct {
	Id             interface{} `json:"id"`
	ShortNum       int         `json:"shor_url"`
	FullUrl        string      `json:"full_url"`
	ExpirationTime int         `json:"expiration_time"`
	IsFrozen       int8        `json:"is_frozen"`
	CreateTime     int         `json:"create_time"`
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

	// key 是原url模糊搜索
	where := make(map[string][]interface{})
	where["is_del"] = append(where["is_del"], "=", 0)
	// 拼接条件进行检索
	if req.Url != "" {
		where["full_url"] = append(where["full_url"], "like", req.Url)
	}
	data := db.GetShortUrlList(where, req.Page, req.Size)

	var total int64
	// 有数据且当page=1时计算结果总条数
	if data != nil && req.Page == 1 {
		// 统计结果总条数
		total = db.GetShortUrlListTotal(where)
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
	c.FormatInterfaceListResp(comm.OK, comm.OK, total, comm.MsgOk, list)
	return

}
