package dbstruct

import (
	"github.com/xormplus/xorm"
)

type ShortNumStruct struct {
	Id         int `xorm:" int pk notnull autoincr"`
	MaxNum     int `xorm:" int notnull default(100)"`
	Step       int `xorm:" int notnull default(100)"`
	Version    int `xorm:" version int notnull default(1)"`
	UpdateTime int `xorm:" updated int notnull default(0)"`
}

func (I *ShortNumStruct) TableName() string {
	return "durl_short_num"
}

const (
	ShortNumIsDelYes = 1
	ShortNumIsDelNo  = 0
)

// 函数名称: ReturnShortNumPeriod
// 功能: 获取号码段
// 输入参数:
//		id
// 输出参数:
//		Step: 步长
//		MaxNum: 号码段开始值
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 20:44 下午 #

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

// 函数名称: InsertFirst
// 功能: 插入号码段第一条默认数据
// 输入参数:
//		id
// 输出参数:
//		Step: 步长
//		MaxNum: 号码段开始值
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 20:44 下午 #

func InsertFirst(engine *xorm.EngineGroup) error {
	var shortNumDetail ShortNumStruct
	shortNumDetail.Id = 1
	shortNumDetail.MaxNum = 100
	shortNumDetail.Step = 100
	_, err := engine.Insert(shortNumDetail)
	return err
}
