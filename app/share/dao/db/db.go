package db

import (
	"fmt"

	"durl/app/share/comm"
	"durl/app/share/dao/db/struct"

	"github.com/beego/beego/v2/core/logs"
	"github.com/xormplus/xorm"
)

type DbService interface {

	// QueueLastId 获取任务最新一条数据的id
	QueueLastId() (id uint32)
	// GetQueueListById 获取需要处理的任务数据列表
	GetQueueListById(id uint32) []*GetQueueListByIdRe
	// GetCacheUrlAllByLimit 查询出符合条件的全部url
	GetCacheUrlAllByLimit(limit int) []*GetCacheUrlAllByLimitRe
	// ReturnShortNumPeriod 获取号码段
	ReturnShortNumPeriod() (Step uint32, MaxNum uint32, err error)

	// InsertUrlOne 插入一条短链
	InsertUrlOne(urlStructReq *InsertUrlOneReq) (err error)
	// DelUrlByShortNum 通过shortNum删除数据
	DelUrlByShortNum(shortNum uint32) (reBool bool, err error)
	// DelUrlById 通过id删除url数据
	DelUrlById(id uint32, shortNum uint32) (reBool bool, err error)
	// UpdateUrlByShortNum 通过shortUrl修改一条数据
	UpdateUrlByShortNum(shortNum uint32, data *map[string]interface{}) (reBool bool, err error)
	// UpdateUrlById 根据id修改url信息
	UpdateUrlById(id uint32, shortNum uint32, data map[string]interface{}) (reBool bool, err error)
	// GetFullUrlByShortNum 通过 ShortNum 获取 完整连接
	GetFullUrlByShortNum(shortNum uint32) *getFullUrlByShortNumReq
	// GetShortUrlList 获取url列表数据
	GetShortUrlList(fields map[string]interface{}, page, size int) []*GetShortUrlListRes
	// GetShortUrlListTotal 查询url列表数据条数
	GetShortUrlListTotal(fields map[string]interface{}) uint32
	// GetShortUrlInfo 获取ShortUrl详情
	GetShortUrlInfo(fields map[string]interface{}) *GetShortUrlListRes
	// GetAllShortUrl 根据条件获取所有Url信息不带分页
	GetAllShortUrl(fields map[string]interface{}) []*GetShortUrlListRes
	// BatchUpdateUrlByIds 根据UrlId 修改Url信息
	BatchUpdateUrlByIds(updateWhere map[string]interface{}, insertShortNum []uint32, updateData map[string]interface{}) (reBool bool, err error)

	// InsertBlacklistOne 添加黑名单数据
	InsertBlacklistOne(urlStructReq *InsertBlacklistOneReq) (err error)
	// GetBlacklistInfo 获取黑名单详情
	GetBlacklistInfo(fields map[string]interface{}) *GetBlacklistListRes
	// UpdateBlacklistById 根据id修改黑名单信息
	UpdateBlacklistById(id string, data map[string]interface{}) (reBool bool, err error)
	// GetBlacklistList 查询黑名单列表数据
	GetBlacklistList(fields map[string]interface{}, page, size int) []*GetBlacklistListRes
	// GetBlacklistListTotal 查询黑名单列表数据条数
	GetBlacklistListTotal(fields map[string]interface{}) uint32
	// DelBlacklistById 通过id删除黑名单数据
	DelBlacklistById(id string) (reBool bool, err error)
	// GetBlacklistAll 获取符合条件的所有黑名单数据
	GetBlacklistAll() []*GetBlacklistListRes
}

type dbService struct {
	EngineGroup *xorm.EngineGroup
}

func NewDbService() DbService {
	return &dbService{
		EngineGroup: Engine,
	}
}

type DBConf struct {
	Type string
	Xorm XormConf
}

type XormConf struct {
	Type  string
	Mysql MysqlConf
}

func (c DBConf) InitDb() {
	c.Xorm.Type = c.Type
	switch c.Type {
	case "mysql":
		InitXormDb(c.Xorm)
		// 检查数据库表结构是否完善,如不完善则自动创建
		CheckMysqlTable()
	default:
		defer fmt.Println(comm.MsgCheckDbType)
		panic(comm.MsgDbTypeError + ", type: " + c.Type)
	}
}

