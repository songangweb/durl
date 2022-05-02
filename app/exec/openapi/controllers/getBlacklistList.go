package controllers

import (
	"durl/app/share/comm"
	"durl/app/share/dao/db"
)

type getBlacklistListReq struct {
	Ip          string `form:"ip"`
	Page        int    `form:"page" valid:"Min(1)"`
	Size        int    `form:"size" valid:"Range(1,500)"`
	CreateTimeL int    `form:"createTimeL"`
	CreateTimeR int    `form:"createTimeR"`
}

// GetBlacklistList
// 函数名称: GetBlacklistList
// 功能: 分页获取url数据
// 输入参数:
//   	shortUrl: 原始url
//		page: 页码  默认1
//		size: 每页展示条数 默认 20  最大500
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #
func (c *OpenApiController) GetBlacklistList() {
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

	engine := db.NewDbService()

	// 统计结果总条数
	total := engine.GetBlacklistListTotal(fields)
	var data []*db.GetBlacklistListRes
	if total != 0 {
		data = engine.GetBlacklistList(fields, req.Page, req.Size)
	}

	c.FormatInterfaceListResp(comm.OK, comm.OK, total, comm.MsgOk, data)
	return

}
