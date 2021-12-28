package dbstruct

import (
	"github.com/xormplus/xorm"
)

type ShortNumStruct struct {
	Id         int8 `xorm:" tinyint pk notnull autoincr"`
	MaxNum     int  `xorm:" int notnull default(100)"`
	Step       int  `xorm:" int notnull default(100)"`
	Version    int  `xorm:" version notnull default(1)"`
	UpdateTime int  `xorm:" updated notnull default(0)"`
}

func (I *ShortNumStruct) TableName() string {
	return "durl_short_num"
}

// ReturnShortNumPeriod 获取号码段
func ReturnShortNumPeriod(engine *xorm.EngineGroup) (int, int, error) {
	var shortNumDetail ShortNumStruct

	// 获取数据
	if has, err := engine.ID(1).Get(&shortNumDetail); nil != err {
		return 0, 0, err
	} else if !has {
		// 插入第一条默认数据
		err := InsertFirst(engine)
		if err != nil {
			return 0, 0, err
		}
		return shortNumDetail.Step, shortNumDetail.MaxNum, err
	}
	// 修改数据
	shortNumDetail.MaxNum += shortNumDetail.Step
	if affected, err := engine.ID(1).Update(&shortNumDetail); nil != err {
		return 0, 0, err
	} else if 0 == affected {
		return 0, 0, err
	}

	return shortNumDetail.Step, shortNumDetail.MaxNum, nil
}

// InsertFirst 插入第一条默认数据
func InsertFirst(engine *xorm.EngineGroup) error {
	var shortNumDetail ShortNumStruct
	shortNumDetail.Id = 1
	shortNumDetail.MaxNum = 100
	shortNumDetail.Step = 100
	_, err := engine.Insert(shortNumDetail)
	return err
}