// QueueLastId 获取任务最新一条数据的id
func (s *dbService) QueueLastId() (id uint32) {
	id, _ = dbstruct.ReturnQueueLastId(s.EngineGroup)
	return id
}

//GetQueueListById 获取需要处理的任务数据列表
func (s *dbService) GetQueueListById(id uint32) []*GetQueueListByIdRe {
	var returnList []*GetQueueListByIdRe

	list, err := dbstruct.GetQueueListById(s.EngineGroup, id)
	if err != nil {
		logs.Error("Action dbStruct.GetQueueListById, err: ", err.Error())
	} else {
		for _, queueStruct := range list {
			var One GetQueueListByIdRe
			One.Id = queueStruct.Id
			One.ShortNum = queueStruct.ShortNum
			returnList = append(returnList, &One)
		}
	}

	return returnList
}

type GetQueueListByIdRe struct {
	Id       uint32 `json:"id"`
	ShortNum uint32 `json:"shortNum"`
}

type GetCacheUrlAllByLimitRe struct {
	ShortNum       uint32 `json:"shortNum"`
	FullUrl        string `json:"fullUrl"`
	ExpirationTime uint32 `json:"expirationTime"`
}

// GetCacheUrlAllByLimit 查询出符合条件的全部url
func (s *dbService) GetCacheUrlAllByLimit(limit int) []*GetCacheUrlAllByLimitRe {

	var returnList []*GetCacheUrlAllByLimitRe

	list, err := dbstruct.GetCacheUrlAllByLimit(s.EngineGroup, limit)
	if err != nil {
		logs.Error("Action dbStruct.GetCacheUrlAllByLimit, err: ", err.Error())
	} else {
		for _, queueStruct := range list {
			var One GetCacheUrlAllByLimitRe
			One.ShortNum = queueStruct.ShortNum
			One.FullUrl = queueStruct.FullUrl
			One.ExpirationTime = queueStruct.ExpirationTime
			returnList = append(returnList, &One)
		}
	}

	return returnList
}

// ReturnShortNumPeriod 获取号码段
func (s *dbService) ReturnShortNumPeriod() (Step uint32, MaxNum uint32, err error) {

	var i int
	for {
		if i >= 10 {
			break
		}
		Step, MaxNum, err = dbstruct.ReturnShortNumPeriod(s.EngineGroup)
		if err != nil {
			logs.Error("Action dbStruct.ReturnShortNumPeriod, err: ", err.Error())
		} else {
			break
		}
		i++
	}

	return Step, MaxNum, nil
}

type InsertUrlOneReq struct {
	ShortNum       uint32 `json:"shortNum"`
	FullUrl        string `json:"fullUrl"`
	ExpirationTime uint32 `json:"expirationTime"`
}

// InsertUrlOne 插入一条数据 shortUrl
func (s *dbService) InsertUrlOne(urlStructReq *InsertUrlOneReq) (err error) {

	var reqOne dbstruct.UrlStruct
	reqOne.ShortNum = urlStructReq.ShortNum
	reqOne.FullUrl = urlStructReq.FullUrl
	reqOne.ExpirationTime = urlStructReq.ExpirationTime
	_, err = dbstruct.InsertUrlOne(s.EngineGroup, reqOne)
	if err != nil {
		logs.Error("Action dbStruct.InsertUrlOne, err: ", err.Error())
	}

	return err
}

// DelUrlByShortNum 通过shortNum删除数据
func (s *dbService) DelUrlByShortNum(shortNum uint32) (reBool bool, err error) {

	reBool, err = dbstruct.DelUrlByShortNum(s.EngineGroup, shortNum)
	if err != nil {
		logs.Error("Action dbStruct.DelUrlByShortNum, err: ", err.Error())
	}

	return reBool, err
}

// 函数名称: DelUrlById
// 功能: 通过id删除url数据
// 输入参数:
//     id: 数据id
//	   shortNum: 短链Key
// 输出参数:
//	   reBool bool
//	   err error
// 返回: 删除结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/24 5:13 下午 #

