package controllers

import (
	"durl/app/share/comm"
	"durl/app/share/dao/cache"
	"durl/app/share/dao/db"
)

type Controller struct {
	comm.BaseController
}

// Prepare
// 函数名称: cacheDetail
// 功能: base操作
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #
func (c *Controller) Prepare() {
	// 过滤黑名单
	ip := c.Ctx.Input.IP()
	cache.BlacklistConnLock.RLock()
	Bbool := cache.Blacklist.Search(ip)
	cache.BlacklistConnLock.RUnlock()
	if Bbool {
		reStatusNotFound(c)
		return
	}
}

// reStatusNotFound
// 函数名称: reStatusNotFound
// 功能: 返回404页面
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #
func reStatusNotFound(c *Controller) {
	c.Abort("404")
}

type UrlConf struct {
	GoodUrlLen int
	BedUrlLen  int
}

// InitCache
// 函数名称: InitCache
// 功能: 初始化缓存
// 输入参数: 缓存配置
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #
func InitCache(c cache.Conf) {
	// 初始化url缓存
	initUrlCache(c)
	//初始化黑名单缓存
	initBlacklist()
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
func InitMsg(msgType string) {
	c := MSGConf{
		Type: msgType,
	}
	c.InitMsg()
}

// initUrlCache
// 函数名称: initUrlCache
// 功能: 初始化url缓存
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #
func initUrlCache(c cache.Conf) {
	// 初始化缓存
	cache.InitUrlCache(c)
	// 获取数据库中需要放到缓存的url
	engine := db.NewDbService()
	UrlList := engine.GetCacheUrlAllByLimit(c.GoodUrlLen)
	// 添加数据到缓存中
	for i := 0; i < len(UrlList); i++ {
		cache.NewUrlListCache().Gadd(UrlList[i].ShortNum, UrlList[i].FullUrl, int64(UrlList[i].ExpirationTime))
	}
}

// initBlacklist
// 函数名称: initBlacklist
// 功能: 初始化黑名单
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #
func initBlacklist() {
	// 初始化缓存
	blacklist := cache.InitBlacklist()
	// 获取所有黑名单列表
	engine := db.NewDbService()
	list := engine.GetBlacklistAll()
	for _, val := range list {
		err := blacklist.Add(val.Ip)
		if err != nil {
		}
	}
	cache.Blacklist = blacklist
}
