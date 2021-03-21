package event

//todo: 把 web 和 websocket 的 Handler 从业务中剥离，提炼出通用模块。

type Event struct {
	Type string
	Data interface{}
}

type HandlerFunc func(event Event)

// Chan 是一个能接收 Event 的 channel
type Chan chan Event

// Chans 是一个包含 Chan 数据的切片
type Chans []Chan