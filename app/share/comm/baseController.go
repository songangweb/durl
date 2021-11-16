package comm

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"net/http"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/json-iterator/go"
)

type BaseController struct {
	web.Controller
}

type AddResp struct {
	BaseResp
	Data Id `json:"data"`
}

type InterfaceResp struct {
	BaseResp
	Data interface{} `json:"data"`
}

type BaseResp struct {
	Code    int    `json:"code"`    // 返回码
	Message string `json:"message"` // 错误描述
}

type Id struct {
	Id int64 `json:"id"`
}

type ListData struct {
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Total    int           `json:"total"`
	Len      int           `json:"len"`
	List     []interface{} `json:"list"`
}

/* 回复应答消息 */
func (b *BaseController) sendResponse(status int, code int, message string) {
	m := make(map[string]interface{})

	m["code"] = code
	m["message"] = message

	str, _ := jsoniter.Marshal(m)

	/* 应答 */
	b.Ctx.ResponseWriter.WriteHeader(status)
	b.Ctx.ResponseWriter.Write(str)

	b.StopRun()
}

func (c *BaseController) FormatResp(httpCode, code int, message string) {
	c.Ctx.Output.SetStatus(httpCode)

	resp := BaseResp{
		Code:    code,
		Message: message,
	}
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *BaseController) FormatInterfaceResp(httpCode, code int, message string, i interface{}) {
	curNow := time.Now()
	defer func() {
		logs.Info("Time used FormatInterfaceResp:", time.Now().Sub(curNow))
	}()

	c.Ctx.Output.SetStatus(httpCode)
	resp := InterfaceResp{
		BaseResp: BaseResp{
			Code:    code,
			Message: message,
		},
		Data: i,
	}

	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *BaseController) FormatInterfaceListResp(httpCode, code, len int, message string, i interface{}) {
	c.Ctx.Output.SetStatus(httpCode)

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

	c.Data["json"] = resp
	c.ServeJSON()
}

type BaseListResp struct {
	Len  int         `json:"len"`
	List interface{} `json:"list"`
}

/* 消息回复 */
func (b *BaseController) SendJson(status int, m map[string]interface{}) {
	str, _ := jsoniter.Marshal(m)

	/* 应答 */
	b.Ctx.ResponseWriter.WriteHeader(status)
	b.Ctx.ResponseWriter.Write(str)

	b.StopRun()
}

/* 发送应答消息 */
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