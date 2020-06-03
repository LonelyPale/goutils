package goutils

import "github.com/LonelyPale/goutils/errors/ecode"

type Msg = Message

type Message struct {
	Code int         `json:"code"`
	Text string      `json:"text,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func NewMessage(code int, text string, data interface{}) *Message {
	return &Message{code, text, data}
}

func NewSuccessMessage(text string, datas ...interface{}) *Message {
	code := ecode.StatusOK
	if len(text) == 0 {
		text = ecode.StatusText(code)
	}

	var data interface{}
	if len(datas) == 1 {
		data = datas[0]
	} else if len(datas) > 1 {
		data = datas
	}

	return &Message{code, text, data}
}

func NewErrorMessage(err error) *Message {
	switch e := err.(type) {
	case ecode.ErrorCode:
		return &Message{e.Code(), e.Error(), e.Details()}
	default:
		emsg := e.Error()
		if len(emsg) > 0 {
			return &Message{Code: ecode.StatusUndefinedError, Text: emsg}
		} else {
			return &Message{Code: ecode.StatusUndefinedError, Text: ecode.StatusText(ecode.StatusUndefinedError), Data: err}
		}
	}
}
