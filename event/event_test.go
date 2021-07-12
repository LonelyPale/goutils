package event

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var testBus = NewBus()

func printEvent(ch string, event Event) {
	fmt.Printf("Channel: %s; Type: %s; Data: %v\n", ch, event.Type, event.Data)
}

func publisTo(typ string, data string) {
	for {
		testBus.Publish(typ, data)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func TestEventBus(t *testing.T) {
	ch1 := make(chan Event)
	ch2 := make(chan Event)
	ch3 := make(chan Event)

	testBus.Subscribe("type1", ch1)
	testBus.Subscribe("type2", ch2)
	testBus.Subscribe("type3", ch3)

	go publisTo("type1", "Hi topic 1")
	go publisTo("type2", "Welcome to topic 2")
	go publisTo("type3", "3333")

	for {
		select {
		case d := <-ch1:
			go printEvent("ch1", d)
		case d := <-ch2:
			go printEvent("ch2", d)
		case d := <-ch3:
			go printEvent("ch3", d)
		}
	}
}

func TestEventBusFunc(t *testing.T) {
	testBus.SubscribeFunc("type1", func(event Event) {
		printEvent("ch1", event)
	})

	testBus.SubscribeFunc("type2", func(event Event) {
		printEvent("ch2", event)
	})

	testBus.SubscribeFunc("type3", func(event Event) {
		printEvent("ch3", event)
	})

	go publisTo("type1", "Hi topic 1")
	go publisTo("type2", "Welcome to topic 2")
	go publisTo("type3", "3333")

	time.Sleep(10 * time.Second)
}
