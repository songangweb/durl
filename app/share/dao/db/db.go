package db

import (
	"context"
	_ "durl/app/share/comm"
	comm "durl/app/share/comm"
	"durl/app/share/dao/db/mongoDb"
	mongoDbStruct "durl/app/share/dao/db/mongoDb/struct"
	"durl/app/share/dao/db/xormDb"
	xormDbStruct "durl/app/share/dao/db/xormDb/struct"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"go.mongodb.org/mongo-driver/bson"
)

type Conf struct {
	Type  string
	Xorm  xormDb.Conf
	Mongo mongoDb.Conf
}

var DbType string

func (c Conf) InitDb() {
	c.Xorm.Type = c.Type
	switch c.Type {
	case "mysql":
		DbType = "xorm"
		xormDb.InitXormDb(c.Xorm)

		// 检查数据库表结构是否完善,如不完善则自动创建
		CheckMysqlTable()
	case "mongo":
		DbType = "mongo"
		mongoDb.InitMongoDb(c.Mongo)

		// 检查数据库表结构是否完善,如不完善则自动创建
		CheckMongoTable()
	default:
		defer fmt.Println(comm.MsgCheckDbType)
		panic(comm.MsgDbTypeError + ", type: " + c.Type)
	}
}

// CheckMysqlTable 检查Mysql表配置
func CheckMysqlTable() {
	// 获取数据表信息
	tables := make(map[string]interface{}, 3)
	NewQueue := xormDbStruct.QueueStruct{}
	tables[NewQueue.TableName()] = NewQueue

	NewShortNum := xormDbStruct.ShortNumStruct{}
	tables[NewShortNum.TableName()] = NewShortNum

	NewUrl := xormDbStruct.UrlStruct{}
	tables[NewUrl.TableName()] = NewUrl

	NewBlacklist := xormDbStruct.BlacklistStruct{}
	tables[NewUrl.TableName()] = NewBlacklist

	for tableName, tableStruct := range tables {
		// 判断表是否已经存在, 如果已经存在则不自动创建
		res, err := xormDb.Engine.IsTableExist(tableName)
		if err != nil {
			defer fmt.Println(comm.MsgCheckDbMysqlConf)
			panic(tableName + " " + comm.MsgInitDbMysqlTable + ", errMsg: " + err.Error())
		}

		if !res {
			// 同步表结构
			err = xormDb.Engine.Charset("utf8mb4").StoreEngine("InnoDB").Sync2(tableStruct)
			if err != nil {
				defer fmt.Println(comm.MsgCheckDbMysqlConf)
				panic(tableName + " " + comm.MsgInitDbMysqlTable + ", errMsg: " + err.Error())
			}

			if tableName == NewShortNum.TableName() {
				// 添加短链号码段表数据
				has, err := xormDb.Engine.ID(1).Exist(&xormDbStruct.ShortNumStruct{})
				if err != nil {
					defer fmt.Println(comm.MsgCheckDbMysqlConf)
					panic(tableName + " " + comm.MsgCheckDbMysqlConf + ", errMsg: " + err.Error())
				}
				if !has {
					err := xormDbStruct.InsertFirst()
					if err != nil {
						defer fmt.Println(comm.MsgCheckDbMysqlConf)
						panic(tableName + " " + comm.MsgInitDbMysqlData + ", errMsg: " + err.Error())
					}
				}
			}
		}
	}

	fmt.Println("业务数据表检查完毕!!")
}

// CheckMongoTable 检查Mongo表配置
func CheckMongoTable() {

	// 获取数据表信息
	NewShortNum := mongoDbStruct.ShortNumStruct{}
	filter := bson.D{}
	err := mongoDb.Engine.Collection(NewShortNum.TableName()).FindOne(context.Background(), filter).Decode(&NewShortNum)
	if err != nil {
		err = mongoDbStruct.InsertFirst()
		if err != nil {
			defer fmt.Println(comm.MsgCheckDbMongoConf)
			panic(NewShortNum.TableName() + " " + comm.MsgInitDbMongoData + ", errMsg: " + err.Error())
		}
		fmt.Println("数据表: durl_short_num 初始化数据完毕!!")
	}
	fmt.Println("业务数据表检查完毕!!")
}

// QueueLastId 获取任务最新一条数据的id
func QueueLastId() (id interface{}) {
	if DbType == "xorm" {
		id, _ = xormDbStruct.ReturnQueueLastId()
	} else {
		id, _ = mongoDbStruct.ReturnQueueLastId()
	}
	return id
}

