package appInit

import (
	"durl/app/exec/backend/routers"
	"durl/app/share/dao/db"
	"durl/app/share/log"

	"github.com/beego/beego/v2/core/config"
)

// Init
// 函数名称: Init
// 功能: 配置初始化
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 5:44 下午 #
func Init() {
	// 获取配置
	AppConf := initConf()

	// 初始化app相关功能
	initApp(AppConf)
}

type Conf struct {
	Db  db.DBConf
	Log log.Conf
}

// initConf
// 函数名称: initConf
// 功能: 配置初始化
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

	// db
	AppConf.Db.Type, _ = config.String(runmode + "::Db_Type")
	// mysql
	AppConf.Db.Xorm.Mysql.Master, _ = config.String(runmode + "::Db_Mysql_Master")
	AppConf.Db.Xorm.Mysql.Slave1, _ = config.String(runmode + "::Db_Mysql_Slave1")
	AppConf.Db.Xorm.Mysql.SetMaxOpen, _ = config.Int(runmode + "::Db_Mysql_SetMaxOpen")
	AppConf.Db.Xorm.Mysql.SetMaxIdle, _ = config.Int(runmode + "::Db_Mysql_SetMaxIdle")

	// log
	AppConf.Log.Conf, _ = config.String(runmode + "::Log_Conf")

	return AppConf
}

// initApp
// 函数名称: initApp
// 功能: 组件初始化
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 5:44 下午 #
func initApp(c *Conf) {

	// 初始化日志服务
	c.Log.InitLog()

	// 数据库初始化
	c.Db.InitDb()

	// 初始化路由组
	routers.RouterHandler()
}
