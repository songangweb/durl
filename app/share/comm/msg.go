package comm

// 日志相关提示
const (
	MsgLogConfIsError = "日志配置错误"
	MsgCheckLogConf   = "请检查log配置是否正确!!"
)

// db相关提示
const (
	MsgDbTypeError = "数据库类型错误"
	MsgCheckDbType = "请检查数据库类型"

	MsgDbMysqlConfError = "mysql数据库配置错误"
	MsgCheckDbMysqlConf = "请检查mysql数据库配置是否正确!!"

	MsgInitDbMysqlTable = "数据表创建失败!!"
	MsgInitDbMysqlData  = "初始化数据创建失败!!"

	MsgDbMongoConfError = "mongo数据库配置错误"
	MsgCheckDbMongoConf = "请检查mongo数据库配置是否正确!!"

	MsgInitDbMongoTable = "数据表创建失败!!"
	MsgInitDbMongoData  = "初始化数据创建失败!!"
)

// MsgCacheError 内存缓存相关提示
const (
	MsgInitializeCacheError = "初始化内存缓存数据池错误"
)

// openApi相关提示
const (
	MsgOk           = "请求成功"
	MsgNotOk        = "请求失败"
	MsgParseFormErr = "请求参数错误"
)
