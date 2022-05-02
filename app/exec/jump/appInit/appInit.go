package appInit

import (
	"durl/app/exec/jump/controllers"
	"durl/app/exec/jump/routers"
	"durl/app/share/dao/cache"
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
	MsgType     string
	Db          db.DBConf
	Log         log.Conf
	MemoryCache cache.Conf
}

// initConf
// 函数名称: initConf
// 功能: 初始化配置
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 5:44 下午 #
func initConf() (AppConf *Conf) {

	AppConf = new(Conf)

	// 获取环境
	runmode, _ := config.String("runmode")

	// 获取消息通信方式
	AppConf.MsgType, _ = config.String("MsgType")

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

	// 初始化缓存
	controllers.InitCache(c.MemoryCache)

	// 初始化消息队列
	go controllers.InitMsg(c.MsgType)

}
