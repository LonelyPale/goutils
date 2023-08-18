package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-spring/spring-web"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-multierror"
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
	mu               sync.RWMutex
	config           *Config                       //配置选项
	register         chan *Conn                    //上线注册连接
	unregister       chan *Conn                    //下线注销连接
	router           map[string]Handler            //所有类型的处理方法
	protocols        map[string]map[*Conn]struct{} //所有在线客户端的连接子协议和连接内存地址
	processorFactory ProcessorFactory              //处理方法工厂
}

func NewContainer(processorFactory ProcessorFactory, config *Config) *Container {
	return &Container{
		config:           config,
		register:         make(chan *Conn),
		unregister:       make(chan *Conn),
		router:           make(map[string]Handler),
		protocols:        make(map[string]map[*Conn]struct{}),
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

	if err := conn.Open(w, r, processor); err != nil {
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
				c.mu.Lock()
				conns, ok := c.protocols[conn.protocol]
				if ok && conns != nil {
					conns[conn] = struct{}{}
				} else {
					c.protocols[conn.protocol] = map[*Conn]struct{}{
						conn: {},
					}
				}
				c.mu.Unlock()
			case conn := <-c.unregister:
				c.mu.Lock()
				conns, ok := c.protocols[conn.protocol]
				if ok && conns != nil {
					if _, ok := conns[conn]; ok {
						delete(conns, conn)
					}

					if err := conn.Close(); err != nil {
						DefaultLogger.Error(err)
					}
				} else {
					DefaultLogger.Error(fmt.Errorf("websocket.Container.Start conn protocol not found: %s", conn.protocol))
				}
				c.mu.Unlock()
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Second * 600)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				builder := strings.Builder{}

				c.mu.RLock()
				for protocol, conns := range c.protocols {
					if protocol == "" {
						protocol = "(empty)"
					}
					builder.WriteString(fmt.Sprintf("websocket protocols info: %s=%d", protocol, len(conns)))
				}
				c.mu.RUnlock()

				info := builder.String()
				if info != "" {
					DefaultLogger.Info(info)
				}
			}
		}
	}()
}

// todo
func (c *Container) Stop() error {
	return nil
}

func (c *Container) Handle(typed string, handler Handler) {
	c.mu.Lock()
	defer c.mu.Unlock()

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
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.router[typed]
}

func (c *Container) Send(protocol string, msg *Message) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	conns, ok := c.protocols[protocol]
	if !ok || conns == nil {
		return fmt.Errorf("websocket.Container.Send conn protocol not found: %s", protocol)
	}

	bs, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	wsmsg := &WSMessage{
		Type: websocket.TextMessage,
		Data: bs,
	}

	var errs error
	for conn := range conns {
		if err := conn.Write(wsmsg); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}
