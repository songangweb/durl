package controllers

import (
	"durl/app/share/comm"
	"durl/app/share/dao/cache"
	"durl/app/share/dao/db"
	"durl/app/share/tool"
	"fmt"
	"github.com/beego/beego/v2/core/config"
	"net/http"
)

// Jump
// 函数名称: Jump
// 功能: 跳转
// 输入参数:
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 5:44 下午 #
func (c *Controller) Jump() {
	shortKey := c.Ctx.Input.Param(":jump")
	shortNum := tool.Base62Decode(shortKey)

	// 判断缓存是否存在数据
	if fullUrl, ok := cache.NewUrlListCache().Gget(shortNum); ok {
		reStatusFound(c, fmt.Sprint(fullUrl))
		return
	}

	// 判断错误url缓存是否存在, 如果存在返回404
	if _, ok := cache.NewUrlListCache().Bget(shortKey); ok {
		cache.NewUrlListCache().Badd(shortKey, "", (tool.TimeNowUnix()+600)*1000)
		reStatusNotFound(c)
		return
	}

	// 查询数据库
	urlDetail := db.NewDbService().GetFullUrlByShortNum(shortNum)
	// 跳转到 404 页面
	if urlDetail == nil {
		cache.NewUrlListCache().Badd(shortKey, "", (tool.TimeNowUnix()+600)*1000)
		reStatusNotFound(c)
		return
	}

	cache.NewUrlListCache().Gadd(shortNum, urlDetail.FullUrl, int64(urlDetail.ExpirationTime*1000))
	cache.NewUrlListCache().Gget(shortNum)

	reStatusFound(c, urlDetail.FullUrl)
	return
}

// reStatusFound
// 函数名称: reStatusFound
// 功能: 返回跳转页面
// 输入参数:
// 输出参数:
// 返回: 返回跳转页面
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 5:44 下午 #
func reStatusFound(c *Controller, fullUrl string) {
	c.Data["fullUrl"] = fullUrl

	// 百度统计key
	// 判断是否开启百度统计
	runmode, _ := config.String("runmode")
	statisticalBool, _ := config.Bool("Statistical")
	if statisticalBool {
		sConf, _ := config.String(runmode + "::Baidu")
		c.Data["shortUrl"] = fullUrl
		if sConf != "" {
			c.Data["Statistical_Baidu_Key"] = sConf
		}
		c.TplName = "jump.html"
		_ = c.Render()
	}

	// 直接跳转 临时重定向
	c.Redirect(fullUrl, http.StatusFound)
}

type cacheDetailResp struct {
	BedUrlCacheKeys  []interface{} `json:"bedUrlCacheKeys"`
	GoodUrlCacheKeys []interface{} `json:"goodUrlCacheKeys"`
}

// CacheDetail
// 函数名称: CacheDetail
// 功能: 返回当前缓存使用情况
// 输入参数:
// 输出参数:
// 返回: 返回当前缓存使用情况
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 5:44 下午 #
func (c *Controller) CacheDetail() {

	BedUrlCacheKeys := cache.BedUrlCache.Keys()
	GoodUrlCacheKeys := cache.GoodUrlCache.Keys()

	data := &cacheDetailResp{
		BedUrlCacheKeys:  BedUrlCacheKeys,
		GoodUrlCacheKeys: GoodUrlCacheKeys,
	}
	c.FormatInterfaceResp(comm.OK, comm.OK, comm.MsgOk, data)
	return
}
