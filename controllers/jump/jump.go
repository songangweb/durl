package jump

import (
	"durl/dao/db"
	"durl/tool"
	"fmt"
	"github.com/beego/beego/v2/core/config"
)

func (c *Controller) Jump() {

	shortKey := c.Ctx.Input.Param(":jump")
	shortNum := tool.Base62Decode(shortKey)

	// 判断缓存是否存在数据
	if fullUrl, _, ok := GoodUrlCache.Get(shortNum); ok {
		reStatusFound(c, fmt.Sprint(fullUrl))
		return
	}

	// 判断错误url缓存是否存在, 如果存在返回404
	if _, _, ok := BedUrlCache.Get(shortKey); ok {
		reStatusNotFoundAndCache(c, shortKey)
		return
	}

	// 查询数据库
	urlDetail := db.GetFullUrlByshortNum(shortNum)

	// 跳转到 404 页面
	if urlDetail == nil {
		reStatusNotFoundAndCache(c, shortKey)
		return
	}

	reStatusFoundAndCache(c, shortNum, urlDetail.FullUrl, urlDetail.ExpirationTime)
	return
}

// 返回404页面,并加入缓存
func reStatusNotFoundAndCache(c *Controller, shortKey string) {
	BedUrlCache.Add(shortKey, "", (tool.TimeNowUnix()+600)*1000)
	c.Abort("404")
}

// 返回跳转页面
func reStatusFound(c *Controller, fullUrl string) {
	c.Data["url"] = fullUrl

	// 百度统计key
	sConf, _ := config.String("Statistical_Baidu")
	if sConf != "" {
		c.Data["Statistical_Baidu_Key"] = sConf
	}

	c.TplName = "jump.html"
}

// 返回跳转页面,并加入缓存
func reStatusFoundAndCache(c *Controller, shortNum int, fullUrl string, expirationTime int64) {
	GoodUrlCache.Add(shortNum, fullUrl, expirationTime*1000)
	GoodUrlCache.Get(shortNum)

	// 跳转页面
	reStatusFound(c, fullUrl)
}
