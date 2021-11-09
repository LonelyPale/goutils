package event

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var testBus = NewBus()

func printEvent(event Event) {
	fmt.Printf("Type: %s; Data: %v\n", event.Type, event.Data)
}

func publisTo(typ string, data string) {
	for {
		testBus.Publish(typ, data)
		time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
	}
}

func TestEventBus(t *testing.T) {
	testBus.Subscribe("type1", func(event Event) {
		printEvent(event)
	})

	testBus.Subscribe("type2", func(event Event) {
		printEvent(event)
	})

	testBus.Subscribe("type3", func(event Event) {
		printEvent(event)
	})

	testBus.Subscribe("type3", func(event Event) {
		fmt.Println("==>", event.Type)
	})

	go publisTo("type1", "Hi topic 1")
	go publisTo("type2", "Welcome to topic 2")
	go publisTo("type3", "3333")

	time.Sleep(3 * time.Second)

	fmt.Println("event close: 1")
	testBus.Close()
	fmt.Println("event close: 2")
}
