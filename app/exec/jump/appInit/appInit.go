package appInit

import (
	"durl/app/exec/jump/cache"
	"durl/app/exec/jump/controllers"
	"durl/app/exec/jump/routers"
	"durl/app/share/dao/db"
	"durl/app/share/log"
	"github.com/beego/beego/v2/core/config"
)

func Init() {
	// 获取配置
	AppConf := initConf()

	// 初始化app相关功能
	initApp(AppConf)
}

type Conf struct {
	Db          db.Conf
	Log         log.Conf
	MemoryCache cache.Conf
}

func initConf() (AppConf *Conf) {

	AppConf = new(Conf)

	// 获取环境
	runmode, _ := config.String("runmode")

	// db
	AppConf.Db.Type, _ = config.String(runmode + "::Db_Type")
	// mysql
	AppConf.Db.Xorm.Mysql.Master, _ = config.String(runmode + "::Db_Mysql_Master")
	AppConf.Db.Xorm.Mysql.Slave1, _ = config.String(runmode + "::Db_Mysql_Slave1")
	AppConf.Db.Xorm.Mysql.SetMaxOpen, _ = config.Int(runmode + "::Db_Mysql_SetMaxOpen")
	AppConf.Db.Xorm.Mysql.SetMaxIdle, _ = config.Int(runmode + "::Db_Mysql_SetMaxIdle")

	// log
	AppConf.Log.Conf, _ = config.String(runmode + "::Log_Conf")

	// memory cache
	AppConf.MemoryCache.GoodUrlLen, _ = config.Int(runmode + "::MemoryCache_GoodUrlLen")
	AppConf.MemoryCache.BedUrlLen, _ = config.Int(runmode + "::MemoryCache_BedUrlLen")

	return AppConf
}

func initApp(c *Conf) {

	// 初始化日志服务
	c.Log.InitLog()

	// 数据库初始化
	c.Db.InitDb()

	// 初始化路由组
	routers.RouterHandler()

	// 初始化url缓存
	controllers.InitUrlCache(c.MemoryCache)

	// 初始化黑名单
	controllers.InitBlacklist()

}
