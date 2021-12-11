package controllers

import (
	"durl/app/exec/jump/mcache"
	"durl/app/share/dao/db"
	"durl/app/share/tool"
	"fmt"
	"github.com/beego/beego/v2/core/config"
)

func (c *Controller) Jump() {

	shortKey := c.Ctx.Input.Param(":jump")
	shortNum := tool.Base62Decode(shortKey)

	// 判断缓存是否存在数据
	if fullUrl, ok := mcache.NewMcache.Gget(shortNum); ok {
		reStatusFound(c, fmt.Sprint(fullUrl))
		return
	}
	// 判断错误url缓存是否存在, 如果存在返回404
	if _, ok := mcache.NewMcache.Gget(shortKey); ok {
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

// 返回404页面,并加入缓存(60秒)
func reStatusNotFoundAndCache(c *Controller, shortKey string) {
	mcache.NewMcache.Badd(shortKey, "", (tool.TimeNowUnix()+600)*1000)
	c.Abort("404")
}

// 返回跳转页面
func reStatusFound(c *Controller, fullUrl string) {
	c.Data["shortUrl"] = fullUrl
	// 百度统计key
	sConf, _ := config.String("Statistical_Baidu")
	if sConf != "" {
		c.Data["Statistical_Baidu_Key"] = sConf
	}

	c.TplName = "jump.html"
	_ = c.Render()
}

// 返回跳转页面,并加入缓存
func reStatusFoundAndCache(c *Controller, shortNum int, fullUrl string, expirationTime int) {
	mcache.NewMcache.Gadd(shortNum, fullUrl, int64(expirationTime*1000))
	mcache.NewMcache.Gget(shortNum)

	// 跳转页面
	reStatusFound(c, fullUrl)
}
