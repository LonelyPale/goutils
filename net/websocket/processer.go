package websocket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/panjf2000/ants/v2"

	"github.com/LonelyPale/goutils/errors"
)

type reader struct {
	hub  Hub
	conn *Conn
	pool *ants.PoolWithFunc
}

type ReaderFactory func(hub Hub) Reader

func NewReader(hub Hub) Reader {
	return &reader{
		hub: hub,
	}
}

func (r *reader) SetConn(conn *Conn) error {
	opts := r.hub.Options()
	pool, err := ants.NewPoolWithFunc(opts.ReadPoolSize, func(i interface{}) {
		readProcess(i, conn, r.hub)
	})
	if err != nil {
		return err
	}

	r.conn = conn
	r.pool = pool
	return nil
}

func (r *reader) OnRead(msg *WSMessage) {
	if err := r.pool.Invoke(msg); err != nil {
		r.OnError(err, msg)
	}
}

func (r *reader) OnError(err error, msg *WSMessage) {
	if err != nil {
		DefaultLogger.Error(err, msg)
	}
}

func (r *reader) OnClose() {
	r.hub.UnRegister(r.conn)
	r.pool.Release()
}

type writer struct {
	hub  Hub
	conn *Conn
	//pool *ants.PoolWithFunc
}

type WriterFactory func(hub Hub) Writer

func NewWriter(hub Hub) Writer {
	return &writer{
		hub: hub,
	}
}

func (w *writer) SetConn(conn *Conn) error {
	//opts := w.hub.Options()
	//pool, err := ants.NewPoolWithFunc(opts.WritePoolSize, func(i interface{}) {
	//	//writeProcess(i, conn, w.hub)
	//})
	//if err != nil {
	//	return err
	//}

	w.conn = conn
	//w.pool = pool
	return nil
}

// todo: 待优化
func (w *writer) OnWrite(msg *WSMessage) {
	//if err := w.pool.Invoke(msg); err != nil {
	//	w.OnError(err, msg)
	//}
}

func (w *writer) OnError(err error, msg *WSMessage) {
	if err != nil {
		DefaultLogger.Error(err, msg)
	}
}

func (w *writer) OnClose() {
	w.hub.UnRegister(w.conn)
	//w.pool.Release()
}

func readProcess(i interface{}, conn *Conn, hub Hub) {
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
			DefaultLogger.Error(err, raw)
			return
		}

		msg := raw.Msg()

		handler := hub.ReaderRoute(msg.Type)
		if handler == nil {
			DefaultLogger.Error(errors.New("unknown message type"))
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

// Deprecated
func writeProcess(i interface{}, conn *Conn, hub Hub) {
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
			DefaultLogger.Error(err, raw)
			return
		}

		msg := raw.Msg()

		handler := hub.WriterRoute(msg.Type)
		if handler == nil {
			DefaultLogger.Error(errors.New("unknown message type"))
			return
		}

		handler.Invoke(conn, msg)
	}
}
