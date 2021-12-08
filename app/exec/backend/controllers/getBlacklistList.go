package controllers

import (
	comm "durl/app/share/comm"
	"durl/app/share/dao/db"
)

type getBlacklistListReq struct {
	Ip          string `form:"ip"`
	Page        int    `form:"page" valid:"Min(1)"`
	Size        int    `form:"size" valid:"Range(1,500)"`
	CreateTimeL int    `from:"createTimeL" valid:"Match(/([0-9]{10}$)|([0])/);Range(0,9999999999)"`
	CreateTimeR int    `from:"createTimeR" valid:"Match(/([0-9]{10}$)|([0])/);Range(0,9999999999)"`
}

type getBlacklistListDataResp struct {
	Id         interface{} `json:"id"`
	Ip         string      `json:"ip"`
	CreateTime int         `json:"createTime"`
	UpdateTime int         `json:"updateTime"`
}

// 函数名称: GetBlacklistList
// 功能: 分页获取url数据
// 输入参数:
//   	shortUrl: 原始url
//		page: 页码  默认0
//		size: 每页展示条数 默认 20  最大500
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #

func (c *BackendController) GetBlacklistList() {
	req := getBlacklistListReq{}
	// 效验请求参数格式
	c.BaseCheckParams(&req)

	// 透传业务搜索字段
	fields := make(map[string]interface{})
	if req.Ip != "" {
		fields["ip"] = req.Ip
	}
	if req.CreateTimeL != 0 {
		fields["createTimeL"] = req.CreateTimeL
	}
	if req.CreateTimeR != 0 {
		fields["createTimeR"] = req.CreateTimeR
	}

	data := db.GetBlacklistList(fields, req.Page, req.Size)

	var total int64
	// 有数据且当page=1时计算结果总条数
	if data != nil && req.Page == 1 {
		// 统计结果总条数
		total = db.GetBlacklistListTotal(fields)
	}

	var list []*getBlacklistListDataResp
	for _, queueStruct := range data {
		var One getBlacklistListDataResp
		One.Id = queueStruct.Id
		One.Ip = queueStruct.Ip
		One.CreateTime = queueStruct.CreateTime
		One.UpdateTime = queueStruct.UpdateTime
		list = append(list, &One)
	}
	c.FormatInterfaceListResp(comm.OK, comm.OK, total, comm.MsgOk, list)
	return

}
