package controllers

import (
	"time"

	"durl/app/share/dao/cache"
	"durl/app/share/dao/db"

	"github.com/beego/beego/v2/server/web"
)

type Controller struct {
	web.Controller
}

func (c *Controller) Prepare() {
	// 过滤黑名单
	ip := c.Ctx.Input.IP()

	cache.BlacklistConnLock.RLock()

	if cache.Blacklist.Search(ip) {
		cache.BlacklistConnLock.RUnlock()

		reStatusNotFound(c)
		return
	}
}

// 返回404页面
func reStatusNotFound(c *Controller) {
	c.Abort("404")
}

type UrlConf struct {
	GoodUrlLen int
	BedUrlLen  int
}

func InitUrlCache(c cache.Conf) {

	// 初始化缓存
	cache.InitUrlCache(c)

	// 获取任务队列表里最新的一条数据id
	engine := db.NewDbService()
	queueId := engine.QueueLastId()

	// 获取数据库中需要放到缓存的url
	UrlList := engine.GetCacheUrlAllByLimit(c.GoodUrlLen)
	// 添加数据到缓存中
	for i := 0; i < len(UrlList); i++ {
		cache.NewUrlListCache().Gadd(UrlList[i].ShortNum, UrlList[i].FullUrl, int64(UrlList[i].ExpirationTime))
	}

	// 开启定时任务获取需要处理的数据
	go taskDisposalQueue(queueId)
}

// 循环获取queue表数据时间 s
const taskQueueTime = 30

// taskDisposalQueue 获取需要处理的数据
func taskDisposalQueue(queueId int) {
	engine := db.NewDbService()
	for {
		list := engine.GetQueueListById(queueId)
		count := len(list)
		if count > 0 {
			queueId = list[count-1].Id
			for _, val := range list {
				shortNum := val.ShortNum
				cache.NewUrlListCache().Gremove(shortNum)
			}
		}
		time.Sleep(taskQueueTime * time.Second)
	}
}

// InitBlacklist 初始化黑名单
func InitBlacklist() {
	// 开启定时任务获取黑名单列表
	go taskBlacklist()
}

// 循环获取黑名单数据时间 s
const taskBlacklistTime = 3

// taskBlacklist 开启定时任务获取黑名单列表
func taskBlacklist() {
	engine := db.NewDbService()

	// 初始化缓存
	for {
		blacklist := cache.InitBlacklist()
		// 获取所有黑名单列表
		list := engine.GetBlacklistAll()
		for _, val := range list {
			blacklist.Add(val.Ip)
		}

		cache.BlacklistConnLock.Lock()
		cache.Blacklist = blacklist
		cache.BlacklistConnLock.Unlock()

		time.Sleep(taskBlacklistTime * time.Second)
	}
}
