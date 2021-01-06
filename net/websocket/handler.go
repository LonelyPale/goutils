package websocket

import (
	"reflect"

	"github.com/gorilla/websocket"

	"github.com/LonelyPale/goutils/encoding/json"
	"github.com/LonelyPale/goutils/errors"
	"github.com/LonelyPale/goutils/validator"
)

func validBindFn(fnType reflect.Type) bool {
	// fn 必须是函数
	// 入参(0-n): *Conn、*Message、*struct{}、*map[]、*slice[]、*string、*number、*bool
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

func BIND(fn interface{}) Handler {
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
	panic(errors.New("fn should be func(*Conn、*Message、*struct{}、*map[]、*slice[]、*string、*number、*bool) anything"))
}

// bindHandler BIND 形式的 websocket 处理接口
type bindHandler struct {
	fn        interface{}
	fnType    reflect.Type
	fnValue   reflect.Value
	bindParam []Param
}

func (b *bindHandler) Invoke(conn *Conn, msg *Message) {
	HandlerInvoke(conn, msg, b.call)
}

func (b *bindHandler) call(conn *Conn, msg *Message) []interface{} {
	in := make([]reflect.Value, len(b.bindParam))

	// 反射创建需要绑定请求参数
	for i := 0; i < len(b.bindParam); i++ {
		param := b.bindParam[i]

		if param.ParamType == ParamConn {
			in[i] = reflect.ValueOf(conn)
			continue
		} else if param.ParamType == ParamMessage {
			in[i] = reflect.ValueOf(msg)
			continue
		}

		var err error
		var bindVal reflect.Value
		if b.bindParam[i].Type.Kind() == reflect.Ptr {
			bindVal = reflect.New(b.bindParam[i].Type.Elem())
		} else {
			bindVal = reflect.New(b.bindParam[i].Type)
		}

		data := msg.Data.(json.RawMessage)
		err = json.Unmarshal(data, bindVal.Interface())
		errors.Panic(err).When(err != nil)

		//验证绑定参数
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

// HandlerInvoke 可自定义的 websocket 执行函数
var HandlerInvoke = defaultHandlerInvoke

// defaultHandlerInvoke 默认的 websocket 执行函数
func defaultHandlerInvoke(conn *Conn, msg *Message, fn func(conn *Conn, msg *Message) []interface{}) {
	// 目前 WebSocket 只能处理 json 格式的数据

	defer func() {
		if r := recover(); r != nil {
			DefaultLogger.Error(r)
		}
	}()

	var result *Message
	out := fn(conn, msg)

	switch len(out) {
	case 0:
		return
	case 1:
		switch v := out[0].(type) {
		case Message:
			result = &v
		case *Message:
			result = v
		case error:
			result = NewErrorMessage(v)
			DefaultLogger.Error(v)
		default:
			result = NewSuccessMessage(v)
		}
	default:
		last := out[len(out)-1]
		if l, ok := last.(error); ok && l != nil {
			result = NewErrorMessage(l)
			DefaultLogger.Error(l)
		} else {
			result = NewSuccessMessage(out)
		}
	}

	result.Type = msg.Type
	result.SN = msg.SN
	bs, err := json.Marshal(result)
	if err != nil {
		DefaultLogger.Error(err)
		return
	}

	if err := conn.Write(&WSMessage{Type: websocket.TextMessage, Data: bs}); err != nil {
		DefaultLogger.Error(err)
	}
}
