package main

import "fmt"

func main() {
	c := make(chan struct{}, 1)

	close(c)

	<-c
	<-c
	fmt.Println(<-c)
}