type GetQueueListByIdRe struct {
	Id       interface{} `json:"id"`
	ShortNum int         `json:"shortNum"`
}

// GetQueueListById 获取需要处理的任务数据列表
func GetQueueListById(id interface{}) []*GetQueueListByIdRe {
	var returnList []*GetQueueListByIdRe

	if DbType == "xorm" {
		list, err := xormDbStruct.GetQueueListById(id)
		if err != nil {
			logs.Error("Action xormDbStruct.GetQueueListById, err: ", err.Error())
		} else {
			for _, queueStruct := range list {
				var One GetQueueListByIdRe
				One.Id = queueStruct.Id
				One.ShortNum = queueStruct.ShortNum
				returnList = append(returnList, &One)
			}
		}
	} else {
		list, err := mongoDbStruct.GetQueueListById(id)
		if err != nil {
			logs.Error("Action mongoDbStruct.GetQueueListById, err: ", err.Error())
		} else {
			for _, queueStruct := range list {
				var One GetQueueListByIdRe
				One.Id = queueStruct.Id
				One.ShortNum = queueStruct.ShortNum
				returnList = append(returnList, &One)
			}
		}
	}

	return returnList
}

type GetCacheUrlAllByLimitRe struct {
	ShortNum       int    `json:"shortNum"`
	FullUrl        string `json:"fullUrl"`
	ExpirationTime int    `json:"expirationTime"`
}

// GetCacheUrlAllByLimit 查询出符合条件的全部url
func GetCacheUrlAllByLimit(limit int) []*GetCacheUrlAllByLimitRe {

	var returnList []*GetCacheUrlAllByLimitRe

	if DbType == "xorm" {
		list, err := xormDbStruct.GetCacheUrlAllByLimit(limit)
		if err != nil {
			logs.Error("Action xormDbStruct.GetCacheUrlAllByLimit, err: ", err.Error())
		} else {
			for _, queueStruct := range list {
				var One GetCacheUrlAllByLimitRe
				One.ShortNum = queueStruct.ShortNum
				One.FullUrl = queueStruct.FullUrl
				One.ExpirationTime = queueStruct.ExpirationTime
				returnList = append(returnList, &One)
			}
		}
	} else {
		list, err := mongoDbStruct.GetCacheUrlAllByLimit(limit)
		if err != nil {
			logs.Error("Action mongoDbStruct.GetCacheUrlAllByLimit err :", err.Error())
		} else {
			for _, queueStruct := range list {
				var One GetCacheUrlAllByLimitRe
				One.ShortNum = queueStruct.ShortNum
				One.FullUrl = queueStruct.FullUrl
				One.ExpirationTime = queueStruct.ExpirationTime
				returnList = append(returnList, &One)
			}
		}
	}

	return returnList
}

// ReturnShortNumPeriod 获取号码段
func ReturnShortNumPeriod() (Step int, MaxNum int, err error) {

	if DbType == "xorm" {
		var i int
		for {
			if i >= 10 {
				break
			}
			Step, MaxNum, err = xormDbStruct.ReturnShortNumPeriod()
			if err != nil {
				logs.Error("Action xormDbStruct.ReturnShortNumPeriod, err: ", err.Error())
			} else {
				break
			}
			i++
		}
	} else {
		var i int
		for {
			if i >= 10 {
				break
			}
			Step, MaxNum, err = mongoDbStruct.ReturnShortNumPeriod()
			if err != nil {
				logs.Error("Action mongoDbStruct.ReturnShortNumPeriod, err: ", err.Error())
			} else {
				break
			}
			i++
		}
	}

	return Step, MaxNum, nil
}

type InsertUrlOneReq struct {
	ShortNum       int    `json:"shortNum"`
	FullUrl        string `json:"fullUrl"`
	ExpirationTime int    `json:"expirationTime"`
}

