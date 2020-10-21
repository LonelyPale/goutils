package controller

import (
	"net/http"

	"github.com/go-spring/spring-web"

	"github.com/LonelyPale/goutils"
	"github.com/LonelyPale/goutils/errors"
	"github.com/LonelyPale/goutils/errors/ecode"
)

func init() {
	SpringWeb.RpcInvoke = WebInvoke
}

// WebInvoke 替换 spring-web 默认的 rpc 执行函数
func WebInvoke(webCtx SpringWeb.WebContext, fn func(SpringWeb.WebContext) interface{}) {
	// 目前 HTTP RPC 只能返回 json 格式的数据
	webCtx.Header("Content-Type", "application/json")

	defer func() {
		if r := recover(); r != nil {
			result, ok := r.(error)
			if !ok {
				result = errors.UnknownError(r)
			}
			_ = webCtx.JSON(http.StatusOK, goutils.NewErrorMessage(result))
		}
	}()

	var result *goutils.Message
	switch v := fn(webCtx).(type) {
	case error:
		result = goutils.NewErrorMessage(v)
	case ecode.ErrorCode:
		result = goutils.NewErrorMessage(v)
	case goutils.Message:
	default:
		result = goutils.NewSuccessMessage(v)
	}
	_ = webCtx.JSON(http.StatusOK, result)
}
