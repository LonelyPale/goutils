package http

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/lonelypale/goutils/thirdparty/gin/middleware"
)

// SetupRouter interface for dealing with api path
type SetupRouter interface {
	Setup(engine *gin.Engine)
}

type ServerOptions struct {
	// in the form "host:port". If empty, ":http" (port 80) or ":https" (port 443) is used.
	Addr     string
	TLS      bool
	CertFile string
	KeyFile  string
}

type Server struct {
	opts   ServerOptions
	engine *gin.Engine
}

func NewServer(opts ServerOptions) *Server {
	gin.DefaultErrorWriter = new(errWriter)
	engine := gin.Default()
	engine.Use(middleware.Cors())

	if opts.Addr == "" {
		if opts.TLS {
			opts.Addr = "0.0.0.0:443"
		} else {
			opts.Addr = "0.0.0.0:80"
		}
	}
	return &Server{opts: opts, engine: engine}
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}

func (s *Server) Run() {
	var err error
	if s.opts.TLS {
		err = s.engine.RunTLS(s.opts.Addr, s.opts.CertFile, s.opts.KeyFile)
	} else {
		err = s.engine.Run(s.opts.Addr)
	}

	if err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}

func (s *Server) AddRouter(router SetupRouter) {
	router.Setup(s.engine)
}

type errWriter struct{}

func (errWriter) Write(p []byte) (n int, err error) {
	info := "[GIN]" + string(p)
	log.Errorf(info)
	return len([]byte(info)), nil
}
