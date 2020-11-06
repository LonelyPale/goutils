package springweb

import (
	"github.com/LonelyPale/goutils/thirdparty/gin/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-spring/spring-boot"
	"github.com/go-spring/spring-gin"
	"github.com/go-spring/spring-web"
	"github.com/go-spring/starter-web"

	"github.com/LonelyPale/goutils"
)

type WebServerConfig struct {
	StarterWeb.WebServerConfig
	Session WebSessionConfig
}

func init() {
	SpringBoot.RegisterNameBeanFn("web-container", func(config WebServerConfig) SpringWeb.WebContainer {
		container := SpringGin.NewContainer(SpringWeb.ContainerConfig{
			Port: config.Port,
		})
		return ginHandler(container, config)
	}).ConditionOnOptionalPropertyValue("web.server.enable", true)

	SpringBoot.RegisterNameBeanFn("ssl-web-container", func(config WebServerConfig) SpringWeb.WebContainer {
		home := SpringBoot.GetStringProperty("app.home")
		key := goutils.Rootify(config.SSLKey, home)
		cert := goutils.Rootify(config.SSLCert, home)
		container := SpringGin.NewContainer(SpringWeb.ContainerConfig{
			EnableSSL: true,
			Port:      config.SSLPort,
			KeyFile:   key,
			CertFile:  cert,
		})
		return ginHandler(container, config)
	}).ConditionOnPropertyValue("web.server.ssl.enable", true)
}

func ginHandler(container *SpringGin.Container, config WebServerConfig) *SpringGin.Container {
	// 使用 gin 原生的中间件
	fLogger := SpringGin.Filter(gin.Logger())
	container.SetLoggerFilter(fLogger)

	// 使用 gin 原生的中间件
	fRecover := SpringGin.Filter(gin.Recovery())
	container.SetRecoveryFilter(fRecover)

	// gin session 中间件
	if config.Session.Enable {
		fSession := SessionFilter(config.Session)
		container.AddFilter(fSession)
	}

	container.AddFilter(SpringGin.Filter(middleware.Cors()))

	return container
}
