package controllers

import (
	"durl/app/share/comm"
	"durl/app/share/dao/db"
)

type BatchFrozenShortUrlReq struct {
	Ids      []uint32 `from:"ids" valid:"Required"`
	IsFrozen uint8   `from:"isFrozen" valid:"Range(0,1)"`
}

type BatchFrozenShortUrlRes struct {
	RequestCount uint32   `json:"requestCount"`
	UpdateCount  uint32   `json:"updateCount"`
	ErrIds       []uint32 `json:"errIds"`
}

// 函数名称: BatchFrozenShortUrl
// 功能: 批量冻结/解冻Url
// 输入参数:
//		BatchFrozenShortUrlReq{}
// 输出参数:
// 返回: BatchFrozenShortUrlRes{}
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/26 2:15 下午 #

func (c *BackendController) BatchFrozenShortUrl() {

	req := BatchFrozenShortUrlReq{}

	c.BaseCheckParams(&req)

	// 查询待操作Url信息
	fields := map[string]interface{}{"id": req.Ids}
	engine := db.NewDbService()
	data := engine.GetAllShortUrl(fields)
	if data == nil {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}

	var updateIds []uint32
	errIds := make([]uint32, 0)
	var insertShortNum []uint32
	// 提交id数量与查询出的数据量不一致
	// 需要以数据库数据为准筛选出差集，准备进行错误返回
	requestCount := uint32(len(req.Ids))
	updateCount := uint32(len(data))
	if updateCount != requestCount {

		// 将请求操作的id 提为key
		mapData := make(map[uint32]interface{})
		for _, v := range data {
			mapData[v.Id] = v.ShortNum
		}

		for _, v := range req.Ids {
			if mapData[v] != nil {
				updateIds = append(updateIds, v)
				insertShortNum = append(insertShortNum, mapData[v].(uint32))
			} else {
				errIds = append(errIds, v)
			}
		}

	} else {
		updateIds = req.Ids
		for _, vv := range data {
			insertShortNum = append(insertShortNum, vv.ShortNum)
		}
	}

	// 正确数据进行批量操作
	// 批量冻结/解冻Url
	updateData := map[string]interface{}{"is_frozen": req.IsFrozen}
	updateWhere := map[string]interface{}{"id": updateIds}

	_, err := engine.BatchUpdateUrlByIds(updateWhere, insertShortNum, updateData)
	if err != nil {
		c.FormatInterfaceResp(comm.OK, comm.OK, comm.MsgOk, &BatchFrozenShortUrlRes{
			RequestCount: requestCount,
			UpdateCount:  0,
			ErrIds:       req.Ids,
		})
		return
	}
	res := BatchFrozenShortUrlRes{
		RequestCount: requestCount,
		UpdateCount:  updateCount,
		ErrIds:       errIds,
	}
	c.FormatInterfaceResp(comm.OK, comm.OK, comm.MsgOk, &res)
	return
}
