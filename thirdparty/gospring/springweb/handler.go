package springweb

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-spring/spring-logger"
	"github.com/go-spring/spring-utils"
	"github.com/go-spring/spring-web"

	"github.com/lonelypale/goutils"
	"github.com/lonelypale/goutils/encoding/json"
	"github.com/lonelypale/goutils/errors"
	"github.com/lonelypale/goutils/validator"
)

func validBindFn(fnType reflect.Type) bool {
	// fn 必须是函数
	// 入参(0-n): context.Context、SpringWeb.WebContext、*struct{`json`}、*struct{`form`}、*struct{`uri`}、*struct{`query`}、*struct{`header`}
	// 出参(0-n): any 无要求
	if fnType.Kind() != reflect.Func {
		return false
	}

	for n := 0; n < fnType.NumIn(); n++ {
		param := NewParam(fnType.In(n))
		if param.ParamType == ParamInvalid {
			return false
		}
	}

	return true
}

// BIND 转换成 BIND 形式的 Web 处理接口
func BIND(fn interface{}) SpringWeb.Handler {
	if fnType := reflect.TypeOf(fn); validBindFn(fnType) {
		bindParam := make([]Param, fnType.NumIn())
		for n := 0; n < fnType.NumIn(); n++ {
			bindParam[n] = NewParam(fnType.In(n))
		}

		return &bindHandler{
			fn:        fn,
			fnType:    fnType,
			fnValue:   reflect.ValueOf(fn),
			bindParam: bindParam,
		}
	}
	panic(errors.New("fn should be func(context.Context、SpringWeb.WebContext、*struct{`json`}、*struct{`form`}、*struct{`uri`}、*struct{`query`}、*struct{`header`}、*struct})anything"))
}

// bindHandler BIND 形式的 Web 处理接口
type bindHandler struct {
	fn        interface{}
	fnType    reflect.Type
	fnValue   reflect.Value
	bindParam []Param
}

func (b *bindHandler) Invoke(ctx SpringWeb.WebContext) {
	WebInvoke(ctx, b.call)
}

func (b *bindHandler) call(ctx SpringWeb.WebContext) []interface{} {
	in := make([]reflect.Value, len(b.bindParam))

	// 反射创建需要绑定请求参数
	for i := 0; i < len(b.bindParam); i++ {
		param := b.bindParam[i]

		if param.ParamType == ParamContext {
			in[i] = reflect.ValueOf(ctx.Context())
			continue
		} else if param.ParamType == ParamWebContext {
			in[i] = reflect.ValueOf(ctx)
			continue
		}

		var err error
		var bindVal reflect.Value
		if b.bindParam[i].Type.Kind() == reflect.Ptr {
			bindVal = reflect.New(b.bindParam[i].Type.Elem())
		} else {
			bindVal = reflect.New(b.bindParam[i].Type)
		}

		ginCtx := ctx.NativeContext().(*gin.Context)
		switch param.ParamType {
		case ParamJsonStruct:
			err = ginCtx.ShouldBindJSON(bindVal.Interface())
		case ParamFormStruct:
			err = ginCtx.ShouldBindWith(bindVal.Interface(), binding.Form)
		case ParamUriStruct:
			err = ginCtx.ShouldBindUri(bindVal.Interface())
		case ParamQueryStruct:
			err = ginCtx.ShouldBindQuery(bindVal.Interface())
		case ParamHeaderStruct:
			err = ginCtx.ShouldBindHeader(bindVal.Interface())
		case ParamStruct, ParamOther:
			err = json.NewDecoder(ctx.Request().Body).Decode(bindVal.Interface())
		}
		errors.Panic(err).When(err != nil)

		//验证绑定参数：gin的bind方法会自己验证参数，所以只用验证ParamStruct
		if param.ParamType == ParamStruct {
			err = validator.Validate(bindVal.Interface())
			errors.Panic(err).When(err != nil)
		}

		if b.bindParam[i].Type.Kind() == reflect.Ptr {
			in[i] = bindVal
		} else {
			in[i] = bindVal.Elem()
		}
	}

	// 执行处理函数，并返回结果
	out := b.fnValue.Call(in)
	result := make([]interface{}, len(out))
	for i, v := range out {
		result[i] = v.Interface()
	}

	return result
}

func (b *bindHandler) FileLine() (file string, line int, fnName string) {
	return SpringUtils.FileLine(b.fn)
}

// WebInvoke 可自定义的 web 执行函数
var WebInvoke = defaultWebInvoke

// todo: 目前 HTTP Web 只能返回 json 格式的数据
// defaultWebInvoke 默认的 web 执行函数
func defaultWebInvoke(webCtx SpringWeb.WebContext, fn func(SpringWeb.WebContext) []interface{}) {
	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(error)
			if !ok {
				e = errors.UnknownError(r)
			}
			SpringLogger.Error(e)
			_ = webCtx.JSON(http.StatusOK, goutils.NewErrorMessage(e))
		}
	}()

	var result *goutils.Message
	out := fn(webCtx)

	//处理 ResultHandler 回调
	for _, handler := range resultHandlers {
		if handler != nil {
			out = handler.Invoke(out)
		}
	}

	switch len(out) {
	case 0:
		result = goutils.NewSuccessMessage()
	case 1:
		switch v := out[0].(type) {
		case goutils.Message:
			result = &v
		case *goutils.Message:
			result = v
		case error:
			if v != nil {
				result = goutils.NewErrorMessage(v)
				SpringLogger.Error(v)
			} else {
				result = goutils.NewSuccessMessage()
			}
		default:
			result = goutils.NewSuccessMessage(v)
		}
	default:
		lastIndex := len(out) - 1
		last := out[lastIndex]
		if l, ok := last.(error); ok {
			result = goutils.NewErrorMessage(l)
			SpringLogger.Error(l)
		} else {
			out = out[:lastIndex]
			if len(out) == 1 {
				result = goutils.NewSuccessMessage(out[0])
			} else {
				result = goutils.NewSuccessMessage(out)
			}
		}
	}

	//_ = webCtx.JSON(http.StatusOK, result)
	b, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	if err := webCtx.Blob(http.StatusOK, SpringWeb.MIMEApplicationJSONCharsetUTF8, b); err != nil {
		panic(err)
	}
}
