package websocket

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Options struct {
	Enable            bool          `value:"${web.server.websocket.enable:=false}"`            //是否启用WebSocket
	Origin            bool          `value:"${web.server.websocket.origin:=true}"`             //是否启用跨域
	ReadDeadline      time.Duration `value:"${web.server.websocket.read_deadline:=0}"`         //消息单次读取超时时间，单位：秒
	WriteDeadline     time.Duration `value:"${web.server.websocket.write_deadline:=0}"`        //消息单次写入超时时间，单位：秒
	ReadBufferSize    int           `value:"${web.server.websocket.read_buffer_size:=20480}"`  //connect read buffer size: 20kb
	WriteBufferSize   int           `value:"${web.server.websocket.write_buffer_size:=20480}"` //connect write buffer size: 20kb
	MaxMessageSize    int64         `value:"${web.server.websocket.max_message_size:=65535}"`  //从消息管道读取消息的最大字节: 65535 byte
	ProcessorPoolSize int           `value:"${web.server.websocket.processor_pool_size:=10}"`  //写协程池大小
	InChanSize        int           `value:"${web.server.websocket.in_chan_size:=100}"`        // 待读管道大小
	OutChanSize       int           `value:"${web.server.websocket.out_chan_size:=100}"`       // 待写管道大小
}

func DefaultOptions() *Options {
	return &Options{
		Enable:            true,
		Origin:            true,
		ReadDeadline:      0,
		WriteDeadline:     0,
		ReadBufferSize:    20480,
		WriteBufferSize:   20480,
		MaxMessageSize:    65535,
		ProcessorPoolSize: 10,
		InChanSize:        100,
		OutChanSize:       100,
	}
}

type Hub interface {
	Options() *Options
	Register(conn *Conn)
	UnRegister(conn *Conn)
	Route(typed string) Handler
}

type Handler interface {
	Invoke(conn *Conn, msg *Message)
}

type HandlerFunc func(conn *Conn, msg *Message)

func (f HandlerFunc) Invoke(conn *Conn, msg *Message) {
	f(conn, msg)
}

type Server struct {
	opts             *Options           //配置选项
	register         chan *Conn         //上线注册连接
	unregister       chan *Conn         //下线注销连接
	conns            map[*Conn]struct{} //所有在线客户端的内存地址
	router           map[string]Handler
	routerMu         sync.RWMutex
	processorFactory ProcessorFactory
}

func NewServer(processorFactory ProcessorFactory, opts *Options) *Server {
	return &Server{
		opts:             opts,
		register:         make(chan *Conn),
		unregister:       make(chan *Conn),
		conns:            make(map[*Conn]struct{}),
		router:           make(map[string]Handler),
		processorFactory: processorFactory,
	}
}

func (s *Server) Open(ctx *gin.Context) {
	conn := NewConn(s.opts)
	processor, err := s.processorFactory(conn, s)
	if err != nil {
		DefaultLogger.Error(err)
		return
	}

	if err := conn.Open(ctx.Writer, ctx.Request, nil, processor); err != nil {
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

func (s *Server) Handle(typed string, handler Handler) {
	s.routerMu.Lock()
	defer s.routerMu.Unlock()

	if typed == "" {
		panic("WebSocket: invalid typed")
	}
	if handler == nil {
		panic("WebSocket: nil handler")
	}
	if _, exist := s.router[typed]; exist {
		panic("WebSocket: multiple registrations for " + typed)
	}

	s.router[typed] = handler
}

func (s *Server) HandleFunc(typed string, handler func(*Conn, *Message)) {
	if handler == nil {
		panic("WebSocket: nil handler")
	}
	s.Handle(typed, HandlerFunc(handler))
}

func (s *Server) Route(typed string) Handler {
	s.routerMu.RLock()
	defer s.routerMu.RUnlock()
	return s.router[typed]
}
