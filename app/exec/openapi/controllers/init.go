package controllers

import (
	"container/list"
	"sync"

	"durl/app/share/dao/db"

	"github.com/beego/beego/v2/server/web"
)

type Controller struct {
	web.Controller
}

type Pool struct {
	step    uint32
	numList *list.List
	lock    sync.Mutex
}

var KeyPool *Pool

func InitCon() {
	KeyPool = &Pool{
		step:    0,
		numList: list.New(),
	}
	KeyPool.ProducerKey()
}

// ProducerKey 申请号码段 放入缓存里
func (p *Pool) ProducerKey() {
	// 申请号码段
	Step, MaxNum, _ := db.NewDbService().ReturnShortNumPeriod()
	if Step != 0 && MaxNum != 0 {
		p.lock.Lock()
		defer p.lock.Unlock()

		p.step = Step
		// 放入到短链池中
		var i uint32
		for i = 0; i < Step; i++ {
			p.numList.PushBack(MaxNum - i)
		}
	}
}

// ReturnShortNumOne 获取单个short_num
func ReturnShortNumOne() (ShortNum uint32) {

	KeyPool.lock.Lock()

	ent := KeyPool.numList.Front()
	KeyPool.numList.Remove(ent)

	KeyPool.lock.Unlock()

	ShortNum, _ = ent.Value.(uint32)
	// 判断是否需要申请新的号码段
	if uint32(KeyPool.numList.Len()) < KeyPool.step {
		// 申请号码段 放入缓存里
		KeyPool.ProducerKey()
	}
	return ShortNum
}
