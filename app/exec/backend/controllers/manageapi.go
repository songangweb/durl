package controllers

import (
	comm "durl/app/share/comm"
	"durl/app/share/dao/db"
	"durl/app/share/tool"
	"github.com/beego/beego/v2/core/logs"
	"reflect"
	"strconv"
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
	ExpirationTime int    `form:"expirationTime" valid:"Match(/([0-9]{10}$)|([0])/);Max(9999999999)"`
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

// 函数名称: DelShortUrl
// 功能: 根据单个短链删除短链接
// 输入参数:
//     key: 短链结果
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/18 5:44 下午 #

func (c *Controller) DelShortUrl() {

	id := c.Ctx.Input.Param(":id")

	// 查询此短链
	var where map[string][]interface{}
	where = make(map[string][]interface{})
	where["id"] = append(where["id"], "=", id)
	where["is_del"] = append(where["is_del"], "=", 0)
	urlInfo := db.GetShortUrlInfo(where)
	if urlInfo.ShortNum == 0 {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}

	// 删除此短链
	_, err := db.DelUrlById(id, urlInfo.ShortNum)
	if err != nil {
		logs.Error("Action DelShortKey, err: ", err.Error())
		c.ErrorMessage(comm.ErrSysDb, comm.MsgNotOk)
		return
	}
	c.FormatResp(comm.OK, comm.OK, comm.MsgOk)
	return
}

type updateShortUrlReq struct {
	Url            string `form:"url" valid:"Required"`
	IsFrozen       int    `form:"isFrozen" valid:"Range(0,1)"`
	ExpirationTime int    `form:"expirationTime" valid:"Match(/([0-9]{10}$)|([0])/);Max(9999999999)"`
}

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

	id := c.Ctx.Input.Param(":id")

	// 查询此短链
	var where map[string][]interface{}
	where = make(map[string][]interface{})
	where["id"] = append(where["id"], "=", id)
	where["is_del"] = append(where["is_del"], "=", 0)
	urlInfo := db.GetShortUrlInfo(where)
	if urlInfo.ShortNum == 0 {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}

	// 初始化需要更新的内容
	updateData := make(map[string]interface{})
	updateData["expiration_time"] = req.ExpirationTime
	updateData["full_url"] = req.Url
	updateData["is_frozen"] = req.IsFrozen

	// 修改此短链信息
	_, err := db.UpdateUrlById(id, urlInfo.ShortNum, updateData)
	if err != nil {
		c.ErrorMessage(comm.ErrSysDb, comm.MsgNotOk)
		return
	}

	c.FormatResp(comm.OK, comm.OK, comm.MsgOk)
	return
}

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
	var where map[string][]interface{}
	where = make(map[string][]interface{})
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

// 函数名称: GetShortUrlInfo
// 功能: 获取url详情
// 输入参数:
//     id: 短链id
// 输出参数:
// 返回: 短链详情
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/26 10:59 上午 #

func (c *Controller) GetShortUrlInfo() {

	id := c.Ctx.Input.Param(":id")

	// 查询此短链
	var where map[string][]interface{}
	where = make(map[string][]interface{})
	where["id"] = append(where["id"], "=", id)
	where["is_del"] = append(where["is_del"], "=", 0)
	urlInfo := db.GetShortUrlInfo(where)
	if urlInfo.ShortNum == 0 {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}
	c.FormatInterfaceResp(comm.OK, comm.OK, comm.MsgOk, urlInfo)
	return
}

// 函数名称: FrozenShortUrl
// 功能: 冻结ShortUrl
// 输入参数:
//     id: 短链id
// 输出参数:
// 返回: 冻结操作结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/26 1:56 下午 #

func (c *Controller) FrozenShortUrl() {

	id := c.Ctx.Input.Param(":id")

	// 查询此短链
	var where map[string][]interface{}
	where = make(map[string][]interface{})
	where["id"] = append(where["id"], "=", id)
	where["is_del"] = append(where["is_del"], "=", 0)
	urlInfo := db.GetShortUrlInfo(where)
	if urlInfo.ShortNum == 0 {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}

	// 冻结/解冻ShortUrl
	updateData := make(map[string]interface{})
	if urlInfo.IsFrozen == 0 {
		updateData["is_frozen"] = 1
	} else {
		updateData["is_frozen"] = 0
	}

	_, err := db.UpdateUrlById(id, urlInfo.ShortNum, updateData)
	if err != nil {
		c.ErrorMessage(comm.ErrSysDb, comm.MsgNotOk)
		return
	}
	c.FormatResp(comm.OK, comm.OK, comm.MsgOk)
	return
}

type BatchFrozenShortUrlReq struct {
	Ids      []string `from:"ids" valid:"Required"`
	IsFrozen int      `from:"isFrozen" valid:"Range(0,1)"`
}

type BatchFrozenShortUrlRes struct {
	RequestCount int
	UpdateCount  int
	ErrIds       []string
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

func (c *Controller) BatchFrozenShortUrl() {

	req := BatchFrozenShortUrlReq{}

	c.BaseCheckParams(&req)

	// 查询待操作Url信息
	var where map[string][]interface{}
	where = make(map[string][]interface{})
	where["id"] = append(where["id"], "in", req.Ids)
	where["is_del"] = append(where["is_del"], "=", 0)
	data := db.GetAllShortUrl(where)

	if data == nil {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}

	var updateIds []string
	var errIds []string
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

	// 正确数据进行批量操作
	// 批量冻结/解冻Url
	updateData := map[string]interface{}{"is_frozen": req.IsFrozen}
	var updateWhere map[string][]interface{}
	updateWhere = make(map[string][]interface{})
	updateWhere["id"] = append(updateWhere["id"], "in", updateIds)

	_, err := db.BatchUpdateUrlByIds(updateWhere, insertShortNum, updateData)
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

type BatchDelShortUrlReq struct {
	Ids []string `from:"ids" valid:"Required"`
}

type BatchDelShortUrlRes struct {
	RequestCount int
	DelCount     int
	ErrIds       []string
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

func (c *Controller) BatchDelShortUrl() {

	req := BatchDelShortUrlReq{}

	c.BaseCheckParams(&req)

	// 查询待操作Url信息
	var where map[string][]interface{}
	where = make(map[string][]interface{})
	where["id"] = append(where["id"], "in", req.Ids)
	where["is_del"] = append(where["is_del"], "=", 0)
	data := db.GetAllShortUrl(where)
	if data == nil {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}

	var updateIds []string
	var errIds []string
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
	updateData := map[string]interface{}{"is_del": 1}
	var updateWhere map[string][]interface{}
	updateWhere = make(map[string][]interface{})
	updateWhere["id"] = append(updateWhere["id"], "in", updateIds)

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
