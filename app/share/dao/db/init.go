package db

import (
	"fmt"

	"durl/app/share/comm"
	"durl/app/share/dao/db/struct"

	"github.com/xormplus/xorm"
)

var Engine *xorm.EngineGroup

func InitXormDb(c XormConf) {
	switch c.Type {
	case "mysql":
		InitMysql(c.Mysql)
	default:
		defer fmt.Println(comm.MsgCheckDbType)
		panic(comm.MsgDbTypeError + ", type: " + c.Type)
	}
}

// CheckMysqlTable 检查Mysql表配置
func CheckMysqlTable() {
	// 获取数据表信息
	tables := make(map[string]interface{}, 4)
	NewQueue := dbstruct.QueueStruct{}
	tables[NewQueue.TableName()] = NewQueue

	NewShortNum := dbstruct.ShortNumStruct{}
	tables[NewShortNum.TableName()] = NewShortNum

	NewUrl := dbstruct.UrlStruct{}
	tables[NewUrl.TableName()] = NewUrl

	NewBlacklist := dbstruct.BlacklistStruct{}
	tables[NewUrl.TableName()] = NewBlacklist

	for tableName, tableStruct := range tables {
		// 判断表是否已经存在, 如果已经存在则不自动创建
		res, err := Engine.IsTableExist(tableName)
		if err != nil {
			defer fmt.Println(comm.MsgCheckDbMysqlConf)
			panic(tableName + " " + comm.MsgInitDbMysqlTable + ", errMsg: " + err.Error())
		}

		if !res {
			// 同步表结构
			err = Engine.Charset("utf8mb4").StoreEngine("InnoDB").Sync2(tableStruct)
			if err != nil {
				defer fmt.Println(comm.MsgCheckDbMysqlConf)
				panic(tableName + " " + comm.MsgInitDbMysqlTable + ", errMsg: " + err.Error())
			}

			if tableName == NewShortNum.TableName() {
				// 添加短链号码段表数据
				has, err := Engine.ID(1).Exist(&dbstruct.ShortNumStruct{})
				if err != nil {
					defer fmt.Println(comm.MsgCheckDbMysqlConf)
					panic(tableName + " " + comm.MsgCheckDbMysqlConf + ", errMsg: " + err.Error())
				}
				if !has {
					err := dbstruct.InsertFirst(Engine)
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
