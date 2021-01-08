package websocket

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/LonelyPale/goutils/errors"
)

const DefaultWebSocketKey = "WebSocket"

var (
	ErrConnClosed = errors.New("websocket connection closed")
)

type Conn struct {
	config    *Config         //配置
	conn      *websocket.Conn //WS连接
	in        chan *WSMessage //待读管道
	out       chan *WSMessage //待写管道
	quit      chan struct{}   //退出信号
	quitMu    sync.RWMutex    //退出信号锁
	processor Processor       //消息处理器
}

func NewConn(config *Config) *Conn {
	return &Conn{
		config: config,
		in:     make(chan *WSMessage, config.InChanSize),
		out:    make(chan *WSMessage, config.OutChanSize),
		quit:   make(chan struct{}),
	}
}

func (c *Conn) Open(w http.ResponseWriter, r *http.Request, h http.Header, processor Processor) error {
	var upGrader = websocket.Upgrader{
		ReadBufferSize:  c.config.ReadBufferSize,                          //读缓冲区
		WriteBufferSize: c.config.WriteBufferSize,                         //写缓冲区
		Subprotocols:    []string{r.Header.Get("Sec-WebSocket-Protocol")}, // 处理 Sec-WebSocket-Protocol Header
		CheckOrigin: func(r *http.Request) bool { // cross origin domain
			return c.config.Origin
		},
	}

	//处理子协议
	//var topics []string
	//var topicsStr string
	//val, ok := ctx.Get(DefaultWebSocketKey)
	//if ok {
	//	protocolMap := val.(map[string]string)
	//	topicsStr = protocolMap["topics"]
	//	topics = strings.Split(topicsStr, "&")
	//}
	//header := http.Header{"Sec-WebSocket-Protocol": []string{"topics#" + topicsStr}}
	//wsConn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, header)

	wsConn, err := upGrader.Upgrade(w, r, h)
	if err != nil {
		return err
	}

	c.conn = wsConn
	c.conn.SetReadLimit(c.config.MaxMessageSize)
	c.processor = processor

	go c.ReadLoop()
	go c.WriteLoop()
	go c.ProcessLoop()

	return nil
}

func (c *Conn) IsClose() bool {
	c.quitMu.RLock()
	defer c.quitMu.RUnlock()

	select {
	case <-c.quit:
		return true
	default:
		return false
	}
}

func (c *Conn) Close() error {
	c.quitMu.Lock()
	defer c.quitMu.Unlock()

	select {
	case <-c.quit:
		return nil
	default:
		close(c.quit)
		close(c.in)
		close(c.out)
		return c.conn.Close()
	}
}

func (c *Conn) ReadLoop() {
	defer c.processor.OnClose()

	for {
		if msg, err := c.read(); err != nil {
			c.processor.OnError(err, msg)
			return
		}
	}
}

func (c *Conn) read() (msg *WSMessage, err error) {
	c.quitMu.RLock()
	defer c.quitMu.RUnlock()

	defer func() {
		if r := recover(); r != nil {
			err = errors.UnknownError(r)
		}
	}()

	select {
	case <-c.quit:
		return nil, ErrConnClosed
	default:
		if c.config.ReadDeadline > 0 {
			if err := c.conn.SetReadDeadline(time.Now().Add(c.config.ReadDeadline * time.Second)); err != nil {
				return nil, err
			}
		}

		messageType, receivedData, err := c.conn.ReadMessage()
		msg = &WSMessage{
			Type: messageType,
			Data: receivedData,
		}
		if err != nil {
			return msg, err
		}

		c.in <- msg
		return msg, nil
	}
}

// 1、当conn关闭后，out管道中还有没处理完的消息(需要业务端自己处理没有发送的消息)；
// 2、当msg已发送到client后，需要callback回写数据库的情况(特殊情况业务端自己处理)；
func (c *Conn) WriteLoop() {
	defer c.processor.OnClose()

	for msg := range c.out {
		if err := c.write(msg); err != nil {
			c.processor.OnError(err, msg)
			return
		}
	}
}

func (c *Conn) write(msg *WSMessage) (err error) {
	c.quitMu.RLock()
	defer c.quitMu.RUnlock()

	defer func() {
		if r := recover(); r != nil {
			err = errors.UnknownError(r)
		}
	}()

	select {
	case <-c.quit:
		return ErrConnClosed
	default:
		if c.config.WriteDeadline > 0 {
			if err := c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteDeadline * time.Second)); err != nil {
				return err
			}
		}

		if err := c.conn.WriteMessage(msg.Type, msg.Data); err != nil {
			return err
		}

		return nil
	}
}

// 1、当conn关闭后，in管道中还有没处理完的数据；
func (c *Conn) ProcessLoop() {
	defer c.processor.OnQuit()

	for msg := range c.in {
		c.process(msg)
	}
}

func (c *Conn) process(msg *WSMessage) {
	defer func() {
		if r := recover(); r != nil {
			c.processor.OnError(errors.UnknownError(r), msg)
		}
	}()
	c.processor.OnMessage(msg)
}

func (c *Conn) Read() *WSMessage {
	return <-c.in
}

func (c *Conn) Write(msg *WSMessage) error {
	c.quitMu.RLock()
	defer c.quitMu.RUnlock()

	select {
	case <-c.quit:
		return ErrConnClosed
	case c.out <- msg:
		return nil
	}
}
