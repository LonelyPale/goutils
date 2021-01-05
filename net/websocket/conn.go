package websocket

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/LonelyPale/goutils/errors"
)

const DefaultWebSocketKey = "WebSocket"

var DefaultLogger Logger = log.New()

var (
	ErrConnClosed = errors.New("websocket connection closed")
)

type Logger interface {
	Error(args ...interface{})
}

type base interface {
	SetConn(conn *Conn) error
	OnError(err error, msg *WSMessage)
	OnClose()
}

type Reader interface {
	base
	OnRead(msg *WSMessage)
}

type Writer interface {
	base
	OnWrite(msg *WSMessage)
}

type WSMessage struct {
	Type int
	Data []byte
}

type Conn struct {
	opts   *Options
	conn   *websocket.Conn
	in     chan *WSMessage //待读管道
	out    chan *WSMessage //待写管道
	outed  chan *WSMessage //已写管道
	quit   chan struct{}   //退出信号
	quitMu sync.RWMutex
	reader Reader
	writer Writer
}

func NewConn(reader Reader, writer Writer, opts *Options) *Conn {
	return &Conn{
		opts:   opts,
		reader: reader,
		writer: writer,
		in:     make(chan *WSMessage, opts.InChanSize),
		out:    make(chan *WSMessage, opts.OutChanSize),
		outed:  make(chan *WSMessage, opts.OutedChanSize),
		quit:   make(chan struct{}),
	}
}

func (c *Conn) Open(ctx *gin.Context) error {
	var upGrader = websocket.Upgrader{
		ReadBufferSize:  c.opts.ReadBufferSize,                             //读缓冲区
		WriteBufferSize: c.opts.WriteBufferSize,                            //写缓冲区
		Subprotocols:    []string{ctx.GetHeader("Sec-WebSocket-Protocol")}, // 处理 Sec-WebSocket-Protocol Header
		CheckOrigin: func(r *http.Request) bool { // cross origin domain
			return c.opts.Origin
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

	wsConn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return err
	}

	c.conn = wsConn
	c.conn.SetReadLimit(c.opts.MaxMessageSize)

	if err := c.reader.SetConn(c); err != nil {
		return err
	}

	if err := c.writer.SetConn(c); err != nil {
		return err
	}

	return nil
}

func (c *Conn) Start() {
	if c.reader != nil {
		go c.ReadLoop()
	}
	if c.writer != nil {
		go c.WriteLoop()
	}
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
		close(c.outed)
		return c.conn.Close()
	}
}

// 1、处理callback事件
func (c *Conn) ReadLoop() {
	defer c.reader.OnClose()

	go c.readProcess()

	for {
		if msg, err := c.read(); err != nil {
			c.reader.OnError(err, msg)
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
		if c.opts.ReadDeadline > 0 {
			if err := c.conn.SetReadDeadline(time.Now().Add(c.opts.ReadDeadline * time.Second)); err != nil {
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

// 1、当conn关闭后，out队列中还有没处理完的数据(需要业务端自己记录处理没有执行已完成callback的msg)；
// 2、当msg已发送到client后，需要callback回写数据库的情况；
// 3、处理callback事件
func (c *Conn) WriteLoop() {
	defer c.writer.OnClose()

	go c.writeProcess()

	for msg := range c.out {
		if err := c.write(msg); err != nil {
			c.writer.OnError(err, msg)
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
		if c.opts.WriteDeadline > 0 {
			if err := c.conn.SetWriteDeadline(time.Now().Add(c.opts.WriteDeadline * time.Second)); err != nil {
				return err
			}
		}

		if err := c.conn.WriteMessage(msg.Type, msg.Data); err != nil {
			return err
		}

		c.outed <- msg
		return nil
	}
}

// 1、当OnRead或OnWrite抛出panic后，程序被中断，但此时in或outed队列中还有没处理完的数据；
// 2、正常中断待处理的in和outed队列；
func (c *Conn) readProcess() {
	for msg := range c.in {
		func(msg *WSMessage) {
			defer func() {
				if r := recover(); r != nil {
					DefaultLogger.Error(r)
				}
			}()
			c.reader.OnRead(msg)
		}(msg)
	}
}

func (c *Conn) writeProcess() {
	for msg := range c.outed {
		func(msg *WSMessage) {
			defer func() {
				if r := recover(); r != nil {
					DefaultLogger.Error(r)
				}
			}()
			c.writer.OnWrite(msg)
		}(msg)
	}
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
