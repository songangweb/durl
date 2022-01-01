package dbstruct

import (
	"github.com/xormplus/xorm"
)

type ShortNumStruct struct {
	Id         uint8  `xorm:" int pk notnull autoincr"`
	MaxNum     uint32 `xorm:" int notnull default(100)"`
	Step       uint32 `xorm:" int notnull default(100)"`
	Version    uint32 `xorm:" version int notnull default(1)"`
	UpdateTime uint32 `xorm:" updated int notnull default(0)"`
}

func (I *ShortNumStruct) TableName() string {
	return "durl_short_num"
}

const (
	ShortNumIsDelYes uint8 = 1
	ShortNumIsDelNo  uint8 = 0
)

// ReturnShortNumPeriod 获取号码段
func ReturnShortNumPeriod(engine *xorm.EngineGroup) (uint32, uint32, error) {
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
