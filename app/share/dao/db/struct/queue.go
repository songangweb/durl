package dbstruct

import (
	"github.com/xormplus/xorm"
)

type QueueStruct struct {
	Id         uint32 `xorm:"int pk notnull autoincr"`
	ShortNum   uint32 `xorm:"int notnull"`
	IsDel      uint8  `xorm:"int notnull default(0)"`
	CreateTime uint32 `xorm:"created int notnull default(0) "`
	UpdateTime uint32 `xorm:"updated int notnull default(0)"`
}

func (I *QueueStruct) TableName() string {
	return "durl_queue"
}

const (
	QueueIsDelYes uint8 = 1
	QueueIsDelNo  uint8 = 0
)

// InsertQueueOne 插入一条数据
func InsertQueueOne(engine *xorm.EngineGroup, req QueueStruct) (uint32, error) {
	Detail := new(QueueStruct)
	Detail.ShortNum = req.ShortNum
	affected, err := engine.Insert(Detail)
	return uint32(affected), err
}

// ReturnQueueLastId 获取最新一条数据的id
func ReturnQueueLastId(engine *xorm.EngineGroup) (uint32, error) {
	QueueDetail := new(QueueStruct)
	_, err := engine.Desc("id").Get(QueueDetail)
	return QueueDetail.Id, err
}

// GetQueueListById 获取需要处理的数据
func GetQueueListById(engine *xorm.EngineGroup, id uint32) ([]*QueueStruct, error) {
	pEveryOne := make([]*QueueStruct, 0)
	err := engine.Where("id > ? and is_del = ?", id, QueueIsDelNo).Find(&pEveryOne)
	return pEveryOne, err
}