func (s *dbService) DelUrlById(id uint32, shortNum uint32) (reBool bool, err error) {

	reBool, err = dbstruct.DelUrlById(s.EngineGroup, id, shortNum)
	if err != nil {
		logs.Error("Action dbStruct.DelUrlById, err: ", err.Error())
	}

	return reBool, err
}

// UpdateUrlByShortNum 修改一条数据 shortUrl
func (s *dbService) UpdateUrlByShortNum(shortNum uint32, data *map[string]interface{}) (reBool bool, err error) {

	reBool, err = dbstruct.UpdateUrlByShortNum(s.EngineGroup, shortNum, data)
	if err != nil {
		logs.Error("Action dbStruct.UpdateUrlByShortNum, err: ", err.Error())
	}
	return reBool, err
}

// 函数名称: UpdateUrlById
// 功能: 根据id修改url信息
// 输入参数:
//     id: url数据id
//     shortNum: 短链key值
//     data: 需要修改的信息内容
// 输出参数:
// 返回: 修改操作结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/25 4:53 下午 #

func (s *dbService) UpdateUrlById(id uint32, shortNum uint32, data map[string]interface{}) (reBool bool, err error) {

	reBool, err = dbstruct.UpdateUrlById(s.EngineGroup, id, shortNum, data)
	if err != nil {
		logs.Error("Action dbStruct.UpdateUrlById, err: ", err.Error())
	}
	return reBool, err
}

type getFullUrlByShortNumReq struct {
	ShortNum       uint32 `json:"shortNum"`
	FullUrl        string `json:"fullUrl"`
	ExpirationTime uint32 `json:"expirationTime"`
}

func (s *dbService) GetFullUrlByShortNum(shortNum uint32) *getFullUrlByShortNumReq {

	var One getFullUrlByShortNumReq
	Detail, err := dbstruct.GetFullUrlByShortNum(s.EngineGroup, shortNum)
	if err != nil {
		logs.Error("Action dbStruct.GetFullUrlByShortNum, err: ", err.Error())
	}
	if Detail != nil {
		One.ShortNum = Detail.ShortNum
		One.FullUrl = Detail.FullUrl
		One.ExpirationTime = Detail.ExpirationTime
		return &One
	}

	return nil
}

// GetShortUrlListRes url列表结构体
type GetShortUrlListRes struct {
	Id             uint32 `json:"id"`
	ShortNum       uint32 `json:"shortNum"`
	FullUrl        string `json:"fullUrl"`
	ExpirationTime uint32 `json:"expirationTime"`
	IsFrozen       uint8  `json:"isFrozen"`
	CreateTime     uint32 `json:"createTime"`
	UpdateTime     uint32 `json:"updateTime"`
}

// 函数名称: GetShortUrlList
// 功能: 查询url列表数据
// 输入参数:
//     where：sql搜索条件
//     page：页码
//     size：每页展示条数
// 输出参数: []*GetShortUrlListRes
// 返回: 返回结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/19 3:27 下午 #

func (s *dbService) GetShortUrlList(fields map[string]interface{}, page, size int) []*GetShortUrlListRes {

	var returnList []*GetShortUrlListRes

	list, err := dbstruct.GetShortUrlList(s.EngineGroup, fields, page, size)
	if err != nil {
		logs.Error("Action dbStruct.GetShortUrlList, err: ", err.Error())
	} else {
		for _, queueStruct := range list {
			var one GetShortUrlListRes
			one.Id = queueStruct.Id
			one.ShortNum = queueStruct.ShortNum
			one.FullUrl = queueStruct.FullUrl
			one.ExpirationTime = queueStruct.ExpirationTime
			one.IsFrozen = queueStruct.IsFrozen
			one.CreateTime = queueStruct.CreateTime
			one.UpdateTime = queueStruct.UpdateTime
			returnList = append(returnList, &one)
		}
	}

	return returnList
}

