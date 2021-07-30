package xormDbStruct

import (
	"durl/dao/db/xormDb"
)

type ShortNumStruct struct {
	Id         int64 `xorm:" int(11) pk notnull autoincr"`
	MaxNum     int   `xorm:" int(11) notnull default(100)"`
	Step       int   `xorm:" int(11) notnull default(1)"`
	Version    int   `xorm:"version notnull"`
	UpdateTime int   `xorm:"updated notnull default(0)"`
}

func (I *ShortNumStruct) TableName() string {
	return "durl_short_num"
}

// ReturnShortNumPeriod 获取号码段
func ReturnShortNumPeriod() (int, int, error) {
	var shortNumDetail ShortNumStruct

	// 获取数据
	if has, err := xormDb.Engine.ID(1).Get(&shortNumDetail); nil != err {
		return 0, 0, err
	} else if !has {
		// 插入第一条默认数据
		err := InsertFirst()
		if err !=nil{
			return 0, 0, err
		}
		return shortNumDetail.Step, shortNumDetail.MaxNum, err
	}
	// 修改数据
	shortNumDetail.MaxNum += shortNumDetail.Step
	if affected, err := xormDb.Engine.ID(1).Update(&shortNumDetail); nil != err {
		return 0, 0, err
	} else if 0 == affected {
		return 0, 0, err
	}

	return shortNumDetail.Step, shortNumDetail.MaxNum, nil
}

// 插入第一条默认数据
func InsertFirst() error {
	var shortNumDetail ShortNumStruct
	shortNumDetail.Id = 1
	shortNumDetail.MaxNum = 100
	shortNumDetail.Step = 1
	_, err := xormDb.Engine.Insert(shortNumDetail)
	return err
}