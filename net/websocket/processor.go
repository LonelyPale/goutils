package websocket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
)

var DefaultLogger Logger = log.New()

type Logger interface {
	Error(args ...interface{})
}

type Processor interface {
	OnMessage(msg *WSMessage)          //消息处理
	OnError(err error, msg *WSMessage) //错误处理
	OnClose()                          //关闭连接时
	OnQuit()                           //退出协程时
}

type processor struct {
	hub  Hub
	conn *Conn
	pool *ants.PoolWithFunc
}

type ProcessorFactory func(conn *Conn, hub Hub) (Processor, error)

func NewProcessor(conn *Conn, hub Hub) (Processor, error) {
	opts := hub.Options()
	pool, err := ants.NewPoolWithFunc(opts.ProcessorPoolSize, func(i interface{}) {
		ProcessMessage(i, conn, hub)
	})
	if err != nil {
		return nil, err
	}

	return &processor{
		hub:  hub,
		conn: conn,
		pool: pool,
	}, nil
}

func (p *processor) OnMessage(msg *WSMessage) {
	if err := p.pool.Invoke(msg); err != nil {
		p.OnError(err, msg)
	}
}

func (p *processor) OnError(err error, msg *WSMessage) {
	if err != nil {
		DefaultLogger.Error(err, msg)
	}
}

func (p *processor) OnClose() {
	p.hub.UnRegister(p.conn)
}

func (p *processor) OnQuit() {
	p.pool.Release()
}

var ProcessMessage = defaultProcessMessage

func defaultProcessMessage(i interface{}, conn *Conn, hub Hub) {
	defer func() {
		if r := recover(); r != nil {
			DefaultLogger.Error(r)
		}
	}()

	ws := i.(*WSMessage)
	switch ws.Type {
	case websocket.TextMessage:
		raw := &RawMessage{}
		if err := json.Unmarshal(ws.Data, raw); err != nil {
			DefaultLogger.Error(err, ws)
			return
		}

		msg := raw.Msg()

		handler := hub.Route(msg.Type)
		if handler == nil {
			DefaultLogger.Error("WebSocket: unknown message type ", msg)
			return
		}

		handler.Invoke(conn, msg)

		//case websocket.BinaryMessage:
		//case websocket.CloseMessage:
		//case websocket.PingMessage:
		//case websocket.PongMessage:
		//default:
	}
}
