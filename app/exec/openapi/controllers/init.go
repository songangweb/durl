package controllers

import (
	"container/list"
	"durl/app/share/dao/db"
	"github.com/beego/beego/v2/server/web"
	"sync"
)

type Controller struct {
	web.Controller
}

type Pool struct {
	step    int
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
	Step, MaxNum, _ := db.NewDbService(db.Engine).ReturnShortNumPeriod()
	if Step != 0 && MaxNum != 0 {
		p.lock.Lock()
		defer p.lock.Unlock()

		p.step = Step
		// 放入到短链池中
		for i := 0; i < Step; i++ {
			p.numList.PushBack(MaxNum - i)
		}
	}
}

// ReturnShortNumOne 获取单个short_num
func ReturnShortNumOne() (ShortNum int) {

	KeyPool.lock.Lock()

	ent := KeyPool.numList.Front()
	KeyPool.numList.Remove(ent)

	KeyPool.lock.Unlock()

	ShortNum, _ = ent.Value.(int)
	// 判断是否需要申请新的号码段
	if KeyPool.numList.Len() < KeyPool.step {
		// 申请号码段 放入缓存里
		KeyPool.ProducerKey()
	}
	return ShortNum
}
