package db

import (
	"context"
	"durl/comm"
	"durl/dao/db/mongoDb"
	mongoDbStruct "durl/dao/db/mongoDb/struct"
	"durl/dao/db/xormDb"
	"durl/dao/db/xormDb/struct"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"go.mongodb.org/mongo-driver/bson"
)

type Conf struct {
	Type  string
	Xorm  xormDb.Conf
	Mongo mongoDb.Conf
}

var dbType string

func (c Conf) InitDb() {
	c.Xorm.Type = c.Type
	switch c.Type {
	case "mysql":
		dbType = "xorm"
		xormDb.InitXormDb(c.Xorm)
		CheckMysqlTable()
	case "mongo":
		dbType = "mongo"
		mongoDb.InitMongoDb(c.Mongo)
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
	tables["durl_queue"] = xormDbStruct.QueueStruct{}
	tables["durl_short_num"] = xormDbStruct.ShortNumStruct{}
	tables["durl_url"] = xormDbStruct.UrlStruct{}

	for tableName, v := range tables {
		err := xormDb.Engine.Charset("utf8mb4").StoreEngine("InnoDB").Sync2(v)
		if err != nil {
			defer fmt.Println(comm.MsgCheckDbMysqlConf)
			panic(tableName + comm.MsgCheckDbMysqlTable + ", err: " + fmt.Errorf("%v", err).Error())
		}
		fmt.Println("数据表: " + tableName + " 同步完毕!!")
		if tableName == "durl_short_num" {
			has, err := xormDb.Engine.ID(1).Exist(&xormDbStruct.ShortNumStruct{})
			if err != nil {
				panic(tableName + comm.MsgCheckDbMysqlConf + ", err: " + fmt.Errorf("%v", err).Error())
			}
			if !has {
				err := xormDbStruct.InsertFirst()
				if err != nil {
					defer fmt.Println(comm.MsgCheckDbMysqlConf)
					panic(tableName + comm.MsgCheckDbMysqlData + ", err: " + fmt.Errorf("%v", err).Error())
				}
			}
		}
	}
	fmt.Println("业务数据表检查完毕!!")
}

// 检查Mongo表配置
func CheckMongoTable() {

	// 获取数据表信息
	m := mongoDbStruct.ShortNumStruct{}
	filter := bson.D{}
	err := mongoDb.Engine.Collection("durl_short_num").FindOne(context.Background(), filter).Decode(&m)
	if err != nil {
		err = mongoDbStruct.InsertFirst()
		if err != nil {
			panic(comm.MsgCheckDbMysqlData + ", err: " + fmt.Errorf("%v", err).Error())
		}
		fmt.Println("数据表: durl_short_num 初始化数据完毕!!")
	}
	fmt.Println("业务数据表检查完毕!!")
}

// QueueLastId 获取任务最新一条数据的id
func QueueLastId() (id interface{}) {
	if dbType == "xorm" {
		id, _ = xormDbStruct.ReturnQueueLastId()
	} else {
		id, _ = mongoDbStruct.ReturnQueueLastId()
	}
	return id
}

type GetQueueListByIdRe struct {
	Id       interface{}
	ShortNum int
}

// GetQueueListById 获取需要处理的任务数据列表
func GetQueueListById(id interface{}) []*GetQueueListByIdRe {
	var returnList []*GetQueueListByIdRe

	if dbType == "xorm" {
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
	ShortNum       int
	FullUrl        string
	ExpirationTime int64
}

// GetCacheUrlAllByLimit 查询出符合条件的全部url
func GetCacheUrlAllByLimit(limit int) []*GetCacheUrlAllByLimitRe {

	var returnList []*GetCacheUrlAllByLimitRe

	if dbType == "xorm" {
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

	if dbType == "xorm" {
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
	ShortNum       int
	FullUrl        string
	ExpirationTime int64
}

// InsertUrlOne 插入一条数据 url
func InsertUrlOne(urlStructReq *InsertUrlOneReq) (err error) {

	if dbType == "xorm" {
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

	if dbType == "xorm" {
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

// UpdateUrlByShortNum 插入一条数据 url
func UpdateUrlByShortNum(shortNum int, data *map[string]interface{}) (reBool bool, err error) {

	if dbType == "xorm" {
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

type getFullUrlByShortNumReq struct {
	ShortNum       int
	FullUrl        string
	ExpirationTime int64
}

func GetFullUrlByshortNum(shortNum int) *getFullUrlByShortNumReq {

	var One getFullUrlByShortNumReq
	if dbType == "xorm" {
		Detail, err := xormDbStruct.GetFullUrlByShortNum(shortNum)
		if err != nil {
			logs.Error("Action xormDbStruct.GetFullUrlByShortNum, err: ", err.Error())
		}
		if Detail != nil {
			One.ShortNum = Detail.ShortNum
			One.FullUrl = Detail.FullUrl
			One.ExpirationTime = Detail.ExpirationTime
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
		}
	}
	return &One
}
