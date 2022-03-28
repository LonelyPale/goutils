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

type ServerConfig struct {
	Debug    bool   `toml:"debug"`
	Addr     string `toml:"addr"` // in the form "host:port". If empty, ":http" (port 80) or ":https" (port 443) is used.
	TLS      bool   `toml:"tls"`
	CertFile string `toml:"cert_file"`
	KeyFile  string `toml:"key_file"`
}

type Server struct {
	config     *ServerConfig
	httpServer *http.Server
	ginEngine  *gin.Engine
}

func NewServer(conf *ServerConfig) *Server {
	if conf.Debug {
		gin.SetMode(gin.DebugMode) // 调试模式
	} else {
		gin.SetMode(gin.ReleaseMode) // 发布模式
	}

	binding.Validator = nil // 关闭 gin 的校验器
	gin.DefaultErrorWriter = new(errWriter)

	engine := gin.Default()
	engine.Use(middleware.Cors())

	if conf.Addr == "" {
		if conf.TLS {
			conf.Addr = "0.0.0.0:443"
		} else {
			conf.Addr = "0.0.0.0:80"
		}
	}

	return &Server{
		config:    conf,
		ginEngine: engine,
		httpServer: &http.Server{
			Addr:         conf.Addr,
			Handler:      engine,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  10 * time.Second,
		},
	}
}

// Start 启动 Web 容器，阻塞
func (s *Server) Start() {
	log.Infof("http server started on %s", s.config.Addr)

	var err error
	if s.config.TLS {
		err = s.httpServer.ListenAndServeTLS(s.config.CertFile, s.config.KeyFile) // http package
		//err = s.ginEngine.RunTLS(s.config.Addr, s.config.CertFile, s.config.KeyFile) // gin package
	} else {
		err = s.httpServer.ListenAndServe() // http package
		//err = s.ginEngine.Run(s.config.Addr) // gin package
	}

	if err != nil && err != http.ErrServerClosed {
		log.WithField("err", err).Panic("failed to http server[Start]")
	}
}

// Stop 停止 Web 容器，阻塞
func (s *Server) Stop() {
	log.Infof("http server stopped on %s", s.config.Addr)

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
