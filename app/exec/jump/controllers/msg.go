package controllers

import (
	"durl/app/share/comm"
	"durl/app/share/dao/cache"
	"durl/app/share/dao/db"
	"fmt"
	"strconv"
	"time"
)

// 循环获取queue表数据时间 单位:s
const taskQueueTime = 30

// 消息类型
const (
	queueTypeShortNumDel  = 1
	queueTypeBlacklistAdd = 2
	queueTypeBlacklistDel = 3
)

type MSGConf struct {
	Type string
}

// InitMsg
// 函数名称: InitMsg
// 功能: 初始化消息队列
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #
func (c MSGConf) InitMsg() {
	fmt.Println("c.Type: ", c.Type)
	switch c.Type {
	case "mysql":
		InitMysqlMsg()
	default:
		defer fmt.Println(comm.MsgCheckMsgType)
		panic(comm.MsgMsgTypeError + ", type: " + c.Type)
	}
}

// InitMysqlMsg
// 函数名称: InitMysqlMsg
// 功能: 初始化消息队列
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #
func InitMysqlMsg() {
	// 获取任务队列表里最新的一条数据id
	engine := db.NewDbService()
	queueId := engine.QueueLastId()
	for {
		list := engine.GetQueueListById(queueId)
		count := len(list)
		if count > 0 {
			queueId = list[count-1].Id
			for _, val := range list {
				fmt.Println("val: ", val)
				// 有状态更新时，通知所有观察者
				for _, oper := range PurchaseOperFuncArr {
					res, err := oper(val.QueueType, val.Data)
					if err != nil {
						fmt.Println("操作失败")
						break
					}
					if res == false {
						fmt.Println("处理失败")
						break
					}
				}
			}
		}
		time.Sleep(taskQueueTime * time.Second)
	}

}

// PurchaseOperFunc
// 函数名称: PurchaseOperFunc
// 功能: 观察者模式 订阅的消息
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #
type PurchaseOperFunc func(queueType int, data string) (res bool, err error)

// PurchaseOperFuncArr
// 函数名称: PurchaseOperFuncArr
// 功能: 注册的观察者
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #
var PurchaseOperFuncArr = []PurchaseOperFunc{
	shortNumDel,
	blacklistAdd,
	blacklistDel,
}

// shortNumDel
// 函数名称: shortNumDel
// 功能: 短链接缓存删除
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #
func shortNumDel(queueType int, data string) (res bool, err error) {
	if queueType == queueTypeShortNumDel {
		shortNum, _ := strconv.ParseInt(data, 10, 32)
		cache.NewUrlListCache().Gremove(int(shortNum))
	}
	return true, nil
}

// blacklistAdd
// 函数名称: blacklistAdd
// 功能: 黑名单ip添加
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #
func blacklistAdd(queueType int, data string) (res bool, err error) {
	fmt.Println("blacklistAdd: ", blacklistAdd)
	if queueType == queueTypeBlacklistAdd {
		cache.BlacklistConnLock.Lock()
		err = cache.Blacklist.Add(data)
		cache.BlacklistConnLock.Unlock()
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// blacklistDel
// 函数名称: blacklistDel
// 功能: 黑名单ip删除
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #
func blacklistDel(queueType int, data string) (res bool, err error) {
	if queueType == queueTypeBlacklistDel {
		cache.BlacklistConnLock.Lock()
		err = cache.Blacklist.Del(data)
		cache.BlacklistConnLock.Unlock()
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
