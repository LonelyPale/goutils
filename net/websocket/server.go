package websocket

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Options struct {
	Enable          bool          `value:"${web.server.websocket.enable:=false}"`            //是否启用WebSocket
	Origin          bool          `value:"${web.server.websocket.origin:=true}"`             //是否启用跨域
	ReadDeadline    time.Duration `value:"${web.server.websocket.read_deadline:=0}"`         //消息单次读取超时时间，单位：秒
	WriteDeadline   time.Duration `value:"${web.server.websocket.write_deadline:=0}"`        //消息单次写入超时时间，单位：秒
	ReadBufferSize  int           `value:"${web.server.websocket.read_buffer_size:=20480}"`  //connect read buffer size: 20kb
	WriteBufferSize int           `value:"${web.server.websocket.write_buffer_size:=20480}"` //connect write buffer size: 20kb
	ReadPoolSize    int           `value:"${web.server.websocket.read_pool_size:=10}"`       //读协程池大小
	WritePoolSize   int           `value:"${web.server.websocket.write_pool_size:=10}"`      //写协程池大小
	MaxMessageSize  int64         `value:"${web.server.websocket.max_message_size:=65535}"`  //从消息管道读取消息的最大字节: 65535 byte
	InChanSize      int           `value:"${web.server.websocket.in_chan_size:=100}"`        // 待读管道大小
	OutChanSize     int           `value:"${web.server.websocket.out_chan_size:=100}"`       // 待写管道大小
	OutedChanSize   int           `value:"${web.server.websocket.outed_chan_size:=100}"`     // 已写管道大小
}

type Hub interface {
	Options() *Options
	Register(conn *Conn)
	UnRegister(conn *Conn)
	ReaderRoute(typed string) Handler
	WriterRoute(typed string) Handler
}

type Handler interface {
	Invoke(conn *Conn, msg *Message)
}

type HandlerFunc func(conn *Conn, msg *Message)

func (f HandlerFunc) Invoke(conn *Conn, msg *Message) {
	f(conn, msg)
}

type Server struct {
	opts           *Options           //配置选项
	register       chan *Conn         //上线注册连接
	unregister     chan *Conn         //下线注销连接
	conns          map[*Conn]struct{} //所有在线客户端的内存地址
	readerFactory  ReaderFactory
	writerFactory  WriterFactory
	readerRouter   map[string]Handler
	readerRouterMu sync.RWMutex
	writerRouter   map[string]Handler
	writerRouterMu sync.RWMutex
}

func NewServer(readerFactory ReaderFactory, writerFactory WriterFactory, opts *Options) *Server {
	return &Server{
		opts:          opts,
		register:      make(chan *Conn),
		unregister:    make(chan *Conn),
		conns:         make(map[*Conn]struct{}),
		readerFactory: readerFactory,
		writerFactory: writerFactory,
		readerRouter:  make(map[string]Handler),
		writerRouter:  make(map[string]Handler),
	}
}

func (s *Server) Open(ctx *gin.Context) {
	reader := s.readerFactory(s)
	writer := s.writerFactory(s)
	conn := NewConn(reader, writer, s.opts)

	if err := conn.Open(ctx); err != nil {
		DefaultLogger.Error(err)
		return
	}

	s.register <- conn
}

func (s *Server) Options() *Options {
	return s.opts
}

func (s *Server) Register(conn *Conn) {
	s.register <- conn
}

func (s *Server) UnRegister(conn *Conn) {
	s.unregister <- conn
}

func (s *Server) Run() {
	for {
		select {
		case conn := <-s.register:
			s.conns[conn] = struct{}{}
			conn.Start()
		case conn := <-s.unregister:
			if _, ok := s.conns[conn]; ok {
				delete(s.conns, conn)
			}

			if err := conn.Close(); err != nil {
				DefaultLogger.Error(err)
			}
		}
	}
}

func (s *Server) Close() error {
	return nil
}

func (s *Server) ReaderHandle(typed string, handler Handler) {
	s.readerRouterMu.Lock()
	defer s.readerRouterMu.Unlock()

	if typed == "" {
		panic("WebSocket: invalid typed")
	}
	if handler == nil {
		panic("WebSocket: nil handler")
	}
	if _, exist := s.readerRouter[typed]; exist {
		panic("WebSocket: multiple registrations for " + typed)
	}

	s.readerRouter[typed] = handler
}

func (s *Server) ReaderHandleFunc(typed string, handler func(*Conn, *Message)) {
	if handler == nil {
		panic("WebSocket: nil handler")
	}
	s.ReaderHandle(typed, HandlerFunc(handler))
}

func (s *Server) WriterHandle(typed string, handler Handler) {
	s.writerRouterMu.Lock()
	defer s.writerRouterMu.Unlock()

	if typed == "" {
		panic("WebSocket: invalid typed")
	}
	if handler == nil {
		panic("WebSocket: nil handler")
	}
	if _, exist := s.writerRouter[typed]; exist {
		panic("WebSocket: multiple registrations for " + typed)
	}

	s.writerRouter[typed] = handler
}

func (s *Server) WriterHandleFunc(typed string, handler func(*Conn, *Message)) {
	if handler == nil {
		panic("WebSocket: nil handler")
	}
	s.WriterHandle(typed, HandlerFunc(handler))
}

func (s *Server) ReaderRoute(typed string) Handler {
	s.readerRouterMu.RLock()
	defer s.readerRouterMu.RUnlock()
	return s.readerRouter[typed]
}

func (s *Server) WriterRoute(typed string) Handler {
	s.writerRouterMu.RLock()
	defer s.writerRouterMu.RUnlock()
	return s.writerRouter[typed]
}
