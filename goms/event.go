package goms

type Event struct {
	Type string
	Data interface{}
}

// EventChannel 是一个能接收 Event 的 channel
type EventChan chan Event

// EventChans 是一个包含 EventChan 数据的切片
type EventChans []EventChan

type EventFunc func(event Event)
