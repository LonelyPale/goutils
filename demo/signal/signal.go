package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt)

	if f := os.NewFile(3, ""); f != nil {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}

	time.Sleep(3 * time.Second)

	<-sigs
}