// 函数名称: GetShortUrlListTotal
// 功能: 查询url列表数据条数
// 输入参数:
//     where：业务搜索条件
// 输出参数:
// 返回: 结果条数
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/23 6:21 下午 #

func (s *dbService) GetShortUrlListTotal(fields map[string]interface{}) uint32 {

	total, err := dbstruct.GetShortUrlListTotal(s.EngineGroup, fields)
	if err != nil {
		logs.Error("Action dbStruct.GetShortUrlListTotal, err: ", err.Error())
	}
	return uint32(total)

}

// 函数名称: GetShortUrlInfo
// 功能: 获取ShortUrl详情
// 输入参数:
//     where: 数据检索条件
// 输出参数: *GetShortUrlListRes
// 返回: 检索条件
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/24 5:09 下午 #

func (s *dbService) GetShortUrlInfo(fields map[string]interface{}) *GetShortUrlListRes {

	var returnRes GetShortUrlListRes
	detail, err := dbstruct.GetShortUrlInfo(s.EngineGroup, fields)
	if err != nil {
		logs.Error("Action dbStruct.GetShortUrlInfo, err: ", err.Error())
	}
	if detail.Id != 0 {
		returnRes.Id = detail.Id
		returnRes.ShortNum = detail.ShortNum
		returnRes.FullUrl = detail.FullUrl
		returnRes.IsFrozen = detail.IsFrozen
		returnRes.CreateTime = detail.CreateTime
		returnRes.UpdateTime = detail.UpdateTime
	}
	return &returnRes
}

// 函数名称: GetAllShortUrl
// 功能: 根据条件获取所有Url信息不带分页
// 输入参数:
//     where 检索条件
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/30 6:12 下午 #

func (s *dbService) GetAllShortUrl(fields map[string]interface{}) []*GetShortUrlListRes {

	var returnRes []*GetShortUrlListRes

	list, err := dbstruct.GetAllShortUrl(s.EngineGroup, fields)
	if err != nil {
		logs.Error("Action dbStruct.GetAllShortUrl, err: ", err.Error())
	} else {
		for _, v := range list {
			var One GetShortUrlListRes
			One.Id = v.Id
			One.ShortNum = v.ShortNum
			One.FullUrl = v.FullUrl
			One.IsFrozen = v.IsFrozen
			One.CreateTime = v.CreateTime
			One.UpdateTime = v.UpdateTime
			returnRes = append(returnRes, &One)
		}
	}
	return returnRes
}

// 函数名称: BatchUpdateUrlByIds
// 功能: 根据UrlId 修改Url信息
// 输入参数:
//		updateWhere: 修改限制条件
//	    insertShortNum: 涉及修改Url的短链Key值
//		updateData: 修改内容
// 输出参数:
// 返回: 修改结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/30 6:19 下午 #

func (s *dbService) BatchUpdateUrlByIds(updateWhere map[string]interface{}, insertShortNum []uint32, updateData map[string]interface{}) (reBool bool, err error) {

	reBool, err = dbstruct.BatchUpdateUrlByIds(s.EngineGroup, updateWhere, insertShortNum, updateData)
	if err != nil {
		logs.Error("Action dbStruct.BatchUpdateUrlByIds, err: ", err.Error())
	}
	return reBool, err
}

type InsertBlacklistOneReq struct {
	Ip string `json:"ip"`
}

// 函数名称: InsertBlacklistOne
// 功能: 添加黑名单数据
// 输入参数:
// 返回: 检索条件
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #

func (s *dbService) InsertBlacklistOne(urlStructReq *InsertBlacklistOneReq) (err error) {

	var reqOne dbstruct.BlacklistStruct
	reqOne.Ip = urlStructReq.Ip
	_, err = dbstruct.InsertBlacklistOne(s.EngineGroup, reqOne)
	if err != nil {
		logs.Error("Action dbStruct.InsertBlacklistOne, err: ", err.Error())
	}

	return err
}

// 黑名单列表结构体
type GetBlacklistListRes struct {
	Id         uint32 `json:"id"`
	Ip         string `json:"ip"`
	CreateTime uint32 `json:"createTime"`
	UpdateTime uint32 `json:"updateTime"`
}

