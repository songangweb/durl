package comm

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	"net/http"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/json-iterator/go"
)

type BaseController struct {
	web.Controller
}

type InterfaceResp struct {
	BaseResp
	Data interface{} `json:"data"`
}

type BaseResp struct {
	Code    int    `json:"code"`    // 返回码
	Message string `json:"message"` // 错误描述
}

type BaseListResp struct {
	Len  int64       `json:"len"`
	List interface{} `json:"list"`
}

type ListData struct {
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Total    int           `json:"total"`
	Len      int           `json:"len"`
	List     []interface{} `json:"list"`
}

// 函数名称: sendResponse
// 功能: 回复应答消息
// 输入参数:
//     httpCode: http状态码
//     code: code返回值
//     message: msg返回值
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021-11-17 15:15:42 #
func (b *BaseController) sendResponse(httpCode int, code int, message string) {
	m := make(map[string]interface{})

	m["code"] = code
	m["message"] = message

	str, _ := jsoniter.Marshal(m)

	/* 应答 */
	b.Ctx.ResponseWriter.WriteHeader(httpCode)
	b.Ctx.ResponseWriter.Write(str)

	b.StopRun()
}

// 函数名称: FormatResp
// 功能: 回复应答消息
// 输入参数:
//     httpCode: http状态码
//     code: code返回值
//     message: msg返回值
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021-11-17 15:15:42 #

func (b *BaseController) FormatResp(httpCode, code int, message string) {
	b.Ctx.Output.SetStatus(httpCode)

	resp := BaseResp{
		Code:    code,
		Message: message,
	}
	b.Data["json"] = resp
	b.ServeJSON()
}

// 函数名称: FormatInterfaceResp
// 功能: 回复应答消息
// 输入参数:
//     httpCode: http状态码
//     code: code返回值
//     message: msg返回值
//     i: 自定义内容data
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021-11-17 15:15:42 #

func (b *BaseController) FormatInterfaceResp(httpCode, code int, message string, i interface{}) {
	curNow := time.Now()
	defer func() {
		logs.Info("Time used FormatInterfaceResp:", time.Now().Sub(curNow))
	}()

	b.Ctx.Output.SetStatus(httpCode)
	resp := InterfaceResp{
		BaseResp: BaseResp{
			Code:    code,
			Message: message,
		},
		Data: i,
	}

	b.Data["json"] = resp
	b.ServeJSON()
}

// 函数名称: FormatInterfaceListResp
// 功能: 回复列表消息
// 输入参数:
//     httpCode: http状态码
//     code: code返回值
//     message: msg返回值
//     i: 自定义内容data
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021-11-17 15:15:42 #

func (b *BaseController) FormatInterfaceListResp(httpCode, code int, len int64, message string, i interface{}) {
	b.Ctx.Output.SetStatus(httpCode)

	// 为空时转换为 []
	if i == nil || fmt.Sprintf("%v", i) == "[]" {
		i = make([]string, 0)
	}

	resp := InterfaceResp{
		BaseResp: BaseResp{
			Code:    code,
			Message: message,
		},
		Data: BaseListResp{
			Len:  len,
			List: i,
		},
	}

	b.Data["json"] = resp
	b.ServeJSON()
}

// 函数名称: ErrorMessage
// 功能: 回复错误消息
// 输入参数:
//     code: code返回值
//     message: msg返回值
// 输出参数:
// 返回: 回复错误消息
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021-11-17 15:15:42 #

func (b *BaseController) ErrorMessage(code int, message string) {

	if code >= ErrBadReq && code < ErrAuth {
		b.BadRequest(code, message)
		return
	} else if code >= ErrAuth && code < ErrForbidden {
		b.Unauthorized(code, message)
		return
	} else if code >= ErrForbidden && code < ErrNotFound {
		b.Forbidden(code, message)
		return
	} else if code >= ErrNotFound && code < ErrMethodNotAllowed {
		b.NotFound(code, message)
		return
	} else if code >= ErrMethodNotAllowed && code < ErrGone {
		b.MethodNotAllowed(code, message)
		return
	} else if code >= ErrGone && code < ErrUnsupportedMediaType {
		b.Gone(code, message)
		return
	} else if code >= ErrUnsupportedMediaType && code < ErrUnprocessableEntity {
		b.UnsupportedMediaType(code, message)
		return
	} else if code >= ErrUnprocessableEntity && code < ErrTooManyReq {
		b.UnprocessableEntity(code, message)
		return
	} else if code >= ErrTooManyReq && code < ErrInternalServerError {
		b.TooManyRequests(code, message)
		return
	} else if code >= ErrInternalServerError && code < ErrSvcUnavailable {
		b.InternalServerError(code, message)
		return
	}
	b.sendResponse(http.StatusOK, code, message)
	return
}

