package websocket

import (
	"fmt"

	"github.com/LonelyPale/goutils/encoding/json"
	"github.com/LonelyPale/goutils/errors/ecode"
)

type Message struct {
	Type string      `json:"type"`           //消息类型
	SN   int         `json:"sn"`             //流水号
	Code int         `json:"code"`           //状态码
	Msg  string      `json:"msg,omitempty"`  //消息
	Data interface{} `json:"data,omitempty"` //结果数据
}

//todo: 优化显示
func (m Message) String() string {
	return fmt.Sprintf(`[Type:%s Data:%s]`, m.Type, m.Data)
}

func NewMessage(code int, msg string, datas ...interface{}) *Message {
	var data interface{}
	if len(datas) == 1 {
		data = datas[0]
	} else if len(datas) > 1 {
		data = datas
	}
	return &Message{Code: code, Msg: msg, Data: data}
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
			return &Message{Code: ecode.StatusError, Msg: emsg}
		} else {
			return &Message{Code: ecode.StatusError, Msg: ecode.StatusText(ecode.StatusError)}
		}
	}
}

type RawMessage struct {
	*Message
	Data json.RawMessage `json:"data,omitempty"`
}

func (r *RawMessage) Msg() *Message {
	r.Message.Data = r.Data
	return r.Message
}

type WSMessage struct {
	Type int
	Data []byte
}

func (w WSMessage) String() string {
	return fmt.Sprintf("[Type:%d Data:%s]", w.Type, w.Data)
}
