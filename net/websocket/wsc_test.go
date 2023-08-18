package websocket

import (
	"fmt"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func TestWscClient(t *testing.T) {
	exampleClient()
}

func exampleClient() {
	done := make(chan bool)
	ws := New("ws://127.0.0.1:8080/ws1")
	ws.WebSocket.RequestHeader.Set("Sec-Websocket-Protocol", "test.push")

	// 可自定义配置，不使用默认配置
	ws.SetConfig(&WscConfig{
		WriteWait:         10 * time.Second, // 写超时
		MaxMessageSize:    2048,             // 支持接受的消息最大长度，默认512字节
		MinRecTime:        2 * time.Second,  // 最小重连时间间隔
		MaxRecTime:        60 * time.Second, // 最大重连时间间隔
		RecFactor:         1.5,              // 每次重连失败继续重连的时间间隔递增的乘数因子，递增到最大重连时间间隔为止
		MessageBufferSize: 1024,             // 消息发送缓冲池大小，默认256
	})

	// 设置回调处理
	ws.OnConnected(func() {
		log.Println("OnConnected: ", ws.WebSocket.Url)
		// 连接成功后，测试每5秒发送消息
		go func() {
			t := time.NewTicker(time.Second * 5)
			for {
				select {
				case <-t.C:
					err := ws.WebSocket.Conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(time.Second))
					if err != nil {
						log.Error(err)
						continue
					}

					msg := &Message{
						Type: "test.other",
						Data: fmt.Sprintf("wsc-client %s", time.Now().Format(time.DateTime)),
					}

					if err := ws.WebSocket.Conn.WriteJSON(msg); err != nil {
						log.Error(err)
						continue
					}

					//time.Sleep(time.Second * 3)
					//ws.Close()
				}
			}
		}()
	})
	ws.OnConnectError(func(err error) {
		log.Println("OnConnectError: ", err.Error())
	})
	ws.OnDisconnected(func(err error) {
		log.Println("OnDisconnected: ", err.Error())
	})
	ws.OnClose(func(code int, text string) {
		log.Println("OnClose: ", code, text)
		done <- true
	})
	ws.OnTextMessageSent(func(message string) {
		log.Println("OnTextMessageSent: ", message)
	})
	ws.OnBinaryMessageSent(func(data []byte) {
		log.Println("OnBinaryMessageSent: ", string(data))
	})
	ws.OnSentError(func(err error) {
		log.Println("OnSentError: ", err.Error())
	})
	ws.OnPingReceived(func(appData string) {
		log.Println("OnPingReceived: ", appData)
	})
	ws.OnPongReceived(func(appData string) {
		log.Println("OnPongReceived: ", appData)
	})
	ws.OnTextMessageReceived(func(message string) {
		log.Println("OnTextMessageReceived: ", message)
	})
	ws.OnBinaryMessageReceived(func(data []byte) {
		log.Println("OnBinaryMessageReceived: ", string(data))
	})

	// 开始连接
	go ws.Connect()

	for {
		select {
		case <-done:
			return
		}
	}
}
