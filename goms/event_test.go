package goms

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var eb = NewEventBus()

func printEvent(ch string, event Event) {
	fmt.Printf("Channel: %s; Type: %s; Data: %v\n", ch, event.Type, event.Data)
}

func publisTo(typ string, data string) {
	for {
		eb.Publish(typ, data)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func TestEventBus(t *testing.T) {
	ch1 := make(chan Event)
	ch2 := make(chan Event)
	ch3 := make(chan Event)

	eb.Subscribe("type1", ch1)
	eb.Subscribe("type2", ch2)
	eb.Subscribe("type3", ch3)

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
	err := eb.SubscribeFunc("type1", func(event Event) {
		printEvent("ch1", event)
	})
	if err != nil {
		t.Fatal(err)
	}

	err = eb.SubscribeFunc("type2", func(event Event) {
		printEvent("ch2", event)
	})
	if err != nil {
		t.Fatal(err)
	}

	err = eb.SubscribeFunc("type3", func(event Event) {
		printEvent("ch3", event)
	})
	if err != nil {
		t.Fatal(err)
	}

	go publisTo("type1", "Hi topic 1")
	go publisTo("type2", "Welcome to topic 2")
	go publisTo("type3", "3333")

	time.Sleep(10 * time.Second)
}
