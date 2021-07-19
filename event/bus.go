package event

import (
	log "github.com/sirupsen/logrus"
)

var (
	DefaultLogger Logger = log.New()
)

type Logger interface {
	Error(args ...interface{})
}

// Bus 存储有关订阅者感兴趣的特定主题的信息
type Bus struct {
	eventRouter *router
}

func NewBus(filters ...Filter) *Bus {
	return &Bus{
		eventRouter: newRouter(filters...),
	}
}

// 发布
func (b *Bus) Publish(topic string, data interface{}) {
	b.eventRouter.publish(&Event{
		Type: topic,
		Data: data,
	})
}

// 订阅
func (b *Bus) Subscribe(topic string, callback Handler) *Token {
	return b.eventRouter.addRoute(topic, callback)
}

// 退订
func (b *Bus) Unsubscribe(topics ...string) {
	b.eventRouter.unsubscribe(topics...)
}

// 安全释放
func (b *Bus) Close() {
	b.eventRouter.done()
}
