package springweb

import (
	"net/http"

	"github.com/go-spring/spring-logger"
	"github.com/go-spring/spring-web"

	"github.com/LonelyPale/goutils"
	"github.com/LonelyPale/goutils/errors"
)

func init() {
	SpringWeb.RpcInvoke = WebInvoke
	SpringWeb.Validator = SpringWeb.NewDefaultValidator()
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
			SpringLogger.Error(result)
		}
	}()

	var result *goutils.Message
	switch v := fn(webCtx).(type) {
	case goutils.Message:
	case error:
		result = goutils.NewErrorMessage(v)
		SpringLogger.Error(v)
	default:
		result = goutils.NewSuccessMessage(v)
	}
	_ = webCtx.JSON(http.StatusOK, result)
}
