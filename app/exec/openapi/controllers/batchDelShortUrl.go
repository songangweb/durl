package controllers

import (
	"durl/app/share/comm"
	"durl/app/share/dao/db"
	dbstruct "durl/app/share/dao/db/struct"
)

type BatchDelShortUrlReq struct {
	Ids []int `from:"ids" valid:"Required"`
}

type BatchDelShortUrlRes struct {
	RequestCount int   `json:"requestCount"`
	DelCount     int   `json:"delCount"`
	ErrIds       []int `json:"errIds"`
}

// BatchDelShortUrl
// 函数名称: BatchDelShortUrl
// 功能: 批量删除ShortUrl
// 输入参数:
//     BatchDelShortUrlReq struct
// 输出参数:
//	   BatchDelShortUrlRes struct
// 返回: 操作结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/1 1:41 下午 #
func (c *OpenApiController) BatchDelShortUrl() {

	req := BatchDelShortUrlReq{}

	c.BaseCheckParams(&req)

	// 查询待操作Url信息
	fields := map[string]interface{}{"id": req.Ids}
	engine := db.NewDbService()
	data := engine.GetAllShortUrl(fields)
	if data == nil {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}

	var updateIds []int
	errIds := make([]int, 0)
	var insertShortNum []int
	// 提交id数量与查询出的数据量不一致
	// 需要以数据库数据为准筛选出差集，准备进行错误返回
	requestCount := len(req.Ids)
	updateCount := len(data)
	if updateCount != requestCount {

		// 将请求操作的id 提为key
		mapData := make(map[int]interface{})
		for _, v := range data {
			mapData[v.Id] = v.ShortNum
		}

		for _, v := range req.Ids {
			if mapData[v] != nil {
				updateIds = append(updateIds, v)
				insertShortNum = append(insertShortNum, mapData[v].(int))
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

	// 进行删除操作
	updateData := map[string]interface{}{"is_del": dbstruct.UrlIsDelYes}
	updateWhere := map[string]interface{}{"id": updateIds}
	updateWhere["id"] = updateIds

	_, err := engine.BatchUpdateUrlByIds(updateWhere, insertShortNum, updateData)
	if err != nil {
		c.FormatInterfaceResp(comm.OK, comm.OK, comm.MsgOk, &BatchDelShortUrlRes{
			RequestCount: requestCount,
			DelCount:     0,
			ErrIds:       req.Ids,
		})
		return
	}
	res := BatchDelShortUrlRes{
		RequestCount: requestCount,
		DelCount:     updateCount,
		ErrIds:       errIds,
	}
	c.FormatInterfaceResp(comm.OK, comm.OK, comm.MsgOk, &res)
	return
}
