package controllers

import (
	"container/list"
	"fmt"
	"sync"

	"durl/app/share/comm"
	"durl/app/share/dao/db"
)

type Controller struct {
	comm.BaseController
}

type Pool struct {
	step    int
	numList *list.List
	lock    sync.Mutex
}

var KeyPool *Pool

// InitCon
// 函数名称: InitCon
// 功能: 初始化号码段
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 5:44 下午 #
func InitCon() {
	KeyPool = &Pool{
		step:    0,
		numList: list.New(),
	}
	KeyPool.ProducerKey()
}

// ProducerKey
// 函数名称: ProducerKey
// 功能: 申请号码段 放入缓存里
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 5:44 下午 #
func (p *Pool) ProducerKey() {
	// 申请号码段
	Step, MaxNum, _ := db.NewDbService().ReturnShortNumPeriod()
	fmt.Println("Step: ", Step)
	fmt.Println("MaxNum: ", MaxNum)
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

// ReturnShortNumOne
// 函数名称: ReturnShortNumOne
// 功能: 获取单个short_num
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 5:44 下午 #
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