// 函数名称: GetBlacklistInfo
// 功能: 获取黑名单详情
// 输入参数:
//     where: 数据检索条件
// 输出参数: *GetBlacklistListRes
// 返回: 检索条件
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #

func (s *dbService) GetBlacklistInfo(fields map[string]interface{}) *GetBlacklistListRes {

	var returnRes GetBlacklistListRes
	detail, err := dbstruct.GetBlacklistInfo(s.EngineGroup, fields)
	if err != nil {
		logs.Error("Action dbStruct.GetBlacklistInfo, err: ", err.Error())
	}
	if detail.Id != 0 {
		returnRes.Id = detail.Id
		returnRes.Ip = detail.Ip
		returnRes.CreateTime = detail.CreateTime
		returnRes.UpdateTime = detail.UpdateTime
	}
	return &returnRes
}

// 函数名称: UpdateBlacklistById
// 功能: 根据id修改黑名单信息
// 输入参数:
//     id: Blacklist数据id
// 输出参数:
// 返回: 修改操作结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #

func (s *dbService) UpdateBlacklistById(id string, data map[string]interface{}) (reBool bool, err error) {

	reBool, err = dbstruct.UpdateBlacklistById(s.EngineGroup, id, data)
	if err != nil {
		logs.Error("Action dbStruct.UpdateBlacklistById, err: ", err.Error())
	}
	return reBool, err
}

// 函数名称: GetBlacklistList
// 功能: 查询黑名单列表数据
// 输入参数:
//     where：sql搜索条件
//     page：页码
//     size：每页展示条数
// 输出参数: []*GetBlacklistListRes
// 返回: 返回结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #

func (s *dbService) GetBlacklistList(fields map[string]interface{}, page, size int) []*GetBlacklistListRes {

	var returnRes []*GetBlacklistListRes
	list, err := dbstruct.GetBlacklistList(s.EngineGroup, fields, page, size)
	if err != nil {
		logs.Error("Action dbStruct.GetShortUrlList, err: ", err.Error())
	} else {
		for _, v := range list {
			var one GetBlacklistListRes
			one.Id = v.Id
			one.Ip = v.Ip
			one.CreateTime = v.CreateTime
			one.UpdateTime = v.UpdateTime
			returnRes = append(returnRes, &one)
		}
	}

	return returnRes
}

// 函数名称: GetBlacklistListTotal
// 功能: 查询黑名单列表数据条数
// 输入参数:
//     where：业务搜索条件
// 输出参数:
// 返回: 结果条数
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #

func (s *dbService) GetBlacklistListTotal(fields map[string]interface{}) uint32 {

	total, err := dbstruct.GetBlacklistListTotal(s.EngineGroup, fields)
	if err != nil {
		logs.Error("Action dbStruct.GetBlacklistListTotal, err: ", err.Error())
	}
	return total
}

// 函数名称: DelBlacklistById
// 功能: 通过id删除黑名单数据
// 输入参数:
//     id: 数据id
// 输出参数:
//	   reBool bool
//	   err error
// 返回: 删除结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #

func (s *dbService) DelBlacklistById(id string) (reBool bool, err error) {

	reBool, err = dbstruct.DelBlacklistById(s.EngineGroup, id)
	if err != nil {
		logs.Error("Action dbStruct.DelUrlById, err: ", err.Error())
	}

	return reBool, err
}

// 函数名称: GetBlacklistAll
// 功能: 获取符合条件的所有黑名单数据
// 输入参数:
// 输出参数:
// 返回: 删除结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/12 5:44 下午 #

func (s *dbService) GetBlacklistAll() []*GetBlacklistListRes {

	var returnList []*GetBlacklistListRes

	list, err := dbstruct.GetBlacklistAll(s.EngineGroup)
	if err != nil {
		logs.Error("Action dbStruct.GetCacheUrlAllByLimit, err: ", err.Error())
	} else {
		for _, queueStruct := range list {
			var One GetBlacklistListRes
			One.Ip = queueStruct.Ip
			returnList = append(returnList, &One)
		}
	}

	return returnList
}
