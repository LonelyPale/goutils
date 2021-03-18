package event

var defaultBus = NewBus()

func Publish(typ string, data interface{}) {
	defaultBus.Publish(typ, data)
}

func Subscribe(typ string, ch Chan) {
	defaultBus.Subscribe(typ, ch)
}

func SubscribeFunc(typ string, fun HandlerFunc) {
	defaultBus.SubscribeFunc(typ, fun)
}