// 1.服务端错误
// 1.1 Internal server error
//  @状态码: 500
//  @状态含义: Internal server error
//  @状态原因: 客户端请求有效, 服务器处理时发生了意外!
//  @错误码

func (b *BaseController) InternalServerError(code int, message string) {
	b.sendResponse(http.StatusInternalServerError, code, message)
}

// 2.客户端错误
// 2.1 Bad request
//  @状态码: 400
//  @状态含义: Bad request
//  @状态原因: 服务器不理解客户端的请求, 未做任何处理!
//  @错误码

func (b *BaseController) BadRequest(code int, message string) {
	b.sendResponse(http.StatusBadRequest, code, message)
}

// 2.2 Unauthorized
//  @状态码: 401
//  @状态含义: Unauthorized
//  @状态原因: 用户未提供身份验证凭据, 或者没有通过身份验证!
//  @错误码

func (b *BaseController) Unauthorized(code int, message string) {
	b.sendResponse(http.StatusUnauthorized, code, message)
}

// 2.3 Forbidden
//  @状态码: 403
//  @状态含义: Forbidden
//  @状态原因: 用户通过了身份验证, 但是不具有访问资源所需的权限!
//  @错误码

func (b *BaseController) Forbidden(code int, message string) {
	b.sendResponse(http.StatusForbidden, code, message)
}

// 2.4 Not found
//  @状态码: 404
//  @状态含义: Not found
//  @状态原因: 所请求的资源不存在, 或不可用!
//  @错误码

func (b *BaseController) NotFound(code int, message string) {
	b.sendResponse(http.StatusNotFound, code, message)
}

// 2.5 Method not allowed
//  @状态码: 405
//  @状态含义: Method not allowed
//  @状态原因: 用户已经通过身份验证, 但是所用的HTTP方法不在他的权限之内!
//  @错误码

func (b *BaseController) MethodNotAllowed(code int, message string) {
	b.sendResponse(http.StatusMethodNotAllowed, code, message)
}

// 2.6 Gone
//  @状态码: 410
//  @状态含义: Gone
//  @状态原因: 所请求的资源已从这个地址转移, 不再可用!
//  @错误码

func (b *BaseController) Gone(code int, message string) {
	b.sendResponse(http.StatusGone, code, message)
}

// 2.7 Unsupported media type
//  @状态码: 415
//  @状态含义: Unsupported media type
//  @状态原因: 客户端要求的返回格式不支持. 比如: API只能返回JSON格式, 但是客户端要求返回XML格式!
//  @错误码

func (b *BaseController) UnsupportedMediaType(code int, message string) {
	b.sendResponse(http.StatusUnsupportedMediaType, code, message)
}

// 2.8 Unprocessable entity
//  @状态码: 422
//  @状态含义: Unprocessable entity
//  @状态原因: 客户端上传的附件无法处理, 导致请求失败!
//  @错误码

func (b *BaseController) UnprocessableEntity(code int, message string) {
	b.sendResponse(http.StatusUnprocessableEntity, code, message)
}

// 2.9 Too many requests
//  @状态码: 429
//  @状态含义: Too many requests
//  @状态原因: 客户端的请求次数超过限额!
//  @错误码

func (b *BaseController) TooManyRequests(code int, message string) {
	b.sendResponse(http.StatusTooManyRequests, code, message)
}

// 函数名称: BaseCheckParams
// 功能: 接口参数基础校验
// 输入参数:
// 输出参数:
// 返回: 返回校验结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/18 5:31 下午 #

func (b *BaseController) BaseCheckParams(req interface{}) {
	controllerName, actionName := b.GetControllerAndAction()
	method := b.Ctx.Request.Method
	var err error
	if method == "GET" {
		err = b.ParseForm(req)
	} else {
		err = json.Unmarshal(b.Ctx.Input.RequestBody, req)
	}

	if err != nil {
		logs.Info("Method "+method+" Controller "+controllerName+" Action "+actionName+", err: ", err.Error())
		b.ErrorMessage(ErrParamMiss, MsgParseFormErr)
	}

	valid := validation.Validation{}
	c, err := valid.Valid(req)
	if err != nil {
		logs.Info("Method "+method+" Controller "+controllerName+" Action "+actionName+", err: ", err.Error())
		b.ErrorMessage(ErrParamMiss, MsgParseFormErr)
	}
	if !c {
		for _, err := range valid.Errors {
			logs.Info("Method "+method+" Controller "+controllerName+" Action "+actionName+", err: ", err.Key, err.Message)
		}
		b.ErrorMessage(ErrParamInvalid, MsgParseFormErr)
	}
}
