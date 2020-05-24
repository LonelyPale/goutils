package goutils

type Msg = Message

type Message struct {
	Code int         `json:"code"`
	Text string      `json:"text,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func NewMessage(code int, text string, data interface{}) *Message {
	return &Message{code, text, data}
}
