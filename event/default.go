package event

var defaultBus = NewBus()

func Publish(topic string, data interface{}) {
	defaultBus.Publish(topic, data)
}

func Subscribe(topic string, callback Handler) {
	defaultBus.Subscribe(topic, callback)
}

func Unsubscribe(topics ...string) {
	defaultBus.Unsubscribe(topics...)
}

func Close() {
	defaultBus.Close()
}
