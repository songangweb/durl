package controllers

import (
	"fmt"

	"durl/app/share/dao/cache"
	"durl/app/share/dao/db"
	"durl/app/share/tool"

	"github.com/beego/beego/v2/core/config"
)

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

//返回跳转页面
func reStatusFound(c *Controller, fullUrl string) {
	c.Data["fullUrl"] = fullUrl

	// 百度统计key
	runmode, _ := config.String("runmode")
	sConf, _ := config.String(runmode + "::Baidu")

	if sConf != "" {
		c.Data["Statistical_Baidu_Key"] = sConf
	}

	c.TplName = "jump.html"
	_ = c.Render()
}
