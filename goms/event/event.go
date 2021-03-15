package event

import "github.com/LonelyPale/goutils/goms"

var eb = goms.NewEventBus()

func Publish(typ string, data interface{}) {
	eb.Publish(typ, data)
}

func Subscribe(typ string, ch goms.EventChan) {
	eb.Subscribe(typ, ch)
}

func SubscribeFunc(typ string, fun goms.EventFunc) error {
	return eb.SubscribeFunc(typ, fun)
}