// InsertUrlOne 插入一条数据 shortUrl
func InsertUrlOne(urlStructReq *InsertUrlOneReq) (err error) {

	if DbType == "xorm" {
		var reqOne xormDbStruct.UrlStruct
		reqOne.ShortNum = urlStructReq.ShortNum
		reqOne.FullUrl = urlStructReq.FullUrl
		reqOne.ExpirationTime = urlStructReq.ExpirationTime
		_, err = xormDbStruct.InsertUrlOne(reqOne)
		if err != nil {
			logs.Error("Action xormDbStruct.InsertUrlOne, err: ", err.Error())
		}
	} else {
		var reqOne mongoDbStruct.UrlStruct
		reqOne.ShortNum = urlStructReq.ShortNum
		reqOne.FullUrl = urlStructReq.FullUrl
		reqOne.ExpirationTime = urlStructReq.ExpirationTime
		_, err = mongoDbStruct.InsertUrlOne(reqOne)
		if err != nil {
			logs.Error("Action mongoDbStruct.InsertUrlOne, err: ", err.Error())
		}
	}

	return err
}

// DelUrlByShortNum 通过shortNum删除数据
func DelUrlByShortNum(shortNum int) (reBool bool, err error) {

	if DbType == "xorm" {
		reBool, err = xormDbStruct.DelUrlByShortNum(shortNum)
		if err != nil {
			logs.Error("Action xormDbStruct.DelUrlByShortNum, err: ", err.Error())
		}
	} else {
		reBool, err = mongoDbStruct.DelUrlByShortNum(shortNum)
		if err != nil {
			logs.Error("Action mongoDbStruct.DelUrlByShortNum, err: ", err.Error())
		}
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

func DelUrlById(id string, shortNum int) (reBool bool, err error) {

	if DbType == "xorm" {
		reBool, err = xormDbStruct.DelUrlById(id, shortNum)
		if err != nil {
			logs.Error("Action xormDbStruct.DelUrlById, err: ", err.Error())
		}
	} else {
		reBool, err = mongoDbStruct.DelUrlById(id, shortNum)
		if err != nil {
			logs.Error("Action mongoDbStruct.DelUrlById, err: ", err.Error())
		}
	}

	return reBool, err
}

// UpdateUrlByShortNum 插入一条数据 shortUrl
func UpdateUrlByShortNum(shortNum int, data *map[string]interface{}) (reBool bool, err error) {

	if DbType == "xorm" {
		reBool, err = xormDbStruct.UpdateUrlByShortNum(shortNum, data)
		if err != nil {
			logs.Error("Action xormDbStruct.UpdateUrlByShortNum, err: ", err.Error())
		}
	} else {
		reBool, err = mongoDbStruct.UpdateUrlByShortNum(shortNum, data)
		if err != nil {
			logs.Error("Action mongoDbStruct.UpdateUrlByShortNum, err: ", err.Error())
		}
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

func UpdateUrlById(id string, shortNum int, data map[string]interface{}) (reBool bool, err error) {

	if DbType == "xorm" {
		reBool, err = xormDbStruct.UpdateUrlById(id, shortNum, data)
		if err != nil {
			logs.Error("Action xormDbStruct.UpdateUrlById, err: ", err.Error())
		}
	} else {
		reBool, err = mongoDbStruct.UpdateUrlById(id, shortNum, data)
		if err != nil {
			logs.Error("Action mongoDbStruct.UpdateUrlById, err: ", err.Error())
		}
	}
	return reBool, err
}

type getFullUrlByShortNumReq struct {
	ShortNum       int    `json:"shortNum"`
	FullUrl        string `json:"fullUrl"`
	ExpirationTime int    `json:"expirationTime"`
}

func GetFullUrlByshortNum(shortNum int) *getFullUrlByShortNumReq {

	var One getFullUrlByShortNumReq
	if DbType == "xorm" {
		Detail, err := xormDbStruct.GetFullUrlByShortNum(shortNum)
		if err != nil {
			logs.Error("Action xormDbStruct.GetFullUrlByShortNum, err: ", err.Error())
		}
		if Detail != nil {
			One.ShortNum = Detail.ShortNum
			One.FullUrl = Detail.FullUrl
			One.ExpirationTime = Detail.ExpirationTime
			return &One
		}
	} else {
		Detail, err := mongoDbStruct.GetFullUrlByShortNum(shortNum)
		if err != nil {
			logs.Error("Action mongoDbStruct.GetFullUrlByShortNum, err: ", err.Error())
		}
		if Detail != nil {
			One.ShortNum = Detail.ShortNum
			One.FullUrl = Detail.FullUrl
			One.ExpirationTime = Detail.ExpirationTime
			return &One
		}
	}

	return nil
}

// url列表结构体
type GetShortUrlListRes struct {
	Id             interface{} `json:"id"`
	ShortNum       int         `json:"shortNum"`
	FullUrl        string      `json:"fullUrl"`
	ExpirationTime int         `json:"expirationTime"`
	IsFrozen       int8        `json:"isFrozen"`
	CreateTime     int         `json:"createTime"`
	UpdateTime     int         `json:"updateTime"`
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

func GetShortUrlList(fields map[string]interface{}, page, size int) []*GetShortUrlListRes {

	var returnList []*GetShortUrlListRes

	if DbType == "xorm" {
		list, err := xormDbStruct.GetShortUrlList(fields, page, size)
		if err != nil {
			logs.Error("Action xormDbStruct.GetShortUrlList, err: ", err.Error())
		} else {
			for _, queueStruct := range list {
				var One GetShortUrlListRes
				One.Id = queueStruct.Id
				One.ShortNum = queueStruct.ShortNum
				One.FullUrl = queueStruct.FullUrl
				One.ExpirationTime = queueStruct.ExpirationTime
				One.IsFrozen = queueStruct.IsFrozen
				One.CreateTime = queueStruct.CreateTime
				One.UpdateTime = queueStruct.UpdateTime
				returnList = append(returnList, &One)
			}
		}
	} else {
		list, err := mongoDbStruct.GetShortUrlList(fields, page, size)
		if err != nil {
			logs.Error("Action mongoDbStruct.GetShortUrlList err :", err.Error())
		} else {
			for _, queueStruct := range list {
				var One GetShortUrlListRes
				One.Id = queueStruct.Id.Hex()
				One.ShortNum = queueStruct.ShortNum
				One.FullUrl = queueStruct.FullUrl
				One.ExpirationTime = queueStruct.ExpirationTime
				One.IsFrozen = queueStruct.IsFrozen
				One.CreateTime = int(queueStruct.CreateTime.T)
				One.UpdateTime = int(queueStruct.UpdateTime.T)
				returnList = append(returnList, &One)
			}
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

func GetShortUrlListTotal(fields map[string]interface{}) int64 {

	if DbType == "xorm" {
		total, err := xormDbStruct.GetShortUrlListTotal(fields)
		if err != nil {
			logs.Error("Action xormDbStruct.GetShortUrlListTotal, err: ", err.Error())
		}
		return total
	} else {
		total, err := mongoDbStruct.GetShortUrlListCount(fields)
		if err != nil {
			logs.Error("Action mongoDbStruct.GetShortUrlListTotal err :", err.Error())
		}
		return total
	}

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

func GetShortUrlInfo(fields map[string]interface{}) *GetShortUrlListRes {

	var Info GetShortUrlListRes
	if DbType == "xorm" {
		detail, err := xormDbStruct.GetShortUrlInfo(fields)
		if err != nil {
			logs.Error("Action xormDbStruct.GetShortUrlInfo, err: ", err.Error())
		}
		if detail != nil {
			Info.Id = detail.Id
			Info.ShortNum = detail.ShortNum
			Info.FullUrl = detail.FullUrl
			Info.ExpirationTime = detail.ExpirationTime
			Info.IsFrozen = detail.IsFrozen
			Info.CreateTime = detail.CreateTime
			Info.UpdateTime = detail.UpdateTime
			return &Info
		}
	} else {
		detail, err := mongoDbStruct.GetShortUrlInfo(fields)
		if err != nil {
			logs.Error("Action mongoDbStruct.GetShortUrlInfo, err: ", err.Error())
		}
		if detail != nil {
			Info.Id = detail.Id
			Info.ShortNum = detail.ShortNum
			Info.FullUrl = detail.FullUrl
			Info.ExpirationTime = detail.ExpirationTime
			Info.IsFrozen = detail.IsFrozen
			Info.CreateTime = int(detail.CreateTime.T)
			Info.UpdateTime = int(detail.UpdateTime.T)
			return &Info
		}
	}

	return nil
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

func GetAllShortUrl(fields map[string]interface{}) []*GetShortUrlListRes {

	var returnList []*GetShortUrlListRes

	if DbType == "xorm" {
		list, err := xormDbStruct.GetAllShortUrl(fields)
		if err != nil {
			logs.Error("Action xormDbStruct.GetAllShortUrl, err: ", err.Error())
		} else {
			for _, queueStruct := range list {
				var One GetShortUrlListRes
				One.Id = queueStruct.Id
				One.ShortNum = queueStruct.ShortNum
				One.FullUrl = queueStruct.FullUrl
				One.ExpirationTime = queueStruct.ExpirationTime
				One.IsFrozen = queueStruct.IsFrozen
				One.CreateTime = queueStruct.CreateTime
				returnList = append(returnList, &One)
			}
		}
	} else {
		list, err := mongoDbStruct.GetAllShortUrl(fields)
		if err != nil {
			logs.Error("Action mongoDbStruct.GetAllShortUrl err :", err.Error())
		} else {
			for _, queueStruct := range list {
				var One GetShortUrlListRes
				One.Id = queueStruct.Id.Hex()
				One.ShortNum = queueStruct.ShortNum
				One.FullUrl = queueStruct.FullUrl
				One.ExpirationTime = queueStruct.ExpirationTime
				One.IsFrozen = queueStruct.IsFrozen
				One.CreateTime = int(queueStruct.CreateTime.T)
				returnList = append(returnList, &One)
			}
		}
	}

	return returnList
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

func BatchUpdateUrlByIds(updateWhere map[string]interface{}, insertShortNum []int, updateData map[string]interface{}) (reBool bool, err error) {

	if DbType == "xorm" {
		reBool, err = xormDbStruct.BatchUpdateUrlByIds(updateWhere, insertShortNum, updateData)
		if err != nil {
			logs.Error("Action xormDbStruct.BatchUpdateUrlByIds, err: ", err.Error())
		}
	} else {
		reBool, err = mongoDbStruct.BatchUpdateUrlByIds(updateWhere, insertShortNum, updateData)
		if err != nil {
			logs.Error("Action mongoDbStruct.BatchUpdateUrlByIds, err: ", err.Error())
		}
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

func InsertBlacklistOne(urlStructReq *InsertBlacklistOneReq) (err error) {

	if DbType == "xorm" {
		var reqOne xormDbStruct.BlacklistStruct
		reqOne.Ip = urlStructReq.Ip
		_, err = xormDbStruct.InsertBlacklistOne(reqOne)
		if err != nil {
			logs.Error("Action xormDbStruct.InsertBlacklistOne, err: ", err.Error())
		}
	} else {
		//var reqOne mongoDbStruct.UrlStruct
		//reqOne.ShortNum = urlStructReq.ShortNum
		//reqOne.FullUrl = urlStructReq.FullUrl
		//reqOne.ExpirationTime = urlStructReq.ExpirationTime
		//_, err = mongoDbStruct.InsertUrlOne(reqOne)
		//if err != nil {
		//	logs.Error("Action mongoDbStruct.InsertUrlOne, err: ", err.Error())
		//}
	}

	return err
}

// 黑名单列表结构体
type GetBlacklistListRes struct {
	Id         interface{} `json:"id"`
	Ip         string      `json:"ip"`
	CreateTime int         `json:"createTime"`
	UpdateTime int         `json:"updateTime"`
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

func GetBlacklistInfo(fields map[string]interface{}) *GetBlacklistListRes {

	var Info GetBlacklistListRes
	if DbType == "xorm" {
		detail, err := xormDbStruct.GetBlacklistInfo(fields)
		if err != nil {
			logs.Error("Action xormDbStruct.GetBlacklistInfo, err: ", err.Error())
		}
		if detail != nil {
			Info.Id = detail.Id
			Info.Ip = detail.Ip
			Info.CreateTime = detail.CreateTime
			Info.UpdateTime = detail.UpdateTime
			return &Info
		}
	} else {
		//detail, err := mongoDbStruct.GetShortUrlInfo(fields)
		//if err != nil {
		//	logs.Error("Action mongoDbStruct.GetShortUrlInfo, err: ", err.Error())
		//}
		//if detail != nil {
		//	Info.Id = detail.Id
		//	Info.ShortNum = detail.ShortNum
		//	Info.FullUrl = detail.FullUrl
		//	Info.ExpirationTime = detail.ExpirationTime
		//	Info.IsFrozen = detail.IsFrozen
		//	Info.CreateTime = int(detail.CreateTime.T)
		//	Info.UpdateTime = int(detail.UpdateTime.T)
		//	return &Info
		//}
	}

	return nil
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

func UpdateBlacklistById(id string, data map[string]interface{}) (reBool bool, err error) {

	if DbType == "xorm" {
		reBool, err = xormDbStruct.UpdateBlacklistById(id, data)
		if err != nil {
			logs.Error("Action xormDbStruct.UpdateBlacklistById, err: ", err.Error())
		}
	} else {
		//reBool, err = mongoDbStruct.UpdateUrlById(id, data)
		//if err != nil {
		//	logs.Error("Action mongoDbStruct.UpdateUrlById, err: ", err.Error())
		//}
	}
	return reBool, err
}

// 函数名称: GetBlacklistList
// 功能: 查询url列表数据
// 输入参数:
//     where：sql搜索条件
//     page：页码
//     size：每页展示条数
// 输出参数: []*GetBlacklistListRes
// 返回: 返回结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #

func GetBlacklistList(fields map[string]interface{}, page, size int) []*GetBlacklistListRes {

	var returnList []*GetBlacklistListRes

	if DbType == "xorm" {
		list, err := xormDbStruct.GetBlacklistList(fields, page, size)
		if err != nil {
			logs.Error("Action xormDbStruct.GetShortUrlList, err: ", err.Error())
		} else {
			for _, queueStruct := range list {
				var One GetBlacklistListRes
				One.Id = queueStruct.Id
				One.Ip = queueStruct.Ip
				One.CreateTime = queueStruct.CreateTime
				One.UpdateTime = queueStruct.UpdateTime
				returnList = append(returnList, &One)
			}
		}
	} else {
		//list, err := mongoDbStruct.GetShortUrlList(fields, page, size)
		//if err != nil {
		//	logs.Error("Action mongoDbStruct.GetShortUrlList err :", err.Error())
		//} else {
		//	for _, queueStruct := range list {
		//		var One GetShortUrlListRes
		//		One.Id = queueStruct.Id.Hex()
		//		One.ShortNum = queueStruct.ShortNum
		//		One.FullUrl = queueStruct.FullUrl
		//		One.ExpirationTime = queueStruct.ExpirationTime
		//		One.IsFrozen = queueStruct.IsFrozen
		//		One.CreateTime = int(queueStruct.CreateTime.T)
		//		One.UpdateTime = int(queueStruct.UpdateTime.T)
		//		returnList = append(returnList, &One)
		//	}
		//}
	}

	return returnList
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

func GetBlacklistListTotal(fields map[string]interface{}) int64 {

	if DbType == "xorm" {
		total, err := xormDbStruct.GetBlacklistListTotal(fields)
		if err != nil {
			logs.Error("Action xormDbStruct.GetBlacklistListTotal, err: ", err.Error())
		}
		return total
	} else {
		//total, err := mongoDbStruct.GetShortUrlListCount(fields)
		//if err != nil {
		//	logs.Error("Action mongoDbStruct.GetShortUrlListTotal err :", err.Error())
		//}
		//return total
	}
	return 1
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

func DelBlacklistById(id string) (reBool bool, err error) {

	if DbType == "xorm" {
		reBool, err = xormDbStruct.DelBlacklistById(id)
		if err != nil {
			logs.Error("Action xormDbStruct.DelUrlById, err: ", err.Error())
		}
	} else {
		//reBool, err = mongoDbStruct.DelUrlById(id, shortNum)
		//if err != nil {
		//	logs.Error("Action mongoDbStruct.DelUrlById, err: ", err.Error())
		//}
	}

	return reBool, err
}



// 函数名称: GetBlacklistAll
// 功能: 获取符合条件的所有黑名单列表
// 输入参数:
// 输出参数:
// 返回: 删除结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/12 5:44 下午 #

func GetBlacklistAll() []*GetBlacklistListRes {

	var returnList []*GetBlacklistListRes

	if DbType == "xorm" {
		list, err := xormDbStruct.GetBlacklistAll()
		if err != nil {
			logs.Error("Action xormDbStruct.GetCacheUrlAllByLimit, err: ", err.Error())
		} else {
			for _, queueStruct := range list {
				var One GetBlacklistListRes
				One.Ip = queueStruct.Ip
				returnList = append(returnList, &One)
			}
		}
	} else {
		//list, err := mongoDbStruct.GetCacheUrlAllByLimit(limit)
		//if err != nil {
		//	logs.Error("Action mongoDbStruct.GetCacheUrlAllByLimit err :", err.Error())
		//} else {
		//	for _, queueStruct := range list {
		//		var One GetCacheUrlAllByLimitRe
		//		One.ShortNum = queueStruct.ShortNum
		//		One.FullUrl = queueStruct.FullUrl
		//		One.ExpirationTime = queueStruct.ExpirationTime
		//		returnList = append(returnList, &One)
		//	}
		//}
	}

	return returnList
}
