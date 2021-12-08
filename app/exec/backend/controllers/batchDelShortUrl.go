package controllers

import (
	comm "durl/app/share/comm"
	"durl/app/share/dao/db"
	"reflect"
	"strconv"
)

type BatchDelShortUrlReq struct {
	Ids []string `from:"ids" valid:"Required"`
}

type BatchDelShortUrlRes struct {
	RequestCount int      `json:"requestCount"`
	DelCount     int      `json:"delCount"`
	ErrIds       []string `json:"errIds"`
}

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

func (c *BackendController) BatchDelShortUrl() {

	req := BatchDelShortUrlReq{}

	c.BaseCheckParams(&req)

	// 查询待操作Url信息
	fields := map[string]interface{}{"id": req.Ids}
	data := db.GetAllShortUrl(fields)
	if data == nil {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}

	var updateIds []string
	errIds := make([]string, 0)
	var insertShortNum []int
	// 提交id数量与查询出的数据量不一致
	// 需要以数据库数据为准筛选出差集，准备进行错误返回
	requestCount := len(req.Ids)
	updateCount := len(data)
	if updateCount != requestCount {

		// 将请求操作的id 提为key
		mapData := make(map[string]interface{})
		if vType := reflect.TypeOf(data[0].Id); vType.Name() == "int" {
			for _, v := range data {
				mapData[strconv.Itoa(v.Id.(int))] = v.ShortNum
			}
		} else {
			for _, v := range data {
				mapData[v.Id.(string)] = v.ShortNum
			}
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
	updateData := map[string]interface{}{"is_del": comm.TureDel}
	updateWhere := map[string]interface{}{"id": updateIds}
	updateWhere["id"] = updateIds

	_, err := db.BatchUpdateUrlByIds(updateWhere, insertShortNum, updateData)
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
