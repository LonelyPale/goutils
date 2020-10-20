package goutils

import "github.com/LonelyPale/goutils/errors/ecode"

type Message struct {
	Code int         `json:"code"`           //状态码
	Msg  string      `json:"msg,omitempty"`  //消息
	Data interface{} `json:"data,omitempty"` //结果数据
}

func NewMessage(code int, msg string, datas ...interface{}) *Message {
	var data interface{}
	if len(datas) == 1 {
		data = datas[0]
	} else if len(datas) > 1 {
		data = datas
	}
	return &Message{code, msg, data}
}

func NewSuccessMessage(datas ...interface{}) *Message {
	return NewMessage(ecode.StatusOK, ecode.StatusText(ecode.StatusOK), datas...)
}

func NewErrorMessage(err error) *Message {
	switch e := err.(type) {
	case ecode.ErrorCode:
		return &Message{Code: e.Code(), Msg: e.Error()}
	default:
		emsg := e.Error()
		if len(emsg) > 0 {
			return &Message{Code: ecode.StatusUndefinedError, Msg: emsg}
		} else {
			return &Message{Code: ecode.StatusUndefinedError, Msg: ecode.StatusText(ecode.StatusUndefinedError)}
		}
	}
}
