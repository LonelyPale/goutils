package main

import (
	"fmt"
	"github.com/LonelyPale/goutils"
	"golang.org/x/net/websocket"
	"sync"
	"time"
)

const logModule = "websocket"

type Msg = goutils.Message

type Client struct {
	Group   string
	Type    string
	Conn    *websocket.Conn
	Message chan *Msg
	exitSig chan int
	mutex   sync.RWMutex
}

func (c *Client) isClose() bool {
	select {
	case <-c.exitSig:
		return true
	default:
		return false
	}
}

func (c *Client) close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isClose() {
		return
	}

	close(c.exitSig)
	//if err := c.Conn.Close(); err != nil {
	//	log.WithFields(log.Fields{
	//		"module": logModule,
	//		"action": "close",
	//		"group":  c.Group,
	//		"type":   c.Type,
	//		"error":  err,
	//	})
	//}
}

func main() {
	client := Client{
		exitSig: make(chan int),
	}

	go func() {
		a, ok := <-client.exitSig
		fmt.Println("print:", a, ok)
	}()

	fmt.Println(client.isClose())

	//client.exitSig <- 123
	fmt.Println(client.isClose())

	client.close()

	fmt.Println(client.isClose())

	client.close()

	fmt.Println(client.isClose())
	fmt.Println(client.isClose())

	time.Sleep(time.Second)

}
