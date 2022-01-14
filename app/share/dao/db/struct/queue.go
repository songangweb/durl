package dbstruct

import (
	"github.com/xormplus/xorm"
)

type QueueStruct struct {
	Id         int `xorm:"int pk notnull autoincr"`
	ShortNum   int `xorm:"int notnull"`
	IsDel      int `xorm:"int notnull default(0)"`
	CreateTime int `xorm:"created int notnull default(0) "`
	UpdateTime int `xorm:"updated int notnull default(0)"`
}

func (I *QueueStruct) TableName() string {
	return "durl_queue"
}

const (
	QueueIsDelYes = 1
	QueueIsDelNo  = 0
)

// 函数名称: InsertQueueOne
// 功能: 插入一条数据
// 输入参数:
//		req: QueueStruct
// 输出参数:
// 返回: bool: 操作结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 20:44 下午 #

func InsertQueueOne(engine *xorm.EngineGroup, req *QueueStruct) (int, error) {
	Detail := new(QueueStruct)
	Detail.ShortNum = req.ShortNum
	affected, err := engine.Insert(Detail)
	return int(affected), err
}

// 函数名称: ReturnQueueLastId
// 功能: 获取最新一条数据的id
// 输入参数:
// 输出参数:
//		id
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 20:44 下午 #

func ReturnQueueLastId(engine *xorm.EngineGroup) (int, error) {
	QueueDetail := new(QueueStruct)
	_, err := engine.Desc("id").Get(QueueDetail)
	return QueueDetail.Id, err
}

// 函数名称: GetQueueListById
// 功能: 获取需要处理的数据
// 输入参数:
//		id
// 输出参数:
//		[]QueueStruct
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 20:44 下午 #

func GetQueueListById(engine *xorm.EngineGroup, id int) ([]*QueueStruct, error) {
	pEveryOne := make([]*QueueStruct, 0)
	err := engine.Where("id > ? and is_del = ?", id, QueueIsDelNo).Find(&pEveryOne)
	return pEveryOne, err
}
