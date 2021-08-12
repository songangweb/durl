package xormDbStruct

import (
	"durl/app/share/dao/db/xormDb"
)

type QueueStruct struct {
	Id         int  `xorm:"int pk notnull autoincr"`
	ShortNum   int  `xorm:"int notnull"`
	IsDel      int8 `xorm:"tinyint notnull default(0)"`
	CreateTime int  `xorm:"created int notnull default(0) "`
	UpdateTime int  `xorm:"updated int notnull default(0)"`
}

func (I *QueueStruct) TableName() string {
	return "durl_queue"
}

// InsertQueueOne 插入一条数据
func InsertQueueOne(req QueueStruct) (int64, error) {
	Detail := new(QueueStruct)
	Detail.ShortNum = req.ShortNum
	affected, err := xormDb.Engine.Insert(Detail)
	return affected, err
}

// ReturnQueueLastId 获取最新一条数据的id
func ReturnQueueLastId() (int, error) {
	QueueDetail := new(QueueStruct)
	_, err := xormDb.Engine.Desc("id").Get(QueueDetail)
	return QueueDetail.Id, err
}

// GetQueueListById 获取需要处理的数据
func GetQueueListById(id interface{}) ([]*QueueStruct, error) {
	pEveryOne := make([]*QueueStruct, 0)
	err := xormDb.Engine.Where("id > ? and is_del = ?", id, 0).Find(&pEveryOne)
	return pEveryOne, err
}

func check() {

}
