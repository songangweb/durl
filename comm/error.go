package comm

//	正常
//  @状态码: 200
//  @状态含义: Normal
//  @状态原因: 无异常
//  @错误码
const (
	OK = 200 // 正常
)

//	Bad request
//  @状态码: 400
//  @状态含义: Bad request
//  @状态原因: 服务器不理解客户端的请求, 未做任何处理!
//  @错误码
const (
	ErrBadReq              = 400001 // 非法请求
	ErrParamMiss           = 400002 // 参数缺失
	ErrParamInvalid        = 400003 // 参数非法
	ERR_HEAD_INVALID       = 400004 // 报头非法
	ERR_BODY_INVALID       = 400005 // 报体非法
	ERR_SIGN_DATA_INVALID  = 400006 // 签名数据非法
	ERR_REPEAT_COMMIT      = 400007 // 禁止重复提交
	ERR_DUPLICATE_ENTRY    = 400008 // 参数冲突
	ERR_PARAM_EXCEED_LIMIT = 400009 // 参数长度超出范围
)

//	Unauthorized
//  @状态码: 401
//  @状态含义: Unauthorized
//  @状态原因: 用户未提供身份验证凭据, 或者没有通过身份验证!
//  @错误码
const (
	ErrAuth = 401001 // 鉴权失败
)

//	Forbidden
//  @状态码: 403
//  @状态含义: Forbidden
//  @状态原因: 用户通过了身份验证, 但是不具有访问资源所需的权限!
//  @错误码
const (
	ErrForbidden         = 403001 // 访问受限
	ERR_USR_FROZEN       = 403002 // 用户被冻结
	ERR_BUSINESS_FROZEN  = 403003 // 业务线被冻结
	ERR_ACTION_FORBIDDEN = 403004 // 操作行为受限
)

//	Not found
//  @状态码: 404
//  @状态含义: Not found
//  @状态原因: 所请求的资源不存在, 或不可用!
//  @错误码
const (
	ErrNotFound = 404001 // 资源不存在
)

//	Method not allowed
//  @状态码: 405
//  @状态含义: Method not allowed
//  @状态原因: 用户已经通过身份验证, 但是所用的HTTP方法不在他的权限之内!
//  @错误码
const (
	ErrMethodNotAllowed = 405001 // 无本HTTP访问权限
)

//	Gone
//  @状态码: 410
//  @状态含义: Gone
//  @状态原因: 所请求的资源已从这个地址转移, 不再可用!
//  @错误码
const (
	ErrGone = 410001 // 资源不再可用
)

//	Unsupported media type
//  @状态码: 415
//  @状态含义: Unsupported media type
//  @状态原因: 客户端要求的返回格式不支持. 比如: API只能返回JSON格式, 但是客户端要求返回XML格式!
//  @错误码
const (
	ErrUnsupportedMediaType = 415001 // 不支持的返回格式
)

//  Unprocessable entity
//  @状态码: 422
//  @状态含义: Unprocessable entity
//  @状态原因: 客户端上传的附件无法处理, 导致请求失败!
//  @错误码
const (
	ErrUnprocessableEntity = 422001 // 不支持的返回格式
)

//  Too many requests
//  @状态码: 429
//  @状态含义: Too many requests
//  @状态原因: 客户端的请求次数超过限额!
//  @错误码
const (
	ErrTooManyReq        = 429001 // 请求次数超过限制
	ERR_BUS_REQ_TOO_MANY = 429002 // 超过业务频控限制
	ERR_APP_REQ_TOO_MANY = 429003 // 请求次数超过限制
)

//	Internal server error
//	@状态码: 500
//  @状态含义: Internal server error
//  @状态原因: 客户端请求有效, 服务器处理时发生了意外!
//  @错误码
const (
	ErrInternalServerError = 500001 // 服务器内部错误
	ErrTaskError           = 500002 // 任务池错误
	ERR_SYS_TOO_BUSY       = 500002 // 系统繁忙
	ERR_SYS_RPC            = 500003 // RPC异常
	ErrSysDb               = 500004 // 数据库异常
	ERR_SYS_MYSQL          = 500005 // MYSQL异常
	ERR_SYS_REDIS          = 500006 // REDIS异常
	ERR_SYS_MONGO          = 500007 // MONGO异常
	ERR_SYS_SEQ_EXHAUSTION = 500008 // 序列号耗尽
	ERR_SYS_HTTP           = 500009 // HTTP请求失败
	ERR_MAKE_DIR           = 500010 // 创建目录失败
	ERR_BACKEND_CONFIG     = 500011 // 后台数据配置错误
)

//  Service unavailable
//  @状态码: 503
//  @状态含义: Service unavailable
//  @状态原因: 服务器无法处理请求, 一般用于网站维护状态!
//  @错误码
const (
	ErrSvcUnavailable = 503001 // 服务不可达
)
