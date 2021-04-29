package jump

import (
	"durl/comm"
	"durl/dao/db"
	"fmt"
	"github.com/beego/beego/v2/server/web"
	"github.com/songangweb/mcache"
	"time"
)

type Controller struct {
	web.Controller
}

type Pool struct {
	step   int
	keyMap []KeyMapOne
}

type KeyMapOne struct {
	num      int
	shortKey string
}

type Conf struct {
	GoodUrlLen int
	BedUrlLen  int
}

// GoodUrlCache url 内存缓存
var GoodUrlCache *mcache.ARCCache

// BedUrlCache bed url 缓存
var BedUrlCache *mcache.LruCache

func (c Conf) InitJump() {

	var err error

	// 获取任务队列表里最新的一条数据id
	queueId := db.QueueLastId()

	// 初始化Cache数据池
	GoodUrlCache, err = mcache.NewARC(c.GoodUrlLen)
	if err != nil {
		defer fmt.Println(comm.MsgInitializeCacheError)
		panic(comm.MsgInitializeCacheError + ", err: " + err.Error())
	}

	// 获取数据库中需要放到缓存的url
	UrlList := db.GetCacheUrlAllByLimit(c.GoodUrlLen)
	// 添加数据到缓存中
	for i := 0; i < len(UrlList); i++ {
		GoodUrlCache.Add(UrlList[i].ShortNum, UrlList[i].FullUrl, UrlList[i].ExpirationTime)
	}

	// 初始化错误urlCache
	BedUrlCache, err = mcache.NewLRU(c.BedUrlLen)
	if err != nil {
		defer fmt.Println(comm.MsgInitializeCacheError)
		panic(comm.MsgInitializeCacheError + ", err: " + err.Error())
	}

	// 开启定时任务获取需要处理的数据
	go taskDisposalQueue(queueId)
}

// taskDisposalQueue 获取需要处理的数据
func taskDisposalQueue(queueId interface{}) {
	for {
		list := db.GetQueueListById(queueId)
		count := len(list)
		if count > 0 {
			queueId = list[count-1].Id
			for _, val := range list {
				shortNum := val.ShortNum
				GoodUrlCache.Remove(shortNum)
			}
		}
		time.Sleep(30 * time.Second)
	}
}
