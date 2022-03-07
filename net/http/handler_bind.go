package http

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"

	"github.com/lonelypale/goutils"
	"github.com/lonelypale/goutils/encoding/json"
	"github.com/lonelypale/goutils/errors"
	"github.com/lonelypale/goutils/validator"
)

func validBindFn(fnType reflect.Type) bool {
	// fn 必须是函数
	// 入参(0-n): *gin.Context、*struct{`json`}、*struct{`form`}、*struct{`uri`}、*struct{`query`}、*struct{`header`}
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
func BIND(fn interface{}, requestHandlers []Handler, responseHandlers []Handler) gin.HandlerFunc {
	if fnType := reflect.TypeOf(fn); validBindFn(fnType) {
		bindParam := make([]Param, fnType.NumIn())
		for n := 0; n < fnType.NumIn(); n++ {
			bindParam[n] = NewParam(fnType.In(n))
		}

		binder := &bindHandler{
			fn:               fn,
			fnType:           fnType,
			fnValue:          reflect.ValueOf(fn),
			bindParam:        bindParam,
			requestHandlers:  requestHandlers,
			responseHandlers: responseHandlers,
		}
		return binder.Invoke
	}

	panic(errors.New("fn should be func(*gin.Context、*struct{`json`}、*struct{`form`}、*struct{`uri`}、*struct{`query`}、*struct{`header`}、*struct})anything"))
}

// bindHandler BIND 形式的 Web 处理接口
type bindHandler struct {
	fn               interface{}
	fnType           reflect.Type
	fnValue          reflect.Value
	bindParam        []Param
	requestHandlers  []Handler
	responseHandlers []Handler
}

func (b *bindHandler) Invoke(ctx *gin.Context) {
	WebInvoke(ctx, b, b.call)
}

func (b *bindHandler) call(ctx *gin.Context) []interface{} {
	in := make([]reflect.Value, len(b.bindParam))

	// 反射创建需要绑定请求参数
	for i := 0; i < len(b.bindParam); i++ {
		param := b.bindParam[i]

		if param.ParamType == ParamGinContext {
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

		switch param.ParamType {
		case ParamJsonStruct:
			err = ctx.ShouldBindJSON(bindVal.Interface())
		case ParamFormStruct:
			err = ctx.ShouldBindWith(bindVal.Interface(), binding.Form)
		case ParamUriStruct:
			err = ctx.ShouldBindUri(bindVal.Interface())
		case ParamQueryStruct:
			err = ctx.ShouldBindQuery(bindVal.Interface())
		case ParamHeaderStruct:
			err = ctx.ShouldBindHeader(bindVal.Interface())
		case ParamStruct, ParamOther:
			err = json.NewDecoder(ctx.Request.Body).Decode(bindVal.Interface())
		default:
			err = errors.New("invalid param")
		}
		if err != nil {
			panic(err)
		}

		//验证绑定参数：gin的bind方法会自己验证参数，所以只用验证ParamStruct
		if param.ParamType == ParamStruct {
			if err = validator.Validate(bindVal.Interface()); err != nil {
				panic(err)
			}
		}

		if b.bindParam[i].Type.Kind() == reflect.Ptr {
			in[i] = bindVal
		} else {
			in[i] = bindVal.Elem()
		}
	}

	//处理 RequestHandler 回调
	params := make([]interface{}, len(in))
	if len(b.requestHandlers) > 0 {
		for i, v := range in {
			params[i] = v.Interface()
		}
	}

	for _, reqHander := range b.requestHandlers {
		if reqHander == nil {
			continue
		}
		var err error
		params, err = reqHander.Invoke(ctx, params)
		if err != nil {
			log.Error(err)
		}
	}

	if len(b.requestHandlers) > 0 {
		for i, p := range params {
			in[i] = reflect.ValueOf(p)
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

// WebInvoke 可自定义的 web 执行函数
var WebInvoke = defaultWebInvoke

// todo: 目前 HTTP Web 只能返回 json 格式的数据
// defaultWebInvoke 默认的 web 执行函数
func defaultWebInvoke(ctx *gin.Context, bind *bindHandler, fn func(*gin.Context) []interface{}) {
	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(error)
			if !ok {
				e = errors.UnknownError(r)
			}
			log.Error(e)
			ctx.JSON(http.StatusOK, goutils.NewErrorMessage(e))
		}
	}()

	var result *goutils.Message
	out := fn(ctx)

	//处理 ResponseHandler 回调
	for _, respHandler := range bind.responseHandlers {
		if respHandler == nil {
			continue
		}
		var err error
		out, err = respHandler.Invoke(ctx, out)
		if err != nil {
			log.Error(err)
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
				log.Error(v)
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
			log.Error(l)
		} else {
			out = out[:lastIndex]
			if len(out) == 1 {
				result = goutils.NewSuccessMessage(out[0])
			} else {
				result = goutils.NewSuccessMessage(out)
			}
		}
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	ctx.Data(http.StatusOK, gin.MIMEJSON, b)
}
