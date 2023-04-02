package gin

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type IRouter interface {
	Setup(*gin.Engine)
}

type Server struct {
	addr   string
	engine *gin.Engine
}

func NewServer(addr string, router IRouter) *Server {
	server := &Server{addr: addr}
	setupEngine(server)
	server.AddRouter(router)
	return server
}

func (s *Server) Run() {
	log.Infof("server run: %s", s.addr)
	log.Fatalf("failed to server: %v", s.engine.Run(s.addr))
}

type errWriter struct{}

func (errWriter) Write(p []byte) (n int, err error) {
	info := "[gin-inner]" + string(p)
	log.Errorf(info)
	return len([]byte(info)), nil
}

func setupEngine(server *Server) {
	gin.DefaultErrorWriter = new(errWriter)

	engine := gin.Default()
	engine.Use(middleware())
	server.engine = engine
}

func (s *Server) AddRouter(router IRouter) {
	router.Setup(s.engine)
}

func middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		startTime := time.Now()

		c.Next()

		log.Infof("[GIN] %3d | %13v | %15s | %s | %s",
			c.Writer.Status(),
			time.Now().Sub(startTime),
			c.ClientIP(),
			c.Request.Method,
			c.Request.RequestURI,
		)
	}
}
