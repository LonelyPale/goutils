package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"

	"github.com/lonelypale/goutils/thirdparty/gin/middleware"
)

// SetupRouter interface for dealing with api path
type SetupRouter interface {
	Setup(engine *gin.Engine)
}

type ServerOptions struct {
	// in the form "host:port". If empty, ":http" (port 80) or ":https" (port 443) is used.
	Debug    bool
	Addr     string
	TLS      bool
	CertFile string
	KeyFile  string
}

type Server struct {
	opts       ServerOptions
	httpServer *http.Server
	ginEngine  *gin.Engine
}

func NewServer(opts ServerOptions) *Server {
	if opts.Debug {
		gin.SetMode(gin.DebugMode) // 调试模式
	} else {
		gin.SetMode(gin.ReleaseMode) // 发布模式
	}

	binding.Validator = nil // 关闭 gin 的校验器
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

	return &Server{
		opts:      opts,
		ginEngine: engine,
		httpServer: &http.Server{
			Addr:         opts.Addr,
			Handler:      engine,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  10 * time.Second,
		},
	}
}

// Start 启动 Web 容器，阻塞
func (s *Server) Start() {
	log.Infof("http server started on %s", s.opts.Addr)

	var err error
	if s.opts.TLS {
		err = s.httpServer.ListenAndServeTLS(s.opts.CertFile, s.opts.KeyFile) // http package
		//err = s.ginEngine.RunTLS(s.opts.Addr, s.opts.CertFile, s.opts.KeyFile) // gin package
	} else {
		err = s.httpServer.ListenAndServe() // http package
		//err = s.ginEngine.Run(s.opts.Addr) // gin package
	}

	if err != nil && err != http.ErrServerClosed {
		log.WithField("err", err).Panic("failed to http server[Start]")
	}
}

// Stop 停止 Web 容器，阻塞
func (s *Server) Stop() {
	log.Infof("http server stopped on %s", s.opts.Addr)

	if err := s.httpServer.Shutdown(context.TODO()); err != nil {
		log.WithField("err", err).Error("failed to http server[Stop]")
	}
}

func (s *Server) Engine() *gin.Engine {
	return s.ginEngine
}

func (s *Server) AddRouter(router SetupRouter) {
	router.Setup(s.ginEngine)
}

type errWriter struct{}

func (errWriter) Write(p []byte) (n int, err error) {
	info := "[GIN]" + string(p)
	log.Errorf(info)
	return len([]byte(info)), nil
}
