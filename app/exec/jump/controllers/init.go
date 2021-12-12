package controllers

import (
	"durl/app/exec/jump/cache"
	"durl/app/share/dao/db"
	"github.com/beego/beego/v2/server/web"
	"time"
)

type Controller struct {
	web.Controller
}

func (c *Controller) Prepare() {
	// 过滤黑名单
	ip := c.Ctx.Input.IP()
	if cache.Blacklist.Search(ip) {
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
	queueId := db.QueueLastId()

	// 获取数据库中需要放到缓存的url
	UrlList := db.GetCacheUrlAllByLimit(c.GoodUrlLen)
	// 添加数据到缓存中
	for i := 0; i < len(UrlList); i++ {
		cache.UrlListCache.Gadd(UrlList[i].ShortNum, UrlList[i].FullUrl, int64(UrlList[i].ExpirationTime))
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
				cache.UrlListCache.Gremove(shortNum)
			}
		}
		time.Sleep(30 * time.Second)
	}
}

// InitBlacklist 初始化黑名单
func InitBlacklist() {
	// 开启定时任务获取黑名单列表
	go taskBlacklist()
}

// taskBlacklist 开启定时任务获取黑名单列表
func taskBlacklist() {
	for {
		// 初始化缓存
		cache.InitBlacklist()
		// 获取所有黑名单列表
		list := db.GetBlacklistAll()
		for _, val := range list {
			cache.Blacklist.Add(val.Ip)
		}
		time.Sleep(60 * time.Second)
	}
}