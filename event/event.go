package event

//todo: 把 web 和 websocket 的 Handler 从业务中剥离，提炼出通用模块。

type Event struct {
	Type string
	Data interface{}
}

//处理事件的回调方法
type Handler func(event Event)

//是否发布数据到对应订阅的过滤器
//pTopic 发布的主题
//sTopic 订阅的主题
type Filter func(pTopic, sTopic string) bool
