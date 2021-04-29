package xormDbStruct

import (
	"durl/dao/db/xormDb"
)

type ShortNumStruct struct {
	Id         int64
	MaxNum     int
	Step       int
	Version    int `xorm:"version"`
	UpdateTime int `xorm:"updated"`
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
		return 0, 0, err
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
