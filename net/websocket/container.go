package websocket

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-spring/spring-web"
)

type Hub interface {
	Config() *Config
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

type Container struct {
	config           *Config            //配置选项
	register         chan *Conn         //上线注册连接
	unregister       chan *Conn         //下线注销连接
	conns            map[*Conn]struct{} //所有在线客户端的内存地址
	router           map[string]Handler
	routerMu         sync.RWMutex
	processorFactory ProcessorFactory
}

func NewContainer(processorFactory ProcessorFactory, config *Config) *Container {
	return &Container{
		config:           config,
		register:         make(chan *Conn),
		unregister:       make(chan *Conn),
		conns:            make(map[*Conn]struct{}),
		router:           make(map[string]Handler),
		processorFactory: processorFactory,
	}
}

func (c *Container) OpenConn(w http.ResponseWriter, r *http.Request) {
	conn := NewConn(c.config)
	processor, err := c.processorFactory(conn, c)
	if err != nil {
		DefaultLogger.Error(err)
		return
	}

	if err := conn.Open(w, r, nil, processor); err != nil {
		DefaultLogger.Error(err)
		return
	}

	c.register <- conn
}

func (c *Container) OpenConnGin(ctx *gin.Context) {
	c.OpenConn(ctx.Writer, ctx.Request)
}

func (c *Container) OpenConnSpring(ctx SpringWeb.WebContext) {
	c.OpenConn(ctx.ResponseWriter(), ctx.Request())
}

func (c *Container) Config() *Config {
	return c.config
}

func (c *Container) Register(conn *Conn) {
	c.register <- conn
}

func (c *Container) UnRegister(conn *Conn) {
	c.unregister <- conn
}

func (c *Container) Start() {
	go func() {
		for {
			select {
			case conn := <-c.register:
				c.conns[conn] = struct{}{}
			case conn := <-c.unregister:
				if _, ok := c.conns[conn]; ok {
					delete(c.conns, conn)
				}

				if err := conn.Close(); err != nil {
					DefaultLogger.Error(err)
				}
			}
		}
	}()
}

//todo
func (c *Container) Stop() error {
	return nil
}

func (c *Container) Handle(typed string, handler Handler) {
	c.routerMu.Lock()
	defer c.routerMu.Unlock()

	if typed == "" {
		panic("WebSocket: invalid typed")
	}
	if handler == nil {
		panic("WebSocket: nil handler")
	}
	if _, exist := c.router[typed]; exist {
		panic("WebSocket: multiple registrations for " + typed)
	}

	c.router[typed] = handler
}

func (c *Container) HandleFunc(typed string, handler func(*Conn, *Message)) {
	if handler == nil {
		panic("WebSocket: nil handler")
	}
	c.Handle(typed, HandlerFunc(handler))
}

func (c *Container) Route(typed string) Handler {
	c.routerMu.RLock()
	defer c.routerMu.RUnlock()
	return c.router[typed]
}
